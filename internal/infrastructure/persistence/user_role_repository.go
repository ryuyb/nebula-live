package persistence

import (
	"context"
	"nebula-live/ent"
	"nebula-live/ent/role"
	"nebula-live/ent/user"
	"nebula-live/ent/userrole"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

type userRoleRepository struct {
	client *ent.Client
}

// NewUserRoleRepository 创建用户角色仓储实例
func NewUserRoleRepository(client *ent.Client) repository.UserRoleRepository {
	return &userRoleRepository{client: client}
}

func (r *userRoleRepository) AssignRole(ctx context.Context, userRole *entity.UserRole) (*entity.UserRole, error) {
	created, err := r.client.UserRole.
		Create().
		SetUserID(userRole.UserID).
		SetRoleID(userRole.RoleID).
		SetNillableAssignedBy(&userRole.AssignedBy).
		Save(ctx)
	
	if err != nil {
		logger.Error("Failed to assign role to user", 
			zap.Uint("user_id", userRole.UserID), 
			zap.Uint("role_id", userRole.RoleID), 
			zap.Error(err))
		return nil, err
	}

	return &entity.UserRole{
		ID:         created.ID,
		UserID:     created.UserID,
		RoleID:     created.RoleID,
		AssignedBy: created.AssignedBy,
		AssignedAt: created.AssignedAt,
	}, nil
}

func (r *userRoleRepository) RemoveRole(ctx context.Context, userID, roleID uint) error {
	_, err := r.client.UserRole.
		Delete().
		Where(
			userrole.UserID(userID),
			userrole.RoleID(roleID),
		).
		Exec(ctx)
	
	if err != nil {
		logger.Error("Failed to remove role from user", 
			zap.Uint("user_id", userID), 
			zap.Uint("role_id", roleID), 
			zap.Error(err))
		return err
	}

	return nil
}

func (r *userRoleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error) {
	roles, err := r.client.Role.
		Query().
		Where(role.HasUserRolesWith(userrole.UserID(userID))).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get user roles", 
			zap.Uint("user_id", userID), 
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

func (r *userRoleRepository) GetRoleUsers(ctx context.Context, roleID uint) ([]*entity.User, error) {
	users, err := r.client.User.
		Query().
		Where(user.HasUserRolesWith(userrole.RoleID(roleID))).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get role users", 
			zap.Uint("role_id", roleID), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.User, len(users))
	for i, userEnt := range users {
		result[i] = &entity.User{
			ID:        userEnt.ID,
			Username:  userEnt.Username,
			Email:     userEnt.Email,
			Password:  userEnt.Password,
			Nickname:  userEnt.Nickname,
			Avatar:    userEnt.Avatar,
			Status:    entity.UserStatus(convertUserStatus(userEnt.Status)),
			CreatedAt: userEnt.CreatedAt,
			UpdatedAt: userEnt.UpdatedAt,
		}
	}
	
	return result, nil
}

func (r *userRoleRepository) HasRole(ctx context.Context, userID, roleID uint) (bool, error) {
	exists, err := r.client.UserRole.
		Query().
		Where(
			userrole.UserID(userID),
			userrole.RoleID(roleID),
		).
		Exist(ctx)
	
	if err != nil {
		logger.Error("Failed to check user role", 
			zap.Uint("user_id", userID), 
			zap.Uint("role_id", roleID), 
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (r *userRoleRepository) HasRoleByName(ctx context.Context, userID uint, roleName string) (bool, error) {
	exists, err := r.client.UserRole.
		Query().
		Where(
			userrole.UserID(userID),
			userrole.HasRoleWith(role.Name(roleName)),
		).
		Exist(ctx)
	
	if err != nil {
		logger.Error("Failed to check user role by name", 
			zap.Uint("user_id", userID), 
			zap.String("role_name", roleName), 
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (r *userRoleRepository) GetUserRoleAssignments(ctx context.Context, userID uint) ([]*entity.UserRole, error) {
	userRoles, err := r.client.UserRole.
		Query().
		Where(userrole.UserID(userID)).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get user role assignments", 
			zap.Uint("user_id", userID), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserRole, len(userRoles))
	for i, ur := range userRoles {
		result[i] = &entity.UserRole{
			ID:         ur.ID,
			UserID:     ur.UserID,
			RoleID:     ur.RoleID,
			AssignedBy: ur.AssignedBy,
			AssignedAt: ur.AssignedAt,
		}
	}
	
	return result, nil
}

// convertUserStatus 转换用户状态
func convertUserStatus(status user.Status) int {
	switch status {
	case user.StatusActive:
		return int(entity.UserStatusActive)
	case user.StatusInactive:
		return int(entity.UserStatusInactive)
	case user.StatusBanned:
		return int(entity.UserStatusBanned)
	default:
		return int(entity.UserStatusInactive)
	}
}