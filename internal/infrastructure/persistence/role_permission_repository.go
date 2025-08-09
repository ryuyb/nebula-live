package persistence

import (
	"context"
	"nebula-live/ent"
	"nebula-live/ent/permission"
	"nebula-live/ent/role"
	"nebula-live/ent/rolepermission"
	"nebula-live/ent/userrole"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

type rolePermissionRepository struct {
	client *ent.Client
}

// NewRolePermissionRepository 创建角色权限仓储实例
func NewRolePermissionRepository(client *ent.Client) repository.RolePermissionRepository {
	return &rolePermissionRepository{client: client}
}

func (r *rolePermissionRepository) AssignPermission(ctx context.Context, rolePermission *entity.RolePermission) (*entity.RolePermission, error) {
	create := r.client.RolePermission.
		Create().
		SetRoleID(rolePermission.RoleID).
		SetPermissionID(rolePermission.PermissionID)
	
	// 只有当AssignedBy不为0时才设置
	if rolePermission.AssignedBy != 0 {
		create = create.SetAssignedBy(rolePermission.AssignedBy)
	}
	
	created, err := create.Save(ctx)
	
	if err != nil {
		logger.Error("Failed to assign permission to role", 
			zap.Uint("role_id", rolePermission.RoleID), 
			zap.Uint("permission_id", rolePermission.PermissionID), 
			zap.Error(err))
		return nil, err
	}

	return &entity.RolePermission{
		ID:           created.ID,
		RoleID:       created.RoleID,
		PermissionID: created.PermissionID,
		AssignedBy:   created.AssignedBy,
		AssignedAt:   created.AssignedAt,
	}, nil
}

func (r *rolePermissionRepository) RemovePermission(ctx context.Context, roleID, permissionID uint) error {
	_, err := r.client.RolePermission.
		Delete().
		Where(
			rolepermission.RoleID(roleID),
			rolepermission.PermissionID(permissionID),
		).
		Exec(ctx)
	
	if err != nil {
		logger.Error("Failed to remove permission from role", 
			zap.Uint("role_id", roleID), 
			zap.Uint("permission_id", permissionID), 
			zap.Error(err))
		return err
	}

	return nil
}

func (r *rolePermissionRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error) {
	permissions, err := r.client.Permission.
		Query().
		Where(permission.HasRolePermissionsWith(rolepermission.RoleID(roleID))).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get role permissions", 
			zap.Uint("role_id", roleID), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Permission, len(permissions))
	for i, permEnt := range permissions {
		result[i] = &entity.Permission{
			ID:          permEnt.ID,
			Name:        permEnt.Name,
			DisplayName: permEnt.DisplayName,
			Description: permEnt.Description,
			Resource:    permEnt.Resource,
			Action:      permEnt.Action,
			IsSystem:    permEnt.IsSystem,
			CreatedAt:   permEnt.CreatedAt,
			UpdatedAt:   permEnt.UpdatedAt,
		}
	}
	
	return result, nil
}

func (r *rolePermissionRepository) GetPermissionRoles(ctx context.Context, permissionID uint) ([]*entity.Role, error) {
	roles, err := r.client.Role.
		Query().
		Where(role.HasRolePermissionsWith(rolepermission.PermissionID(permissionID))).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get permission roles", 
			zap.Uint("permission_id", permissionID), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Role, len(roles))
	for i, roleEnt := range roles {
		result[i] = &entity.Role{
			ID:          roleEnt.ID,
			Name:        roleEnt.Name,
			DisplayName: roleEnt.DisplayName,
			Description: roleEnt.Description,
			IsSystem:    roleEnt.IsSystem,
			CreatedAt:   roleEnt.CreatedAt,
			UpdatedAt:   roleEnt.UpdatedAt,
		}
	}
	
	return result, nil
}

func (r *rolePermissionRepository) HasPermission(ctx context.Context, roleID, permissionID uint) (bool, error) {
	exists, err := r.client.RolePermission.
		Query().
		Where(
			rolepermission.RoleID(roleID),
			rolepermission.PermissionID(permissionID),
		).
		Exist(ctx)
	
	if err != nil {
		logger.Error("Failed to check role permission", 
			zap.Uint("role_id", roleID), 
			zap.Uint("permission_id", permissionID), 
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (r *rolePermissionRepository) HasPermissionByName(ctx context.Context, roleID uint, permissionName string) (bool, error) {
	exists, err := r.client.RolePermission.
		Query().
		Where(
			rolepermission.RoleID(roleID),
			rolepermission.HasPermissionWith(permission.Name(permissionName)),
		).
		Exist(ctx)
	
	if err != nil {
		logger.Error("Failed to check role permission by name", 
			zap.Uint("role_id", roleID), 
			zap.String("permission_name", permissionName), 
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (r *rolePermissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*entity.Permission, error) {
	permissions, err := r.client.Permission.
		Query().
		Where(
			permission.HasRolePermissionsWith(
				rolepermission.HasRoleWith(
					role.HasUserRolesWith(userrole.UserID(userID)),
				),
			),
		).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get user permissions", 
			zap.Uint("user_id", userID), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Permission, len(permissions))
	for i, permEnt := range permissions {
		result[i] = &entity.Permission{
			ID:          permEnt.ID,
			Name:        permEnt.Name,
			DisplayName: permEnt.DisplayName,
			Description: permEnt.Description,
			Resource:    permEnt.Resource,
			Action:      permEnt.Action,
			IsSystem:    permEnt.IsSystem,
			CreatedAt:   permEnt.CreatedAt,
			UpdatedAt:   permEnt.UpdatedAt,
		}
	}
	
	return result, nil
}

func (r *rolePermissionRepository) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	exists, err := r.client.Permission.
		Query().
		Where(
			permission.Resource(resource),
			permission.Action(action),
			permission.HasRolePermissionsWith(
				rolepermission.HasRoleWith(
					role.HasUserRolesWith(userrole.UserID(userID)),
				),
			),
		).
		Exist(ctx)
	
	if err != nil {
		logger.Error("Failed to check user permission", 
			zap.Uint("user_id", userID), 
			zap.String("resource", resource), 
			zap.String("action", action), 
			zap.Error(err))
		return false, err
	}

	return exists, nil
}