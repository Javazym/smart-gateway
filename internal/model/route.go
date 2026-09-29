package model

import "time"

// Route 请求转发路由规则表。
type Route struct {
	ID               uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                                  // 主键ID
	Name             string    `gorm:"column:name;type:varchar(64);not null" json:"name"`                                                  // 路由名
	MatchType        string    `gorm:"column:match_type;type:varchar(16);not null" json:"match_type"`                                      // 匹配类型: path/prefix/domain/header
	MatchPattern     string    `gorm:"column:match_pattern;type:varchar(255);not null" json:"match_pattern"`                               // 匹配表达式
	UpstreamGroupID  uint64    `gorm:"column:upstream_group_id;type:bigint unsigned;not null;index:idx_upstream" json:"upstream_group_id"` // 转发目标分组
	RateLimitID      uint64    `gorm:"column:rate_limit_id;type:bigint unsigned" json:"rate_limit_id"`                                     // 关联限流策略
	CircuitBreakerID uint64    `gorm:"column:circuit_breaker_id;type:bigint unsigned" json:"circuit_breaker_id"`                           // 关联熔断配置
	StripPrefix      string    `gorm:"column:strip_prefix;type:varchar(128)" json:"strip_prefix"`                                          // 转发时剥离的前缀
	Priority         int       `gorm:"column:priority;type:int;not null;default:0;index:idx_priority" json:"priority"`                     // 匹配优先级,值小优先
	Status           int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`                                        // 状态: 1启用 0停用
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`               // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`               // 更新时间
}

// RouteDto 路由列表查询参数。
type RouteDto struct {
	Page   int    `json:"page" form:"page"`
	Size   int    `json:"size" form:"size"`
	Status *int   `json:"status" form:"status"`
	Name   string `json:"name" form:"name"`
}
