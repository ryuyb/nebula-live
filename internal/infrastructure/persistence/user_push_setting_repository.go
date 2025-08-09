package persistence

import (
	"context"
	"nebula-live/ent"
	"nebula-live/ent/userpushsetting"
	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

type userPushSettingRepository struct {
	client *ent.Client
}

// NewUserPushSettingRepository 创建用户推送设置仓储实例
func NewUserPushSettingRepository(client *ent.Client) repository.UserPushSettingRepository {
	return &userPushSettingRepository{
		client: client,
	}
}

// convertToEntity 转换EntGo实体到Domain实体
func (r *userPushSettingRepository) convertToEntity(entSetting *ent.UserPushSetting) *entity.UserPushSetting {
	return &entity.UserPushSetting{
		ID:         entSetting.ID,
		UserID:     entSetting.UserID,
		Provider:   entSetting.Provider.String(),
		Enabled:    entSetting.Enabled,
		DeviceID:   entSetting.DeviceID,
		DeviceName: entSetting.DeviceName,
		Settings:   entSetting.Settings,
		CreatedAt:  entSetting.CreatedAt,
		UpdatedAt:  entSetting.UpdatedAt,
	}
}

// Create 创建用户推送设置
func (r *userPushSettingRepository) Create(ctx context.Context, setting *entity.UserPushSetting) (*entity.UserPushSetting, error) {
	entSetting, err := r.client.UserPushSetting.
		Create().
		SetUserID(setting.UserID).
		SetProvider(userpushsetting.Provider(setting.Provider)).
		SetEnabled(setting.Enabled).
		SetDeviceID(setting.DeviceID).
		SetNillableDeviceName(&setting.DeviceName).
		SetSettings(setting.Settings).
		Save(ctx)

	if err != nil {
		logger.Error("Failed to create user push setting",
			zap.Uint("user_id", setting.UserID),
			zap.String("provider", setting.Provider),
			zap.Error(err))
		return nil, err
	}

	logger.Info("User push setting created successfully",
		zap.Uint("id", entSetting.ID),
		zap.Uint("user_id", setting.UserID),
		zap.String("provider", setting.Provider))

	return r.convertToEntity(entSetting), nil
}

// GetByID 根据ID获取用户推送设置
func (r *userPushSettingRepository) GetByID(ctx context.Context, id uint) (*entity.UserPushSetting, error) {
	entSetting, err := r.client.UserPushSetting.
		Query().
		Where(userpushsetting.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		logger.Error("Failed to get user push setting by ID",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	return r.convertToEntity(entSetting), nil
}

// GetByUserIDAndProvider 根据用户ID和提供商获取推送设置
func (r *userPushSettingRepository) GetByUserIDAndProvider(ctx context.Context, userID uint, provider string) ([]*entity.UserPushSetting, error) {
	entSettings, err := r.client.UserPushSetting.
		Query().
		Where(
			userpushsetting.UserID(userID),
			userpushsetting.ProviderEQ(userpushsetting.Provider(provider)),
		).
		Order(ent.Desc(userpushsetting.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to get user push settings by user ID and provider",
			zap.Uint("user_id", userID),
			zap.String("provider", provider),
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserPushSetting, len(entSettings))
	for i, entSetting := range entSettings {
		result[i] = r.convertToEntity(entSetting)
	}

	return result, nil
}

// GetByUserID 获取用户的所有推送设置
func (r *userPushSettingRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error) {
	entSettings, err := r.client.UserPushSetting.
		Query().
		Where(userpushsetting.UserID(userID)).
		Order(ent.Desc(userpushsetting.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to get user push settings",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserPushSetting, len(entSettings))
	for i, entSetting := range entSettings {
		result[i] = r.convertToEntity(entSetting)
	}

	return result, nil
}

// GetEnabledByUserID 获取用户的所有启用的推送设置
func (r *userPushSettingRepository) GetEnabledByUserID(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error) {
	entSettings, err := r.client.UserPushSetting.
		Query().
		Where(
			userpushsetting.UserID(userID),
			userpushsetting.EnabledEQ(true),
		).
		Order(ent.Desc(userpushsetting.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to get enabled user push settings",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserPushSetting, len(entSettings))
	for i, entSetting := range entSettings {
		result[i] = r.convertToEntity(entSetting)
	}

	return result, nil
}

// GetEnabledByUserIDAndProvider 获取用户在指定提供商的启用推送设置
func (r *userPushSettingRepository) GetEnabledByUserIDAndProvider(ctx context.Context, userID uint, provider string) ([]*entity.UserPushSetting, error) {
	entSettings, err := r.client.UserPushSetting.
		Query().
		Where(
			userpushsetting.UserID(userID),
			userpushsetting.ProviderEQ(userpushsetting.Provider(provider)),
			userpushsetting.EnabledEQ(true),
		).
		Order(ent.Desc(userpushsetting.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to get enabled user push settings by provider",
			zap.Uint("user_id", userID),
			zap.String("provider", provider),
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserPushSetting, len(entSettings))
	for i, entSetting := range entSettings {
		result[i] = r.convertToEntity(entSetting)
	}

	return result, nil
}

// Update 更新用户推送设置
func (r *userPushSettingRepository) Update(ctx context.Context, setting *entity.UserPushSetting) (*entity.UserPushSetting, error) {
	entSetting, err := r.client.UserPushSetting.
		UpdateOneID(setting.ID).
		SetEnabled(setting.Enabled).
		SetNillableDeviceName(&setting.DeviceName).
		SetSettings(setting.Settings).
		Save(ctx)

	if err != nil {
		logger.Error("Failed to update user push setting",
			zap.Uint("id", setting.ID),
			zap.Error(err))
		return nil, err
	}

	logger.Info("User push setting updated successfully",
		zap.Uint("id", setting.ID),
		zap.Uint("user_id", setting.UserID))

	return r.convertToEntity(entSetting), nil
}

// Delete 删除用户推送设置
func (r *userPushSettingRepository) Delete(ctx context.Context, id uint) error {
	err := r.client.UserPushSetting.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to delete user push setting",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	logger.Info("User push setting deleted successfully",
		zap.Uint("id", id))

	return nil
}

// DeleteByUserIDAndDeviceID 根据用户ID和设备ID删除推送设置
func (r *userPushSettingRepository) DeleteByUserIDAndDeviceID(ctx context.Context, userID uint, provider, deviceID string) error {
	_, err := r.client.UserPushSetting.
		Delete().
		Where(
			userpushsetting.UserID(userID),
			userpushsetting.ProviderEQ(userpushsetting.Provider(provider)),
			userpushsetting.DeviceIDEQ(deviceID),
		).
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to delete user push setting by device ID",
			zap.Uint("user_id", userID),
			zap.String("provider", provider),
			zap.String("device_id", deviceID),
			zap.Error(err))
		return err
	}

	logger.Info("User push setting deleted by device ID",
		zap.Uint("user_id", userID),
		zap.String("provider", provider),
		zap.String("device_id", deviceID))

	return nil
}

// ExistsByProviderAndDeviceID 检查设备是否已存在
func (r *userPushSettingRepository) ExistsByProviderAndDeviceID(ctx context.Context, provider, deviceID string) (bool, error) {
	exists, err := r.client.UserPushSetting.
		Query().
		Where(
			userpushsetting.ProviderEQ(userpushsetting.Provider(provider)),
			userpushsetting.DeviceIDEQ(deviceID),
		).
		Exist(ctx)

	if err != nil {
		logger.Error("Failed to check user push setting existence",
			zap.String("provider", provider),
			zap.String("device_id", deviceID),
			zap.Error(err))
		return false, err
	}

	return exists, nil
}

// List 获取用户推送设置列表（带分页）
func (r *userPushSettingRepository) List(ctx context.Context, userID uint, offset, limit int) ([]*entity.UserPushSetting, error) {
	entSettings, err := r.client.UserPushSetting.
		Query().
		Where(userpushsetting.UserID(userID)).
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(userpushsetting.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Error("Failed to list user push settings",
			zap.Uint("user_id", userID),
			zap.Int("offset", offset),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserPushSetting, len(entSettings))
	for i, entSetting := range entSettings {
		result[i] = r.convertToEntity(entSetting)
	}

	return result, nil
}

// Count 获取用户推送设置总数
func (r *userPushSettingRepository) Count(ctx context.Context, userID uint) (int64, error) {
	count, err := r.client.UserPushSetting.
		Query().
		Where(userpushsetting.UserID(userID)).
		Count(ctx)

	if err != nil {
		logger.Error("Failed to count user push settings",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return 0, err
	}

	return int64(count), nil
}