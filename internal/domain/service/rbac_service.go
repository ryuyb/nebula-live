package service

import (
	"context"
	"errors"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"
	"time"

	"go.uber.org/zap"
)

var (
	// RBAC相关错误
	ErrRoleAlreadyExists            = errors.New("role already exists")
	ErrRoleNotFound                 = errors.New("role not found")
	ErrSystemRoleCannotDelete       = errors.New("system role cannot be deleted")
	ErrPermissionAlreadyExists      = errors.New("permission already exists")
	ErrPermissionNotFound           = errors.New("permission not found")
	ErrSystemPermissionCannotDelete = errors.New("system permission cannot be deleted")
	ErrUserRoleAlreadyExists        = errors.New("user role already exists")
	ErrUserRoleNotFound             = errors.New("user role not found")
	ErrRolePermissionAlreadyExists  = errors.New("role permission already exists")
	ErrRolePermissionNotFound       = errors.New("role permission not found")
)

// RBACService RBAC服务接口
type RBACService interface {
	// 角色管理
	CreateRole(ctx context.Context, name, displayName, description string, isSystem bool) (*entity.Role, error)
	GetRoleByID(ctx context.Context, id uint) (*entity.Role, error)
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	ListRoles(ctx context.Context, offset, limit int) ([]*entity.Role, error)
	UpdateRole(ctx context.Context, id uint, displayName, description string) (*entity.Role, error)
	DeleteRole(ctx context.Context, id uint) error

	// 权限管理
	CreatePermission(ctx context.Context, name, displayName, description, resource, action string, isSystem bool) (*entity.Permission, error)
	GetPermissionByID(ctx context.Context, id uint) (*entity.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	ListPermissions(ctx context.Context, offset, limit int) ([]*entity.Permission, error)
	UpdatePermission(ctx context.Context, id uint, displayName, description string) (*entity.Permission, error)
	DeletePermission(ctx context.Context, id uint) error

	// 用户角色管理
	AssignRoleToUser(ctx context.Context, userID, roleID, assignerID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error)
	HasRole(ctx context.Context, userID uint, roleName string) (bool, error)

	// 角色权限管理
	AssignPermissionToRole(ctx context.Context, roleID, permissionID, assignerID uint) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error)

	// 权限验证
	HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]*entity.Permission, error)

	// 初始化系统数据
	InitializeSystemData(ctx context.Context) error
}

type rbacService struct {
	roleRepo           repository.RoleRepository
	permissionRepo     repository.PermissionRepository
	userRoleRepo       repository.UserRoleRepository
	rolePermissionRepo repository.RolePermissionRepository
}

// NewRBACService 创建RBAC服务实例
func NewRBACService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	rolePermissionRepo repository.RolePermissionRepository,
) RBACService {
	return &rbacService{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		userRoleRepo:       userRoleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

// 角色管理
func (s *rbacService) CreateRole(ctx context.Context, name, displayName, description string, isSystem bool) (*entity.Role, error) {
	// 检查角色名称是否已存在
	exists, err := s.roleRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrRoleAlreadyExists
	}

	role := &entity.Role{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		IsSystem:    isSystem,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.roleRepo.Create(ctx, role)
}

func (s *rbacService) GetRoleByID(ctx context.Context, id uint) (*entity.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *rbacService) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	role, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *rbacService) ListRoles(ctx context.Context, offset, limit int) ([]*entity.Role, error) {
	return s.roleRepo.List(ctx, offset, limit)
}

func (s *rbacService) UpdateRole(ctx context.Context, id uint, displayName, description string) (*entity.Role, error) {
	role, err := s.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	role.DisplayName = displayName
	role.Description = description
	role.UpdatedAt = time.Now()

	return s.roleRepo.Update(ctx, role)
}

func (s *rbacService) DeleteRole(ctx context.Context, id uint) error {
	role, err := s.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return ErrSystemRoleCannotDelete
	}

	return s.roleRepo.Delete(ctx, id)
}

// 权限管理
func (s *rbacService) CreatePermission(ctx context.Context, name, displayName, description, resource, action string, isSystem bool) (*entity.Permission, error) {
	// 检查权限名称是否已存在
	exists, err := s.permissionRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPermissionAlreadyExists
	}

	permission := &entity.Permission{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Resource:    resource,
		Action:      action,
		IsSystem:    isSystem,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.permissionRepo.Create(ctx, permission)
}

