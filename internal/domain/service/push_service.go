package service

import (
	"context"
	"errors"

	"nebula-live/internal/domain/entity"
	"nebula-live/internal/pkg/push"
	"nebula-live/pkg/logger"

	"go.uber.org/zap"
)

// Push notification service errors
var (
	ErrPushServiceUnavailable = errors.New("push service is unavailable")
	ErrInvalidPushProvider    = errors.New("invalid push provider")
)

// PushService defines the interface for push notification service
type PushService interface {
	// SendToUserDevices sends push notifications to all enabled devices of a user
	SendToUserDevices(ctx context.Context, userID uint, message *push.PushMessage) ([]*push.PushResponse, error)
	
	// SendToUserDevicesByProvider sends push notifications to user devices of specific provider
	SendToUserDevicesByProvider(ctx context.Context, userID uint, provider string, message *push.PushMessage) ([]*push.PushResponse, error)
}

// pushService implements PushService
type pushService struct {
	userPushSettingService UserPushSettingService
}

// NewPushService creates a new push service
func NewPushService(userPushSettingService UserPushSettingService) PushService {
	return &pushService{
		userPushSettingService: userPushSettingService,
	}
}


// SendToUserDevices sends push notifications to all enabled devices of a user
func (s *pushService) SendToUserDevices(ctx context.Context, userID uint, message *push.PushMessage) ([]*push.PushResponse, error) {
	logger.Info("Sending push notification to user devices",
		zap.Uint("user_id", userID),
		zap.String("title", message.Title))

	if s.userPushSettingService == nil {
		logger.Error("Push service is not properly initialized")
		return nil, ErrPushServiceUnavailable
	}

	// 获取用户的所有启用推送设置
	userSettings, err := s.userPushSettingService.GetEnabledUserSettings(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user push settings",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	if len(userSettings) == 0 {
		logger.Info("No enabled push settings found for user",
			zap.Uint("user_id", userID))
		return []*push.PushResponse{}, nil
	}

	var responses []*push.PushResponse
	
	for _, setting := range userSettings {
		// 创建消息副本并应用用户设置
		userMessage := *message
		userMessage.DeviceID = setting.DeviceID
		
		// 应用用户特定设置
		if err := s.applyUserSettings(setting, &userMessage); err != nil {
			logger.Error("Failed to apply user settings",
				zap.Uint("user_id", userID),
				zap.Uint("setting_id", setting.ID),
				zap.Error(err))
			continue
		}

		// 基于用户设置创建推送客户端
		pushClient, err := s.createPushClientForSetting(setting)
		if err != nil {
			logger.Error("Failed to create push client for setting",
				zap.Uint("user_id", userID),
				zap.Uint("setting_id", setting.ID),
				zap.Error(err))
			continue
		}

		// 发送推送通知
		response, err := pushClient.SendMessage(ctx, setting.Provider, &userMessage)
		if err != nil {
			logger.Error("Failed to send push notification to user device",
				zap.Uint("user_id", userID),
				zap.String("provider", setting.Provider),
				zap.String("device_id", setting.DeviceID),
				zap.Error(err))
			// 创建错误响应
			response = &push.PushResponse{
				Success:  false,
				Error:    err.Error(),
				Provider: setting.Provider,
			}
		}
		
		if response != nil {
			responses = append(responses, response)
		}
	}

	logger.Info("User push notification batch completed",
		zap.Uint("user_id", userID),
		zap.Int("total_devices", len(userSettings)),
		zap.Int("responses", len(responses)))

	return responses, nil
}

