package persistence

import (
	"context"
	"nebula-live/ent"
	"nebula-live/ent/permission"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

type permissionRepository struct {
	client *ent.Client
}

// NewPermissionRepository 创建权限仓储实例
func NewPermissionRepository(client *ent.Client) repository.PermissionRepository {
	return &permissionRepository{client: client}
}

func (r *permissionRepository) Create(ctx context.Context, permEntity *entity.Permission) (*entity.Permission, error) {
	created, err := r.client.Permission.
		Create().
		SetName(permEntity.Name).
		SetDisplayName(permEntity.DisplayName).
		SetNillableDescription(&permEntity.Description).
		SetResource(permEntity.Resource).
		SetAction(permEntity.Action).
		SetIsSystem(permEntity.IsSystem).
		Save(ctx)
	
	if err != nil {
		logger.Error("Failed to create permission", 
			zap.String("name", permEntity.Name), 
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(created), nil
}

func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*entity.Permission, error) {
	permEnt, err := r.client.Permission.
		Query().
		Where(permission.ID(id)).
		Only(ctx)
	
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Error("Failed to get permission by ID", 
			zap.Uint("id", id), 
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(permEnt), nil
}

func (r *permissionRepository) GetByName(ctx context.Context, name string) (*entity.Permission, error) {
	permEnt, err := r.client.Permission.
		Query().
		Where(permission.Name(name)).
		Only(ctx)
	
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Error("Failed to get permission by name", 
			zap.String("name", name), 
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(permEnt), nil
}

func (r *permissionRepository) List(ctx context.Context, offset, limit int) ([]*entity.Permission, error) {
	permissions, err := r.client.Permission.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(permission.FieldCreatedAt)).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to list permissions", 
			zap.Int("offset", offset), 
			zap.Int("limit", limit), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Permission, len(permissions))
	for i, perm := range permissions {
		result[i] = r.convertToEntity(perm)
	}
	
	return result, nil
}

func (r *permissionRepository) Update(ctx context.Context, permEntity *entity.Permission) (*entity.Permission, error) {
	updated, err := r.client.Permission.
		UpdateOneID(permEntity.ID).
		SetDisplayName(permEntity.DisplayName).
		SetNillableDescription(&permEntity.Description).
		Save(ctx)
	
	if err != nil {
		logger.Error("Failed to update permission", 
			zap.Uint("id", permEntity.ID), 
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(updated), nil
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	err := r.client.Permission.
		DeleteOneID(id).
		Exec(ctx)
	
	if err != nil {
		logger.Error("Failed to delete permission", 
			zap.Uint("id", id), 
			zap.Error(err))
		return err
	}

	return nil
}

func (r *permissionRepository) GetSystemPermissions(ctx context.Context) ([]*entity.Permission, error) {
	permissions, err := r.client.Permission.
		Query().
		Where(permission.IsSystem(true)).
		Order(ent.Asc(permission.FieldName)).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get system permissions", zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Permission, len(permissions))
	for i, permEnt := range permissions {
		result[i] = r.convertToEntity(permEnt)
	}
	
	return result, nil
}

func (r *permissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	exists, err := r.client.Permission.
		Query().
		Where(permission.Name(name)).
		Exist(ctx)
	
	if err != nil {
		logger.Error("Failed to check permission existence", 
			zap.String("name", name), 
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (r *permissionRepository) GetByResource(ctx context.Context, resource string) ([]*entity.Permission, error) {
	permissions, err := r.client.Permission.
		Query().
		Where(permission.Resource(resource)).
		Order(ent.Asc(permission.FieldAction)).
		All(ctx)
	
	if err != nil {
		logger.Error("Failed to get permissions by resource", 
			zap.String("resource", resource), 
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.Permission, len(permissions))
	for i, perm := range permissions {
		result[i] = r.convertToEntity(perm)
	}
	
	return result, nil
}

// convertToEntity 将EntGo实体转换为领域实体
func (r *permissionRepository) convertToEntity(permEnt *ent.Permission) *entity.Permission {
	return &entity.Permission{
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