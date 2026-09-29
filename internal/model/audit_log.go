package model

import "time"

// AuditLog 管理操作审计日志表。
type AuditLog struct {
	ID         uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                                         // 主键ID
	OperatorID *uint64   `gorm:"column:operator_id;type:bigint unsigned;index:idx_operator" json:"operator_id"`                             // 操作人ID
	Username   *string   `gorm:"column:username;type:varchar(64)" json:"username"`                                                          // 操作人登录名
	Action     string    `gorm:"column:action;type:varchar(64);not null" json:"action"`                                                     // 操作: ban_device/unban/update_ratelimit/...
	TargetType *string   `gorm:"column:target_type;type:varchar(32)" json:"target_type"`                                                    // 目标类型
	TargetID   *string   `gorm:"column:target_id;type:varchar(64)" json:"target_id"`                                                        // 目标ID
	Detail     *string   `gorm:"column:detail;type:text" json:"detail"`                                                                     // 操作详情(JSON)
	SourceIP   *string   `gorm:"column:source_ip;type:varchar(64)" json:"source_ip"`                                                        // 操作来源IP
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"` // 创建时间
}

// AuditLogDto 审计日志查询参数。
type AuditLogDto struct {
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
	Username string `json:"username" form:"username"`
	Action   string `json:"action" form:"action"`
}
