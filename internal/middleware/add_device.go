package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
	"smart-gateway/internal/util/robin"
)

// DeviceFingerprintKey 设备指纹在 gin 上下文中的 key，供限流/日志/封禁等后续环节使用。
const DeviceFingerprintKey contextKey = "device_fingerprint_id"

// AddDevice 设备识别中间件：基于 HTTP 层特征计算设备指纹并登记到 devices 表。
type AddDevice struct {
	ds *service.DeviceService
}

func NewAddDevice(ds *service.DeviceService) *AddDevice {
	return &AddDevice{ds: ds}
}

func (a *AddDevice) HandleAddDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		headerHash := util.HeaderHash(c.Request)
		userAgent := c.Request.UserAgent()
		fingerprintID := util.FingerprintID(userAgent, headerHash)
		clientIP := robin.GetClientIP(c.Request)

		device := &model.Device{
			FingerprintID:  fingerprintID,
			HTTPHeaderHash: &headerHash,
			UserAgent:      &userAgent,
			LastSeenIP:     &clientIP,
			DeviceStatus:   "normal",
			FirstSeenAt:    now,
			LastSeenAt:     now,
			UpdatedAt:      now,
		}

		if err := a.ds.Register(c.Request.Context(), device); err != nil {
			// 识别/登记失败仅降级，不阻塞转发。
			log.Printf("register device failed: %v", err)
		}

		c.Set(string(DeviceFingerprintKey), fingerprintID)
		c.Next()
	}
}
