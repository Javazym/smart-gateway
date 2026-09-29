package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util/robin"
)

// Proxy 反向代理转发中间件，读取负载均衡选中的上游节点并转发请求。
type Proxy struct {
	rs *service.RequestService
}

func NewProxy(rs *service.RequestService) *Proxy {
	return &Proxy{rs: rs}
}

func (p *Proxy) HandleProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodeRaw, ok := c.Get(string(UpstreamAddrKey))
		if !ok {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "no upstream address"})
			return
		}
		node, ok := nodeRaw.(*robin.Node)
		if !ok || node == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "invalid upstream address"})
			return
		}

		// 前缀剥离：命中路由配置 strip_prefix 时，转发前剥离该前缀。
		if routeRaw, ok := c.Get(string(RouteKey)); ok {
			if route, ok := routeRaw.(*model.Route); ok && route != nil && route.StripPrefix != "" {
				if strings.HasPrefix(c.Request.URL.Path, route.StripPrefix) {
					c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, route.StripPrefix)
					c.Request.URL.RawPath = ""
					if c.Request.URL.Path == "" {
						c.Request.URL.Path = "/"
					}
				}
			}
		}

		// 向下游透传真实客户端 IP。
		c.Request.Header.Set("X-Real-IP", robin.GetClientIP(c.Request))

		// 最小连接数：转发前后维护活跃连接计数。
		var balancer robin.Balancer
		if b, ok := c.Get(string(BalancerKey)); ok {
			balancer, _ = b.(robin.Balancer)
		}

		target := node.Scheme + "://" + node.Addr
		proxy := p.rs.NewProxy(target)
		if proxy == nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"msg": "invalid upstream target"})
			return
		}
		if balancer != nil {
			balancer.ConnInc(node.Addr)
			defer balancer.ConnDec(node.Addr)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
