package entity

import (
	"time"
)

// Role 角色实体
type Role struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`         // 角色名称，如：admin, user
	DisplayName string    `json:"display_name"` // 显示名称，如：管理员, 普通用户
	Description string    `json:"description"`  // 角色描述
	IsSystem    bool      `json:"is_system"`    // 是否为系统角色（系统角色不可删除）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission 权限实体
type Permission struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`         // 权限名称，如：user:read, user:write, admin:manage
	DisplayName string    `json:"display_name"` // 显示名称，如：查看用户, 修改用户, 管理系统
	Description string    `json:"description"`  // 权限描述
	Resource    string    `json:"resource"`     // 资源名称，如：user, post, system
	Action      string    `json:"action"`       // 操作名称，如：read, write, delete, manage
	IsSystem    bool      `json:"is_system"`    // 是否为系统权限（系统权限不可删除）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole 用户角色关联实体
type UserRole struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	RoleID     uint      `json:"role_id"`
	AssignedBy uint      `json:"assigned_by"` // 分配者的用户ID
	AssignedAt time.Time `json:"assigned_at"`
}

// RolePermission 角色权限关联实体
type RolePermission struct {
	ID           uint      `json:"id"`
	RoleID       uint      `json:"role_id"`
	PermissionID uint      `json:"permission_id"`
	AssignedBy   uint      `json:"assigned_by"` // 分配者的用户ID
	AssignedAt   time.Time `json:"assigned_at"`
}

// 系统预定义角色常量
const (
	RoleNameAdmin = "admin" // 管理员
	RoleNameUser  = "user"  // 普通用户
)

// 系统预定义权限常量
const (
	// 用户管理权限
	PermissionUserRead   = "user:read"
	PermissionUserWrite  = "user:write"
	PermissionUserDelete = "user:delete"
	PermissionUserManage = "user:manage"

	// 角色管理权限
	PermissionRoleRead   = "role:read"
	PermissionRoleWrite  = "role:write"
	PermissionRoleDelete = "role:delete"
	PermissionRoleManage = "role:manage"

	// 权限管理权限
	PermissionPermissionRead   = "permission:read"
	PermissionPermissionWrite  = "permission:write"
	PermissionPermissionDelete = "permission:delete"
	PermissionPermissionManage = "permission:manage"

	// 系统管理权限
	PermissionSystemManage = "system:manage"
)

// IsSystemRole 检查是否为系统角色
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}

// CanDelete 检查角色是否可以被删除
func (r *Role) CanDelete() bool {
	return !r.IsSystem
}

// IsSystemPermission 检查是否为系统权限
func (p *Permission) IsSystemPermission() bool {
	return p.IsSystem
}

// CanDelete 检查权限是否可以被删除
func (p *Permission) CanDelete() bool {
	return !p.IsSystem
}

// GetPermissionKey 获取权限的完整标识
func (p *Permission) GetPermissionKey() string {
	return p.Resource + ":" + p.Action
}