func (s *rbacService) GetPermissionByID(ctx context.Context, id uint) (*entity.Permission, error) {
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}
	return permission, nil
}

func (s *rbacService) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	permission, err := s.permissionRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}
	return permission, nil
}

func (s *rbacService) ListPermissions(ctx context.Context, offset, limit int) ([]*entity.Permission, error) {
	return s.permissionRepo.List(ctx, offset, limit)
}

func (s *rbacService) UpdatePermission(ctx context.Context, id uint, displayName, description string) (*entity.Permission, error) {
	permission, err := s.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	permission.DisplayName = displayName
	permission.Description = description
	permission.UpdatedAt = time.Now()

	return s.permissionRepo.Update(ctx, permission)
}

func (s *rbacService) DeletePermission(ctx context.Context, id uint) error {
	permission, err := s.GetPermissionByID(ctx, id)
	if err != nil {
		return err
	}

	if permission.IsSystem {
		return ErrSystemPermissionCannotDelete
	}

	return s.permissionRepo.Delete(ctx, id)
}

// 用户角色管理
func (s *rbacService) AssignRoleToUser(ctx context.Context, userID, roleID, assignerID uint) error {
	// 检查是否已经分配
	exists, err := s.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserRoleAlreadyExists
	}

	userRole := &entity.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignerID,
		AssignedAt: time.Now(),
	}

	_, err = s.userRoleRepo.AssignRole(ctx, userRole)
	return err
}

func (s *rbacService) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return s.userRoleRepo.RemoveRole(ctx, userID, roleID)
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error) {
	return s.userRoleRepo.GetUserRoles(ctx, userID)
}

func (s *rbacService) HasRole(ctx context.Context, userID uint, roleName string) (bool, error) {
	return s.userRoleRepo.HasRoleByName(ctx, userID, roleName)
}

// 角色权限管理
func (s *rbacService) AssignPermissionToRole(ctx context.Context, roleID, permissionID, assignerID uint) error {
	// 检查是否已经分配
	exists, err := s.rolePermissionRepo.HasPermission(ctx, roleID, permissionID)
	if err != nil {
		return err
	}
	if exists {
		return ErrRolePermissionAlreadyExists
	}

	rolePermission := &entity.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		AssignedBy:   assignerID,
		AssignedAt:   time.Now(),
	}

	_, err = s.rolePermissionRepo.AssignPermission(ctx, rolePermission)
	return err
}

func (s *rbacService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	return s.rolePermissionRepo.RemovePermission(ctx, roleID, permissionID)
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error) {
	return s.rolePermissionRepo.GetRolePermissions(ctx, roleID)
}

// 权限验证
func (s *rbacService) HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	return s.rolePermissionRepo.CheckUserPermission(ctx, userID, resource, action)
}

func (s *rbacService) GetUserPermissions(ctx context.Context, userID uint) ([]*entity.Permission, error) {
	return s.rolePermissionRepo.GetUserPermissions(ctx, userID)
}

// 初始化系统数据
func (s *rbacService) InitializeSystemData(ctx context.Context) error {
	logger.Info("Initializing RBAC system data...")

	// 创建系统角色
	if err := s.createSystemRoles(ctx); err != nil {
		logger.Error("Failed to create system roles", zap.Error(err))
		return err
	}

	// 创建系统权限
	if err := s.createSystemPermissions(ctx); err != nil {
		logger.Error("Failed to create system permissions", zap.Error(err))
		return err
	}

	// 分配权限给角色
	if err := s.assignPermissionsToRoles(ctx); err != nil {
		logger.Error("Failed to assign permissions to roles", zap.Error(err))
		return err
	}

	logger.Info("RBAC system data initialized successfully")
	return nil
}

// createSystemRoles 创建系统角色
func (s *rbacService) createSystemRoles(ctx context.Context) error {
	systemRoles := []struct {
		name        string
		displayName string
		description string
	}{
		{entity.RoleNameAdmin, "管理员", "拥有系统管理权限的管理员"},
		{entity.RoleNameUser, "普通用户", "普通用户角色"},
	}

	for _, roleData := range systemRoles {
		exists, err := s.roleRepo.ExistsByName(ctx, roleData.name)
		if err != nil {
			return err
		}
		if !exists {
			_, err := s.CreateRole(ctx, roleData.name, roleData.displayName, roleData.description, true)
			if err != nil {
				return err
			}
			logger.Info("Created system role", zap.String("name", roleData.name))
		}
	}

	return nil
}

