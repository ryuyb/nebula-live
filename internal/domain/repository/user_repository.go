package repository

import (
	"context"

	"nebula-live/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error
	
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	
	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// Update 更新用户信息
	Update(ctx context.Context, user *entity.User) error
	
	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
	
	// List 获取用户列表
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
	
	// Count 获取用户总数
	Count(ctx context.Context) (int64, error)
	
	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}