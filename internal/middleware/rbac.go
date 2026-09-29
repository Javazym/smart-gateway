package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/util"
)

// RequireRole 三级角色访问控制：admin 全部，operator 除用户管理写操作外的全部，viewer 只读。
func RequireRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr, _ := role.(string)
		method := c.Request.Method
		path := c.Request.URL.Path

		switch roleStr {
		case "admin":
			c.Next()
		case "operator":
			if strings.HasPrefix(path, "/api/users") && method != http.MethodGet {
				c.AbortWithStatusJSON(http.StatusForbidden, util.Error(http.StatusForbidden, "无权限操作"))
				return
			}
			c.Next()
		case "viewer":
			if method != http.MethodGet {
				c.AbortWithStatusJSON(http.StatusForbidden, util.Error(http.StatusForbidden, "无权限操作"))
				return
			}
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, util.Error(http.StatusForbidden, "无权限操作"))
		}
	}
}

// RequireAdmin 仅允许 admin 角色访问。
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if roleStr, _ := role.(string); roleStr != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, util.Error(http.StatusForbidden, "无权限操作"))
			return
		}
		c.Next()
	}
}
