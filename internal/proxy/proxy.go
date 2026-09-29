package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Handler 创建一个把请求反向代理到 upstream 的 Gin 处理器。
func Handler(upstream string) (gin.HandlerFunc, error) {
	target, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}
	p := httputil.NewSingleHostReverseProxy(target)
	return func(c *gin.Context) {
		p.ServeHTTP(c.Writer, c.Request)
	}, nil
}