// SendToUserDevicesByProvider sends push notifications to user devices of specific provider
func (s *pushService) SendToUserDevicesByProvider(ctx context.Context, userID uint, provider string, message *push.PushMessage) ([]*push.PushResponse, error) {
	logger.Info("Sending push notification to user devices by provider",
		zap.Uint("user_id", userID),
		zap.String("provider", provider),
		zap.String("title", message.Title))

	if s.userPushSettingService == nil {
		logger.Error("Push service is not properly initialized")
		return nil, ErrPushServiceUnavailable
	}

	// 获取用户指定提供商的启用推送设置
	userSettings, err := s.userPushSettingService.GetEnabledUserSettingsByProvider(ctx, userID, provider)
	if err != nil {
		logger.Error("Failed to get user push settings by provider",
			zap.Uint("user_id", userID),
			zap.String("provider", provider),
			zap.Error(err))
		return nil, err
	}

	if len(userSettings) == 0 {
		logger.Info("No enabled push settings found for user and provider",
			zap.Uint("user_id", userID),
			zap.String("provider", provider))
		return []*push.PushResponse{}, nil
	}

	var responses []*push.PushResponse

	for _, setting := range userSettings {
		// 创建消息副本并应用用户设置
		userMessage := *message
		userMessage.DeviceID = setting.DeviceID
		
		// 应用用户特定设置
		if err := s.applyUserSettings(setting, &userMessage); err != nil {
			logger.Error("Failed to apply user settings",
				zap.Uint("user_id", userID),
				zap.Uint("setting_id", setting.ID),
				zap.Error(err))
			continue
		}

		// 基于用户设置创建推送客户端
		pushClient, err := s.createPushClientForSetting(setting)
		if err != nil {
			logger.Error("Failed to create push client for setting",
				zap.Uint("user_id", userID),
				zap.Uint("setting_id", setting.ID),
				zap.Error(err))
			continue
		}

		// 发送推送通知
		response, err := pushClient.SendMessage(ctx, setting.Provider, &userMessage)
		if err != nil {
			logger.Error("Failed to send push notification to user device",
				zap.Uint("user_id", userID),
				zap.String("provider", setting.Provider),
				zap.String("device_id", setting.DeviceID),
				zap.Error(err))
			// 创建错误响应
			response = &push.PushResponse{
				Success:  false,
				Error:    err.Error(),
				Provider: setting.Provider,
			}
		}
		
		if response != nil {
			responses = append(responses, response)
		}
	}

	logger.Info("User push notification batch by provider completed",
		zap.Uint("user_id", userID),
		zap.String("provider", provider),
		zap.Int("total_devices", len(userSettings)),
		zap.Int("responses", len(responses)))

	return responses, nil
}

// createPushClientForSetting creates a push client based on user setting
func (s *pushService) createPushClientForSetting(setting *entity.UserPushSetting) (*push.Client, error) {
	switch setting.Provider {
	case "bark":
		barkSettings, err := setting.GetBarkSettings()
		if err != nil {
			return nil, err
		}
		
		// 创建Bark配置
		barkConfig := push.BarkConfig{
			BaseURL: "https://api.day.app", // 默认服务器
			Enabled: true,
		}
		
		// 如果用户设置了自定义服务器
		if barkSettings != nil && barkSettings.BaseURL != "" {
			barkConfig.BaseURL = barkSettings.BaseURL
		}
		
		clientConfig := push.ClientConfig{
			Bark: barkConfig,
		}
		
		return push.NewClient(clientConfig), nil
	default:
		return nil, errors.New("unsupported push provider: " + setting.Provider)
	}
}

// applyUserSettings applies user-specific settings to the push message
func (s *pushService) applyUserSettings(setting *entity.UserPushSetting, message *push.PushMessage) error {
	switch setting.Provider {
	case "bark":
		barkSettings, err := setting.GetBarkSettings()
		if err != nil {
			return err
		}
		if barkSettings != nil {
			// 应用用户的Bark设置，如果消息中没有指定的话
			if message.Sound == "" && barkSettings.Sound != "" {
				message.Sound = barkSettings.Sound
			}
			if message.Icon == "" && barkSettings.Icon != "" {
				message.Icon = barkSettings.Icon
			}
			if message.Group == "" && barkSettings.Group != "" {
				message.Group = barkSettings.Group
			}
			if message.Level == "" && barkSettings.Level != "" {
				message.Level = push.PushLevel(barkSettings.Level)
			}
			if !message.AutoCopy && barkSettings.AutoCopy {
				message.AutoCopy = barkSettings.AutoCopy
			}
			if !message.Call && barkSettings.Call {
				message.Call = barkSettings.Call
			}
		}
	}
	return nil
}
