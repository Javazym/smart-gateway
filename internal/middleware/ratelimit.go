package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
	"smart-gateway/internal/util/limiter"
	"smart-gateway/internal/util/robin"
)

// RateLimit 限流中间件。
type RateLimit struct {
	rs  *service.RequestService
	rdb *clients.RedisDB
}

func NewRateLimit(rs *service.RequestService, rdb *clients.RedisDB) *RateLimit {
	return &RateLimit{rs: rs, rdb: rdb}
}

func (r *RateLimit) LimitRate() gin.HandlerFunc {
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
		policy, err := r.rs.HandlerRateRequest(c.Request.Context(), route)
		// 查询失败、未配置策略或策略停用时放行，交由后续链路处理。
		if err != nil || policy == nil || policy.Status == 0 {
			c.Next()
			return
		}

		l := limiter.NewLimiter(r.rdb, policy)
		if l == nil {
			// 未知限流算法，放行。
			c.Next()
			return
		}

		key := rateLimitKey(c, policy)
		allow, err := l.Allow(c.Request.Context(), key)
		if err != nil {
			// Redis 异常时按放行处理，避免限流组件故障拖垮链路。
			c.Next()
			return
		}
		if !allow {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"msg": "当前系统繁忙，请稍后再试"})
			return
		}
		c.Next()
	}
}

// rateLimitKey 依据策略的限流维度生成 Redis key。
func rateLimitKey(c *gin.Context, policy *model.RateLimitPolicy) string {
	_, matchPattern, _ := util.SplitUrlByHostSlash(c.Request.RequestURI)
	switch policy.KeyType {
	case "ip":
		return "rate:" + policy.Name + ":ip:" + robin.GetClientIP(c.Request)
	case "route":
		return "rate:" + policy.Name + ":route:" + matchPattern
	case "device":
		// 设备指纹已由 AddDevice 中间件写入上下文，缺失时退化为按 IP 限流。
		if v, ok := c.Get(string(DeviceFingerprintKey)); ok {
			if fp, ok := v.(string); ok && fp != "" {
				return "rate:" + policy.Name + ":device:" + fp
			}
		}
		return "rate:" + policy.Name + ":ip:" + robin.GetClientIP(c.Request)
	default:
		return "rate:" + policy.Name + ":" + matchPattern
	}
}
