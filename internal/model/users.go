package model

import "time"

// User 管理后台用户表。
type User struct {
	ID          uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`                    // 主键ID
	Username    string     `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_username" json:"username"`    // 登录名
	Password    string     `gorm:"column:password;type:varchar(255);not null" json:"password"`                           // 密码(bcrypt哈希)
	Role        string     `gorm:"column:role;type:varchar(32);not null;default:admin" json:"role"`                      // 角色: admin/operator/viewer
	Status      int8       `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`                          // 状态: 1启用 0禁用
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:datetime" json:"last_login_at"`                              // 最近登录时间
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录成功返回的令牌与用户信息。
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserDto struct {
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
	Status   *int   `json:"status" form:"status"`
}

// UserCreateRequest 新建用户请求。
type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
	Status   int8   `json:"status"`
}

// UserUpdateRequest 更新用户请求，Password 非空则重置密码。
type UserUpdateRequest struct {
	Role     string `json:"role"`
	Status   *int8  `json:"status"`
	Password string `json:"password"`
}