// createSystemPermissions 创建系统权限
func (s *rbacService) createSystemPermissions(ctx context.Context) error {
	systemPermissions := []struct {
		name        string
		displayName string
		description string
		resource    string
		action      string
	}{
		// 用户管理权限
		{entity.PermissionUserRead, "查看用户", "查看用户信息的权限", "user", "read"},
		{entity.PermissionUserWrite, "修改用户", "修改用户信息的权限", "user", "write"},
		{entity.PermissionUserDelete, "删除用户", "删除用户的权限", "user", "delete"},
		{entity.PermissionUserManage, "管理用户", "完全管理用户的权限", "user", "manage"},

		// 角色管理权限
		{entity.PermissionRoleRead, "查看角色", "查看角色信息的权限", "role", "read"},
		{entity.PermissionRoleWrite, "修改角色", "修改角色信息的权限", "role", "write"},
		{entity.PermissionRoleDelete, "删除角色", "删除角色的权限", "role", "delete"},
		{entity.PermissionRoleManage, "管理角色", "完全管理角色的权限", "role", "manage"},

		// 权限管理权限
		{entity.PermissionPermissionRead, "查看权限", "查看权限信息的权限", "permission", "read"},
		{entity.PermissionPermissionWrite, "修改权限", "修改权限信息的权限", "permission", "write"},
		{entity.PermissionPermissionDelete, "删除权限", "删除权限的权限", "permission", "delete"},
		{entity.PermissionPermissionManage, "管理权限", "完全管理权限的权限", "permission", "manage"},

		// 系统管理权限
		{entity.PermissionSystemManage, "系统管理", "系统管理权限", "system", "manage"},
	}

	for _, permData := range systemPermissions {
		exists, err := s.permissionRepo.ExistsByName(ctx, permData.name)
		if err != nil {
			return err
		}
		if !exists {
			_, err := s.CreatePermission(ctx, permData.name, permData.displayName, permData.description, permData.resource, permData.action, true)
			if err != nil {
				return err
			}
			logger.Info("Created system permission", zap.String("name", permData.name))
		}
	}

	return nil
}

// assignPermissionsToRoles 分配权限给角色
func (s *rbacService) assignPermissionsToRoles(ctx context.Context) error {
	// 管理员拥有所有权限
	adminRole, err := s.GetRoleByName(ctx, entity.RoleNameAdmin)
	if err != nil {
		return err
	}

	// 给管理员分配所有系统权限
	systemPermissions, err := s.permissionRepo.GetSystemPermissions(ctx)
	if err != nil {
		return err
	}

	for _, permission := range systemPermissions {
		exists, err := s.rolePermissionRepo.HasPermission(ctx, adminRole.ID, permission.ID)
		if err != nil {
			return err
		}
		if !exists {
			// 系统初始化时，使用空的assignerID（系统分配）
			rolePermission := &entity.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: permission.ID,
				AssignedAt:   time.Now(),
			}
			_, err = s.rolePermissionRepo.AssignPermission(ctx, rolePermission)
			if err != nil && err != ErrRolePermissionAlreadyExists {
				return err
			}
		}
	}

	// 普通用户拥有基本权限
	userRole, err := s.GetRoleByName(ctx, entity.RoleNameUser)
	if err != nil {
		return err
	}

	userPermissions := []string{
		entity.PermissionUserRead, // 只能查看用户信息
	}

	for _, permName := range userPermissions {
		permission, err := s.GetPermissionByName(ctx, permName)
		if err != nil {
			continue
		}

		exists, err := s.rolePermissionRepo.HasPermission(ctx, userRole.ID, permission.ID)
		if err != nil {
			return err
		}
		if !exists {
			// 系统初始化时，使用空的assignerID（系统分配）
			rolePermission := &entity.RolePermission{
				RoleID:       userRole.ID,
				PermissionID: permission.ID,
				AssignedAt:   time.Now(),
			}
			_, err = s.rolePermissionRepo.AssignPermission(ctx, rolePermission)
			if err != nil && err != ErrRolePermissionAlreadyExists {
				return err
			}
		}
	}

	return nil
}
