package repository

import (
	"context"
	"nebula-live/internal/domain/entity"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) (*entity.Role, error)
	
	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, id uint) (*entity.Role, error)
	
	// GetByName 根据名称获取角色
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	
	// List 获取角色列表
	List(ctx context.Context, offset, limit int) ([]*entity.Role, error)
	
	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) (*entity.Role, error)
	
	// Delete 删除角色
	Delete(ctx context.Context, id uint) error
	
	// GetSystemRoles 获取所有系统角色
	GetSystemRoles(ctx context.Context) ([]*entity.Role, error)
	
	// ExistsByName 检查角色名称是否存在
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	// Create 创建权限
	Create(ctx context.Context, permission *entity.Permission) (*entity.Permission, error)
	
	// GetByID 根据ID获取权限
	GetByID(ctx context.Context, id uint) (*entity.Permission, error)
	
	// GetByName 根据名称获取权限
	GetByName(ctx context.Context, name string) (*entity.Permission, error)
	
	// List 获取权限列表
	List(ctx context.Context, offset, limit int) ([]*entity.Permission, error)
	
	// Update 更新权限
	Update(ctx context.Context, permission *entity.Permission) (*entity.Permission, error)
	
	// Delete 删除权限
	Delete(ctx context.Context, id uint) error
	
	// GetSystemPermissions 获取所有系统权限
	GetSystemPermissions(ctx context.Context) ([]*entity.Permission, error)
	
	// ExistsByName 检查权限名称是否存在
	ExistsByName(ctx context.Context, name string) (bool, error)
	
	// GetByResource 根据资源获取权限列表
	GetByResource(ctx context.Context, resource string) ([]*entity.Permission, error)
}

// UserRoleRepository 用户角色关联仓储接口
type UserRoleRepository interface {
	// AssignRole 分配角色给用户
	AssignRole(ctx context.Context, userRole *entity.UserRole) (*entity.UserRole, error)
	
	// RemoveRole 移除用户的角色
	RemoveRole(ctx context.Context, userID, roleID uint) error
	
	// GetUserRoles 获取用户的所有角色
	GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error)
	
	// GetRoleUsers 获取角色的所有用户
	GetRoleUsers(ctx context.Context, roleID uint) ([]*entity.User, error)
	
	// HasRole 检查用户是否有指定角色
	HasRole(ctx context.Context, userID, roleID uint) (bool, error)
	
	// HasRoleByName 检查用户是否有指定名称的角色
	HasRoleByName(ctx context.Context, userID uint, roleName string) (bool, error)
	
	// GetUserRoleAssignments 获取用户角色分配记录
	GetUserRoleAssignments(ctx context.Context, userID uint) ([]*entity.UserRole, error)
}

// RolePermissionRepository 角色权限关联仓储接口
type RolePermissionRepository interface {
	// AssignPermission 分配权限给角色
	AssignPermission(ctx context.Context, rolePermission *entity.RolePermission) (*entity.RolePermission, error)
	
	// RemovePermission 移除角色的权限
	RemovePermission(ctx context.Context, roleID, permissionID uint) error
	
	// GetRolePermissions 获取角色的所有权限
	GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error)
	
	// GetPermissionRoles 获取权限的所有角色
	GetPermissionRoles(ctx context.Context, permissionID uint) ([]*entity.Role, error)
	
	// HasPermission 检查角色是否有指定权限
	HasPermission(ctx context.Context, roleID, permissionID uint) (bool, error)
	
	// HasPermissionByName 检查角色是否有指定名称的权限
	HasPermissionByName(ctx context.Context, roleID uint, permissionName string) (bool, error)
	
	// GetUserPermissions 获取用户的所有权限（通过角色）
	GetUserPermissions(ctx context.Context, userID uint) ([]*entity.Permission, error)
	
	// CheckUserPermission 检查用户是否有指定权限
	CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
}