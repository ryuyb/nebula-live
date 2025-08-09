package repository

import (
	"context"
	"nebula-live/internal/domain/entity"
)

// UserPushSettingRepository 用户推送设置仓储接口
type UserPushSettingRepository interface {
	// Create 创建用户推送设置
	Create(ctx context.Context, setting *entity.UserPushSetting) (*entity.UserPushSetting, error)
	
	// GetByID 根据ID获取用户推送设置
	GetByID(ctx context.Context, id uint) (*entity.UserPushSetting, error)
	
	// GetByUserIDAndProvider 根据用户ID和提供商获取推送设置
	GetByUserIDAndProvider(ctx context.Context, userID uint, provider string) ([]*entity.UserPushSetting, error)
	
	// GetByUserID 获取用户的所有推送设置
	GetByUserID(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error)
	
	// GetEnabledByUserID 获取用户的所有启用的推送设置
	GetEnabledByUserID(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error)
	
	// GetEnabledByUserIDAndProvider 获取用户在指定提供商的启用推送设置
	GetEnabledByUserIDAndProvider(ctx context.Context, userID uint, provider string) ([]*entity.UserPushSetting, error)
	
	// Update 更新用户推送设置
	Update(ctx context.Context, setting *entity.UserPushSetting) (*entity.UserPushSetting, error)
	
	// Delete 删除用户推送设置
	Delete(ctx context.Context, id uint) error
	
	// DeleteByUserIDAndDeviceID 根据用户ID和设备ID删除推送设置
	DeleteByUserIDAndDeviceID(ctx context.Context, userID uint, provider, deviceID string) error
	
	// ExistsByProviderAndDeviceID 检查设备是否已存在
	ExistsByProviderAndDeviceID(ctx context.Context, provider, deviceID string) (bool, error)
	
	// List 获取用户推送设置列表（带分页）
	List(ctx context.Context, userID uint, offset, limit int) ([]*entity.UserPushSetting, error)
	
	// Count 获取用户推送设置总数
	Count(ctx context.Context, userID uint) (int64, error)
}