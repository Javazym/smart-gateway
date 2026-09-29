package model

import "time"

// Ban 封禁/黑名单记录表。
type Ban struct {
	ID          uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                    // 主键ID
	DeviceID    *uint64    `gorm:"column:device_id;type:bigint unsigned" json:"device_id"`                               // 封禁的设备ID
	BanType     string     `gorm:"column:ban_type;type:varchar(16);not null" json:"ban_type"`                            // 封禁类型: device/ip
	TargetValue string     `gorm:"column:target_value;type:varchar(128);not null;index:idx_target" json:"target_value"`  // 封禁目标(指纹ID或IP)
	Reason      *string    `gorm:"column:reason;type:varchar(255)" json:"reason"`                                        // 封禁原因
	OperatorID  *uint64    `gorm:"column:operator_id;type:bigint unsigned" json:"operator_id"`                           // 操作人ID
	Status      int8       `gorm:"column:status;type:tinyint;not null;default:1;index:idx_status" json:"status"`         // 状态: 1生效 0已解除
	BannedAt    time.Time  `gorm:"column:banned_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"banned_at"`   // 封禁时间
	UnbannedAt  *time.Time `gorm:"column:unbanned_at;type:datetime" json:"unbanned_at"`                                  // 解封时间
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
}

// BanDto 封禁列表查询参数。
type BanDto struct {
	Page   int  `json:"page" form:"page"`
	Size   int  `json:"size" form:"size"`
	Status *int `json:"status" form:"status"`
}
