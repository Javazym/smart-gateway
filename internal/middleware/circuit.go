package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util/circuit"
)

// Circuit 熔断中间件：请求前判定是否放行，请求后上报成功/失败。
type Circuit struct {
	rs  *service.RequestService
	rdb *clients.RedisDB
}

func NewCircuit(rs *service.RequestService, rdb *clients.RedisDB) *Circuit {
	return &Circuit{rs: rs, rdb: rdb}
}

func (cir *Circuit) HandleCircuit() gin.HandlerFunc {
	return func(c *gin.Context) {
		routeRaw, ok := c.Get(string(RouteKey))
		if !ok {
			c.Next()
			return
		}
		route, ok := routeRaw.(*model.Route)
		if !ok || route == nil {
			c.Next()
			return
		}
		cfg, err := cir.rs.HandleCircuitRequest(c.Request.Context(), route)
		// 未配置熔断器、配置被停用或查询失败时放行，交由后续链路处理。
		if err != nil || cfg == nil || cfg.Status == 0 {
			c.Next()
			return
		}

		breaker := circuit.NewCBreaker(cir.rdb, *cfg)
		if !breaker.Allow(c.Request.Context()) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "服务熔断中，请稍后再试"})
			return
		}

		// 应用下游调用超时。
		if timeout := breaker.Timeout(); timeout > 0 {
			ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
			defer cancel()
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()

		// 依据下游响应状态码上报成功/失败。
		// 注意：请求 context 可能已因下游超时被取消，回调需用独立 context，否则失败不被计数。
		if c.Writer.Status() >= http.StatusInternalServerError {
			breaker.OnFail(context.Background())
		} else {
			breaker.OnSuccess(context.Background())
		}
	}
}
