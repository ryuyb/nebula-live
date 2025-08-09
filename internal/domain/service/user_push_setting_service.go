package service

import (
	"context"
	"errors"

	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

// 用户推送设置服务相关错误
var (
	ErrUserPushSettingNotFound     = errors.New("user push setting not found")
	ErrUserPushSettingExists       = errors.New("user push setting already exists")
	ErrInvalidUserPushSetting      = errors.New("invalid user push setting")
	ErrDeviceAlreadyExists         = errors.New("device already exists")
	ErrUserPushSettingUnavailable  = errors.New("user push setting service unavailable")
)

// UserPushSettingService 用户推送设置服务接口
type UserPushSettingService interface {
	// CreateSetting 创建用户推送设置
	CreateSetting(ctx context.Context, userID uint, provider, deviceID, deviceName string, settings map[string]interface{}) (*entity.UserPushSetting, error)
	
	// GetSetting 获取用户推送设置
	GetSetting(ctx context.Context, userID, settingID uint) (*entity.UserPushSetting, error)
	
	// GetUserSettings 获取用户的所有推送设置
	GetUserSettings(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error)
	
	// GetEnabledUserSettings 获取用户的所有启用推送设置
	GetEnabledUserSettings(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error)
	
	// GetEnabledUserSettingsByProvider 获取用户指定提供商的启用推送设置
	GetEnabledUserSettingsByProvider(ctx context.Context, userID uint, provider string) ([]*entity.UserPushSetting, error)
	
	// UpdateSetting 更新用户推送设置
	UpdateSetting(ctx context.Context, userID uint, setting *entity.UserPushSetting) (*entity.UserPushSetting, error)
	
	// EnableSetting 启用推送设置
	EnableSetting(ctx context.Context, userID, settingID uint) error
	
	// DisableSetting 禁用推送设置
	DisableSetting(ctx context.Context, userID, settingID uint) error
	
	// DeleteSetting 删除推送设置
	DeleteSetting(ctx context.Context, userID, settingID uint) error
	
	// DeleteByDeviceID 根据设备ID删除推送设置
	DeleteByDeviceID(ctx context.Context, userID uint, provider, deviceID string) error
	
	// ListSettings 获取用户推送设置列表（带分页）
	ListSettings(ctx context.Context, userID uint, page, limit int) ([]*entity.UserPushSetting, int64, error)
	
	// ValidateDeviceID 验证设备ID是否可用
	ValidateDeviceID(ctx context.Context, provider, deviceID string) error
}

// userPushSettingService 实现用户推送设置服务
type userPushSettingService struct {
	userPushSettingRepo repository.UserPushSettingRepository
	userRepo            repository.UserRepository
}

// NewUserPushSettingService 创建用户推送设置服务
func NewUserPushSettingService(
	userPushSettingRepo repository.UserPushSettingRepository,
	userRepo repository.UserRepository,
) UserPushSettingService {
	return &userPushSettingService{
		userPushSettingRepo: userPushSettingRepo,
		userRepo:            userRepo,
	}
}

// CreateSetting 创建用户推送设置
func (s *userPushSettingService) CreateSetting(ctx context.Context, userID uint, provider, deviceID, deviceName string, settings map[string]interface{}) (*entity.UserPushSetting, error) {
	logger.Info("Creating user push setting",
		zap.Uint("user_id", userID),
		zap.String("provider", provider),
		zap.String("device_id", deviceID))

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 检查设备ID是否已存在
	exists, err := s.userPushSettingRepo.ExistsByProviderAndDeviceID(ctx, provider, deviceID)
	if err != nil {
		logger.Error("Failed to check device existence",
			zap.String("provider", provider),
			zap.String("device_id", deviceID),
			zap.Error(err))
		return nil, err
	}
	if exists {
		logger.Warn("Device already exists",
			zap.String("provider", provider),
			zap.String("device_id", deviceID))
		return nil, ErrDeviceAlreadyExists
	}

	// 创建推送设置
	setting := &entity.UserPushSetting{
		UserID:     userID,
		Provider:   provider,
		Enabled:    true, // 默认启用
		DeviceID:   deviceID,
		DeviceName: deviceName,
		Settings:   settings,
	}

	if !setting.IsValid() {
		return nil, ErrInvalidUserPushSetting
	}

	createdSetting, err := s.userPushSettingRepo.Create(ctx, setting)
	if err != nil {
		logger.Error("Failed to create user push setting",
			zap.Uint("user_id", userID),
			zap.String("provider", provider),
			zap.Error(err))
		return nil, err
	}

	logger.Info("User push setting created successfully",
		zap.Uint("id", createdSetting.ID),
		zap.Uint("user_id", userID),
		zap.String("provider", provider))

	return createdSetting, nil
}

// GetSetting 获取用户推送设置
func (s *userPushSettingService) GetSetting(ctx context.Context, userID, settingID uint) (*entity.UserPushSetting, error) {
	setting, err := s.userPushSettingRepo.GetByID(ctx, settingID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, ErrUserPushSettingNotFound
	}

	// 检查设置是否属于该用户
	if setting.UserID != userID {
		return nil, ErrUserPushSettingNotFound
	}

	return setting, nil
}

// GetUserSettings 获取用户的所有推送设置
func (s *userPushSettingService) GetUserSettings(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error) {
	return s.userPushSettingRepo.GetByUserID(ctx, userID)
}

// GetEnabledUserSettings 获取用户的所有启用推送设置
func (s *userPushSettingService) GetEnabledUserSettings(ctx context.Context, userID uint) ([]*entity.UserPushSetting, error) {
	return s.userPushSettingRepo.GetEnabledByUserID(ctx, userID)
}

// GetEnabledUserSettingsByProvider 获取用户指定提供商的启用推送设置
func (s *userPushSettingService) GetEnabledUserSettingsByProvider(ctx context.Context, userID uint, provider string) ([]*entity.UserPushSetting, error) {
	return s.userPushSettingRepo.GetEnabledByUserIDAndProvider(ctx, userID, provider)
}

// UpdateSetting 更新用户推送设置
func (s *userPushSettingService) UpdateSetting(ctx context.Context, userID uint, setting *entity.UserPushSetting) (*entity.UserPushSetting, error) {
	// 检查设置是否属于该用户
	existingSetting, err := s.GetSetting(ctx, userID, setting.ID)
	if err != nil {
		return nil, err
	}
	if existingSetting == nil {
		return nil, ErrUserPushSettingNotFound
	}

	// 更新设置
	return s.userPushSettingRepo.Update(ctx, setting)
}

// EnableSetting 启用推送设置
func (s *userPushSettingService) EnableSetting(ctx context.Context, userID, settingID uint) error {
	setting, err := s.GetSetting(ctx, userID, settingID)
	if err != nil {
		return err
	}

	setting.Enable()
	_, err = s.userPushSettingRepo.Update(ctx, setting)
	return err
}

// DisableSetting 禁用推送设置
func (s *userPushSettingService) DisableSetting(ctx context.Context, userID, settingID uint) error {
	setting, err := s.GetSetting(ctx, userID, settingID)
	if err != nil {
		return err
	}

	setting.Disable()
	_, err = s.userPushSettingRepo.Update(ctx, setting)
	return err
}

// DeleteSetting 删除推送设置
func (s *userPushSettingService) DeleteSetting(ctx context.Context, userID, settingID uint) error {
	// 验证设置属于该用户
	_, err := s.GetSetting(ctx, userID, settingID)
	if err != nil {
		return err
	}

	return s.userPushSettingRepo.Delete(ctx, settingID)
}

// DeleteByDeviceID 根据设备ID删除推送设置
func (s *userPushSettingService) DeleteByDeviceID(ctx context.Context, userID uint, provider, deviceID string) error {
	return s.userPushSettingRepo.DeleteByUserIDAndDeviceID(ctx, userID, provider, deviceID)
}

// ListSettings 获取用户推送设置列表（带分页）
func (s *userPushSettingService) ListSettings(ctx context.Context, userID uint, page, limit int) ([]*entity.UserPushSetting, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	settings, err := s.userPushSettingRepo.List(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userPushSettingRepo.Count(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return settings, total, nil
}

// ValidateDeviceID 验证设备ID是否可用
func (s *userPushSettingService) ValidateDeviceID(ctx context.Context, provider, deviceID string) error {
	exists, err := s.userPushSettingRepo.ExistsByProviderAndDeviceID(ctx, provider, deviceID)
	if err != nil {
		return err
	}
	if exists {
		return ErrDeviceAlreadyExists
	}
	return nil
}