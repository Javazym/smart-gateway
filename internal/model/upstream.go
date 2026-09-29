package model

import (
	"time"
)

// UpstreamGroup 上游分组表（负载均衡）。
type UpstreamGroup struct {
	ID                  uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                      // 主键ID
	Name                string     `gorm:"column:name;type:varchar(64);not null;uniqueIndex:uk_name" json:"name"`                  // 分组名
	LBAlgorithm         string     `gorm:"column:lb_algorithm;type:varchar(32);not null;default:round_robin" json:"lb_algorithm"`  // 负载均衡算法: round_robin/weighted_rr/ip_hash/least_conn/consistent_hash
	HealthCheckPath     *string    `gorm:"column:health_check_path;type:varchar(255)" json:"health_check_path"`                    // 健康检查路径
	HealthCheckInterval int        `gorm:"column:health_check_interval;type:int;not null;default:10" json:"health_check_interval"` // 健康检查间隔(秒)
	Status              int8       `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`                            // 状态: 1启用 0停用
	CreatedAt           time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`   // 创建时间
	UpdatedAt           time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`   // 更新时间
	Upstreams           []Upstream `gorm:"foreignKey:GroupID;references:ID" json:"upstreams"`                                      // 下游节点
}

// Upstream 上游服务节点表。
type Upstream struct {
	ID           uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                    // 主键ID
	GroupID      uint64    `gorm:"column:group_id;type:bigint unsigned;not null;index:idx_group" json:"group_id"`        // 所属上游分组
	Name         string    `gorm:"column:name;type:varchar(64);not null" json:"name"`                                    // 节点名
	Scheme       string    `gorm:"column:scheme;type:varchar(8);not null;default:http" json:"scheme"`                    // 协议: http/https
	Host         string    `gorm:"column:host;type:varchar(255);not null" json:"host"`                                   // 主机
	Port         int       `gorm:"column:port;type:int;not null" json:"port"`                                            // 端口
	Weight       int       `gorm:"column:weight;type:int;not null;default:1" json:"weight"`                              // 负载均衡权重
	MaxConns     int       `gorm:"column:max_conns;type:int;not null;default:0" json:"max_conns"`                        // 最大连接数,0不限制
	HealthStatus string    `gorm:"column:health_status;type:varchar(16);not null;default:unknown" json:"health_status"`  // 健康状态: healthy/unhealthy/unknown
	Status       int       `gorm:"column:status;type:tinyint;not null;default:1;index:idx_status" json:"status"`         // 状态: 1启用 0停用
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

type GroupDto struct {
	Page   int  `json:"page" form:"page"`
	Size   int  `json:"size" form:"size"`
	Status *int `json:"status" form:"status"`
}
