package model

import "time"

// RequestLog 请求日志表（实时折线图与审计）。
type RequestLog struct {
	ID                  uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                           // 主键ID
	RequestID           string    `gorm:"column:request_id;type:varchar(64);not null" json:"request_id"`                               // 请求唯一ID
	RouteID             uint64    `gorm:"column:route_id;type:bigint unsigned;index:idx_route" json:"route_id"`                        // 命中的路由
	DeviceFingerprintID string    `gorm:"column:device_fingerprint_id;type:varchar(64);index:idx_device" json:"device_fingerprint_id"` // 设备指纹
	SourceIP            string    `gorm:"column:source_ip;type:varchar(64);not null" json:"source_ip"`                                 // 来源IP
	Method              string    `gorm:"column:method;type:varchar(8);not null" json:"method"`                                        // 请求方法
	Path                string    `gorm:"column:path;type:varchar(512);not null" json:"path"`                                          // 请求路径
	UpstreamID          uint64    `gorm:"column:upstream_id;type:bigint unsigned;index:idx_upstream" json:"upstream_id"`               // 转发到的上游节点
	StatusCode          int       `gorm:"column:status_code;type:int;not null;default:0" json:"status_code"`                           // 响应状态码
	LatencyMs           int       `gorm:"column:latency_ms;type:int;not null;default:0" json:"latency_ms"`                             // 耗时(毫秒)
	ResponseSize        int64     `gorm:"column:response_size;type:bigint;not null;default:0" json:"response_size"`                    // 响应大小
	ReceivedAt          time.Time `gorm:"column:received_at;type:datetime;not null;index:idx_received_at" json:"received_at"`          // 请求时间(折线图X轴)
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`        // 创建时间
}

// RequestLogDto 请求日志查询参数。
type RequestLogDto struct {
	Page                int    `json:"page" form:"page"`
	Size                int    `json:"size" form:"size"`
	DeviceFingerprintID string `json:"device_fingerprint_id" form:"device_fingerprint_id"`
	UpstreamID          uint64 `json:"upstream_id" form:"upstream_id"`
	RouteID             uint64 `json:"route_id" form:"route_id"`
	SourceIP            string `json:"source_ip" form:"source_ip"`
}
