package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // 密码不在JSON中显示
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusInactive
	UserStatusBanned
)

// String 返回用户状态的字符串表示
func (s UserStatus) String() string {
	switch s {
	case UserStatusActive:
		return "active"
	case UserStatusInactive:
		return "inactive"
	case UserStatusBanned:
		return "banned"
	default:
		return "unknown"
	}
}

// IsActive 检查用户是否处于活跃状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsBanned 检查用户是否被禁用
func (u *User) IsBanned() bool {
	return u.Status == UserStatusBanned
}

// Activate 激活用户
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

// Ban 禁用用户
func (u *User) Ban() {
	u.Status = UserStatusBanned
	u.UpdatedAt = time.Now()
}
