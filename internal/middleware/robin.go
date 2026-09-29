package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util/robin"
)

// 自定义key，避免string类型冲突
type contextKey string

const (
	UpstreamAddrKey contextKey = "upstream_addr"
	RouteKey        contextKey = "route"
	BalancerKey     contextKey = "balancer"
)

type Robin struct {
	rs  *service.RequestService
	rdb *clients.RedisDB
}

func NewRobin(rs *service.RequestService, rdb *clients.RedisDB) *Robin {
	return &Robin{rs: rs, rdb: rdb}
}

func (r *Robin) HandleRobin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 匹配路由，未命中与上游不可用需区分状态码。
		route, err := r.rs.MatchRoute(c.Request.Context(), c.Request)
		if err != nil {
			if errors.Is(err, service.ErrRouteNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": "route not found"})
				return
			}
			log.Printf("match route failed: %v", err)
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "match route failed"})
			return
		}
		c.Set(string(RouteKey), route)

		upg, err := r.rs.HandlerRobinRequest(c.Request.Context(), route)
		if err != nil {
			log.Printf("get upstream group failed: %v", err)
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "get upstream group failed"})
			return
		}

		rb := robin.NewBalancer(upg, r.rdb)
		adr, err := rb.Pick(c.Request)
		if err != nil {
			log.Printf("pick upstream node failed: %v", err)
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "pick upstream node failed"})
			return
		}

		c.Set(string(UpstreamAddrKey), adr)
		c.Set(string(BalancerKey), rb)
		c.Next()
	}
}
