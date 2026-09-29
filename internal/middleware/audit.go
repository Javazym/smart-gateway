package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util/robin"
)

// Audit 审计日志中间件：对管理 API 的写操作异步记录审计日志。
type Audit struct {
	as *service.AuditLogService
}

func NewAudit(as *service.AuditLogService) *Audit {
	return &Audit{as: as}
}

func (a *Audit) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !isMutating(c.Request.Method) {
			return
		}

		var operatorID *uint64
		if v, ok := c.Get("user_id"); ok {
			if id, ok := v.(uint64); ok {
				operatorID = &id
			}
		}
		var username *string
		if v, ok := c.Get("username"); ok {
			if s, ok := v.(string); ok {
				username = &s
			}
		}

		targetType, targetID := auditTarget(c.Request.URL.Path)
		sourceIP := robin.GetClientIP(c.Request)

		entry := &model.AuditLog{
			OperatorID: operatorID,
			Username:   username,
			Action:     strings.ToLower(c.Request.Method) + "_" + targetType,
			TargetType: &targetType,
			TargetID:   targetID,
			SourceIP:   &sourceIP,
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = a.as.Create(ctx, entry)
		}()
	}
}

func isMutating(method string) bool {
	switch method {
	case "POST", "PUT", "DELETE", "PATCH":
		return true
	default:
		return false
	}
}

// auditTarget 从 /api/{resource}[/{id}] 解析出资源类型与目标 ID。
func auditTarget(path string) (string, *string) {
	p := strings.TrimPrefix(path, "/api/")
	parts := strings.Split(p, "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", nil
	}
	typ := parts[0]
	var id *string
	if len(parts) >= 2 && parts[1] != "" {
		id = &parts[1]
	}
	return typ, id
}
