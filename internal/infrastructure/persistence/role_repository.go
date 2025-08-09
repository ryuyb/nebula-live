package persistence

import (
	"context"
	"nebula-live/ent"
	"nebula-live/ent/role"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

type roleRepository struct {
	client *ent.Client
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(client *ent.Client) repository.RoleRepository {
	return &roleRepository{client: client}
}

func (r *roleRepository) Create(ctx context.Context, roleEntity *entity.Role) (*entity.Role, error) {
	created, err := r.client.Role.
		Create().
		SetName(roleEntity.Name).
		SetDisplayName(roleEntity.DisplayName).
		SetNillableDescription(&roleEntity.Description).
		SetIsSystem(roleEntity.IsSystem).
		Save(ctx)

	if err != nil {
		logger.Error("Failed to create role",
			zap.String("name", roleEntity.Name),
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(created), nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	roleEnt, err := r.client.Role.
		Query().
		Where(role.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Error("Failed to get role by ID",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(roleEnt), nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	roleEnt, err := r.client.Role.
		Query().
		Where(role.Name(name)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Error("Failed to get role by name",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(roleEnt), nil
}

func (r *roleRepository) List(ctx context.Context, offset, limit int) ([]*entity.Role, error) {
	roles, err := r.client.Role.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(role.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to list roles",
			zap.Int("offset", offset),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Role, len(roles))
	for i, roleEnt := range roles {
		result[i] = r.convertToEntity(roleEnt)
	}

	return result, nil
}

func (r *roleRepository) Update(ctx context.Context, roleEntity *entity.Role) (*entity.Role, error) {
	updated, err := r.client.Role.
		UpdateOneID(roleEntity.ID).
		SetDisplayName(roleEntity.DisplayName).
		SetNillableDescription(&roleEntity.Description).
		Save(ctx)

	if err != nil {
		logger.Error("Failed to update role",
			zap.Uint("id", roleEntity.ID),
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(updated), nil
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	err := r.client.Role.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to delete role",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	return nil
}

func (r *roleRepository) GetSystemRoles(ctx context.Context) ([]*entity.Role, error) {
	roles, err := r.client.Role.
		Query().
		Where(role.IsSystem(true)).
		Order(ent.Asc(role.FieldName)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to get system roles", zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Role, len(roles))
	for i, roleEnt := range roles {
		result[i] = r.convertToEntity(roleEnt)
	}

	return result, nil
}

func (r *roleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	exists, err := r.client.Role.
		Query().
		Where(role.Name(name)).
		Exist(ctx)

	if err != nil {
		logger.Error("Failed to check role existence",
			zap.String("name", name),
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

// convertToEntity 将EntGo实体转换为领域实体
func (r *roleRepository) convertToEntity(roleEnt *ent.Role) *entity.Role {
	return &entity.Role{
		ID:          roleEnt.ID,
		Name:        roleEnt.Name,
		DisplayName: roleEnt.DisplayName,
		Description: roleEnt.Description,
		IsSystem:    roleEnt.IsSystem,
		CreatedAt:   roleEnt.CreatedAt,
		UpdatedAt:   roleEnt.UpdatedAt,
	}
}
