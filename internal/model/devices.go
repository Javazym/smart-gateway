package model

import "time"

// Device 设备指纹识别表。
type Device struct {
	ID             uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                                   // 主键ID
	FingerprintID  string    `gorm:"column:fingerprint_id;type:varchar(64);not null;uniqueIndex:uk_fingerprint" json:"fingerprint_id"`    // 综合指纹哈希(JA4+OS+Header哈希)
	JA4Hash        *string   `gorm:"column:ja4_hash;type:varchar(64)" json:"ja4_hash"`                                                    // JA4指纹
	JA3Hash        *string   `gorm:"column:ja3_hash;type:varchar(64)" json:"ja3_hash"`                                                    // JA3指纹(兼容)
	OSFingerprint  *string   `gorm:"column:os_fingerprint;type:varchar(128)" json:"os_fingerprint"`                                       // TCP OS指纹(TTL/窗口等)
	HTTPHeaderHash *string   `gorm:"column:http_header_hash;type:varchar(64)" json:"http_header_hash"`                                    // 完整HTTP头哈希
	UserAgent      *string   `gorm:"column:user_agent;type:varchar(512)" json:"user_agent"`                                               // User-Agent
	LastSeenIP     *string   `gorm:"column:last_seen_ip;type:varchar(64);index:idx_ip" json:"last_seen_ip"`                               // 最近来源IP
	DeviceStatus   string    `gorm:"column:device_status;type:varchar(16);not null;default:normal;index:idx_status" json:"device_status"` // 设备状态: normal/suspected/banned
	BanReason      *string   `gorm:"column:ban_reason;type:varchar(255)" json:"ban_reason"`                                               // 封禁原因
	FirstSeenAt    time.Time `gorm:"column:first_seen_at;type:datetime;not null" json:"first_seen_at"`                                    // 首次出现时间
	LastSeenAt     time.Time `gorm:"column:last_seen_at;type:datetime;not null" json:"last_seen_at"`                                      // 最近出现时间
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`                // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`                // 更新时间
}

// DeviceDto 设备列表查询参数。
type DeviceDto struct {
	Page   int    `json:"page" form:"page"`
	Size   int    `json:"size" form:"size"`
	Status string `json:"status" form:"status"` // normal/suspected/banned，空表示全部
}
