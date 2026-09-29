package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util/robin"
)

// Ban 封禁拦截中间件：在请求链路中拦截被封禁的 IP 或设备指纹。
type Ban struct {
	bs *service.BanService
}

func NewBan(bs *service.BanService) *Ban {
	return &Ban{bs: bs}
}

func (b *Ban) HandleBan() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := robin.GetClientIP(c.Request)
		fingerprintID := ""
		if v, ok := c.Get(string(DeviceFingerprintKey)); ok {
			fingerprintID, _ = v.(string)
		}

		banned, err := b.bs.IsBanned(c.Request.Context(), clientIP, fingerprintID)
		if err != nil {
			// 封禁查询失败时降级放行，避免封禁组件故障拖垮链路。
			c.Next()
			return
		}
		if banned {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "请求已被封禁"})
			return
		}
		c.Next()
	}
}
