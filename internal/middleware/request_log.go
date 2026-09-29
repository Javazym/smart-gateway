package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util/robin"
)

// RequestLog 请求日志中间件：在整条转发链路结束后异步落库。
type RequestLog struct {
	ls *service.RequestLogService
}

func NewRequestLog(ls *service.RequestLogService) *RequestLog {
	return &RequestLog{ls: ls}
}

func (l *RequestLog) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		origPath := c.Request.URL.Path
		c.Next()
		l.record(c, start, origPath)
	}
}

func (l *RequestLog) record(c *gin.Context, start time.Time, origPath string) {
	var routeID uint64
	if v, ok := c.Get(string(RouteKey)); ok {
		if r, ok := v.(*model.Route); ok && r != nil {
			routeID = r.ID
		}
	}
	var upstreamID uint64
	if v, ok := c.Get(string(UpstreamAddrKey)); ok {
		if n, ok := v.(*robin.Node); ok && n != nil {
			upstreamID = n.ID
		}
	}
	var fingerprintID string
	if v, ok := c.Get(string(DeviceFingerprintKey)); ok {
		fingerprintID, _ = v.(string)
	}

	entry := &model.RequestLog{
		RequestID:           newRequestID(),
		RouteID:             routeID,
		DeviceFingerprintID: fingerprintID,
		SourceIP:            robin.GetClientIP(c.Request),
		Method:              c.Request.Method,
		Path:                origPath,
		UpstreamID:          upstreamID,
		StatusCode:          c.Writer.Status(),
		LatencyMs:           int(time.Since(start).Milliseconds()),
		ResponseSize:        int64(c.Writer.Size()),
		ReceivedAt:          start,
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = l.ls.Create(ctx, entry)
	}()
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}
