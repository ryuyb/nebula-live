package persistence

import (
	"context"

	"nebula-live/ent"
	"nebula-live/ent/user"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/internal/domain/service"
)

// userRepository 用户仓储实现
type userRepository struct {
	client *ent.Client
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(client *ent.Client) repository.UserRepository {
	return &userRepository{
		client: client,
	}
}

// entUserToDomainUser 将ent.User转换为domain.User
func entUserToDomainUser(entUser *ent.User) *entity.User {
	if entUser == nil {
		return nil
	}
	
	var status entity.UserStatus
	switch entUser.Status {
	case user.StatusActive:
		status = entity.UserStatusActive
	case user.StatusInactive:
		status = entity.UserStatusInactive
	case user.StatusBanned:
		status = entity.UserStatusBanned
	default:
		status = entity.UserStatusActive
	}
	
	return &entity.User{
		ID:        entUser.ID,
		Username:  entUser.Username,
		Email:     entUser.Email,
		Password:  entUser.Password,
		Nickname:  entUser.Nickname,
		Avatar:    entUser.Avatar,
		Status:    status,
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
	}
}

// domainUserStatusToEntStatus 将domain UserStatus转换为ent status
func domainUserStatusToEntStatus(status entity.UserStatus) user.Status {
	switch status {
	case entity.UserStatusActive:
		return user.StatusActive
	case entity.UserStatusInactive:
		return user.StatusInactive
	case entity.UserStatusBanned:
		return user.StatusBanned
	default:
		return user.StatusActive
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, u *entity.User) error {
	entUser, err := r.client.User.
		Create().
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetNillableNickname(&u.Nickname).
		SetNillableAvatar(&u.Avatar).
		SetStatus(domainUserStatusToEntStatus(u.Status)).
		Save(ctx)
	if err != nil {
		return err
	}
	
	// 更新ID
	u.ID = entUser.ID
	u.CreatedAt = entUser.CreatedAt
	u.UpdatedAt = entUser.UpdatedAt
	
	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}
	
	return entUserToDomainUser(entUser), nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(user.Username(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}
	
	return entUserToDomainUser(entUser), nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}
	
	return entUserToDomainUser(entUser), nil
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, u *entity.User) error {
	_, err := r.client.User.
		UpdateOneID(u.ID).
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetNillableNickname(&u.Nickname).
		SetNillableAvatar(&u.Avatar).
		SetStatus(domainUserStatusToEntStatus(u.Status)).
		SetUpdatedAt(u.UpdatedAt).
		Save(ctx)
	return err
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	err := r.client.User.
		DeleteOneID(id).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return service.ErrUserNotFound
		}
	}
	return err
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	entUsers, err := r.client.User.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	
	users := make([]*entity.User, len(entUsers))
	for i, entUser := range entUsers {
		users[i] = entUserToDomainUser(entUser)
	}
	
	return users, nil
}

// Count 获取用户总数
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.client.User.
		Query().
		Count(ctx)
	return int64(count), err
}

// ExistsByUsername 检查用户名是否已存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.client.User.
		Query().
		Where(user.Username(username)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}