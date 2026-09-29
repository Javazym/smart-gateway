package model

import "time"

// CircuitBreaker 熔断配置表。
type CircuitBreaker struct {
	ID                  uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                       // 主键ID
	Name                string    `gorm:"column:name;type:varchar(64);not null;uniqueIndex:uk_name" json:"name"`                   // 熔断器名
	FailureThreshold    int       `gorm:"column:failure_threshold;type:int;not null;default:5" json:"failure_threshold"`           // 连续失败阈值(触发熔断)
	SuccessThreshold    int       `gorm:"column:success_threshold;type:int;not null;default:3" json:"success_threshold"`           // 半开状态成功阈值(恢复)
	TimeoutMs           int       `gorm:"column:timeout_ms;type:int;not null;default:3000" json:"timeout_ms"`                      // 调用超时(毫秒)
	OpenDurationSeconds int       `gorm:"column:open_duration_seconds;type:int;not null;default:30" json:"open_duration_seconds"`  // 熔断打开时长(秒)
	HalfOpenMaxRequests int       `gorm:"column:half_open_max_requests;type:int;not null;default:1" json:"half_open_max_requests"` // 半开状态最大探测请求数
	Status              int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`                             // 状态: 1启用 0停用
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`    // 创建时间
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`    // 更新时间
}

// CircuitBreakerDto 熔断配置列表查询参数。
type CircuitBreakerDto struct {
	Page   int  `json:"page" form:"page"`
	Size   int  `json:"size" form:"size"`
	Status *int `json:"status" form:"status"`
}
