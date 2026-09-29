package model

import "time"

// RateLimitPolicy 限流策略表。
type RateLimitPolicy struct {
	ID            uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                    // 主键ID
	Name          string    `gorm:"column:name;type:varchar(64);not null;uniqueIndex:uk_name" json:"name"`                // 策略名
	Algorithm     string    `gorm:"column:algorithm;type:varchar(32);not null;default:token_bucket" json:"algorithm"`     // 算法: token_bucket/leaky_bucket/sliding_window/fixed_window
	Rate          int       `gorm:"column:rate;type:int;not null" json:"rate"`                                            // 限流阈值(次数)
	Burst         int       `gorm:"column:burst;type:int;not null;default:0" json:"burst"`                                // 突发容量(令牌桶)
	WindowSeconds int       `gorm:"column:window_seconds;type:int;not null;default:1" json:"window_seconds"`              // 时间窗口(秒)
	KeyType       string    `gorm:"column:key_type;type:varchar(16);not null;default:ip" json:"key_type"`                 // 限流维度: ip/device/route
	Status        int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`                          // 状态: 1启用 0停用
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

// RateLimitPolicyDto 限流策略列表查询参数。
type RateLimitPolicyDto struct {
	Page   int  `json:"page" form:"page"`
	Size   int  `json:"size" form:"size"`
	Status *int `json:"status" form:"status"`
}
