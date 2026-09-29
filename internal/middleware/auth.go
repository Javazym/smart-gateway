package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/util"
)

// Auth 校验 Bearer JWT 令牌，通过后把用户信息写入上下文。
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Error(http.StatusUnauthorized, "unauthorized"))
			return
		}
		claims, err := util.ParseToken(jwtSecret, strings.TrimPrefix(auth, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Error(http.StatusUnauthorized, "unauthorized"))
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
