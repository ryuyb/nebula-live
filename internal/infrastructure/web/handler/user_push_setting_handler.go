package handler

import (
	"nebula-live/internal/domain/service"
	"nebula-live/internal/infrastructure/web/dto"
	"nebula-live/pkg/auth"
	apierrors "nebula-live/pkg/errors"
	"nebula-live/pkg/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserPushSettingHandler 用户推送设置处理器
type UserPushSettingHandler struct {
	userPushSettingService service.UserPushSettingService
}

// NewUserPushSettingHandler 创建用户推送设置处理器
func NewUserPushSettingHandler(userPushSettingService service.UserPushSettingService) *UserPushSettingHandler {
	return &UserPushSettingHandler{
		userPushSettingService: userPushSettingService,
	}
}

// CreateSetting 创建用户推送设置
// POST /api/v1/push-settings
func (h *UserPushSettingHandler) CreateSetting(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	var req dto.CreateUserPushSettingRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "Failed to parse request body"),
		)
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Validation failed", err.Error()),
		)
	}

	setting, err := h.userPushSettingService.CreateSetting(
		c.Context(),
		userID,
		req.Provider,
		req.DeviceID,
		req.DeviceName,
		req.Settings,
	)

	if err != nil {
		logger.Error("Failed to create user push setting", 
			zap.Uint("user_id", userID), 
			zap.Error(err))
		
		switch err {
		case service.ErrDeviceAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(
				apierrors.NewAPIError(fiber.StatusConflict, "Device already exists", "Device with this ID already registered"),
			)
		case service.ErrInvalidUserPushSetting:
			return c.Status(fiber.StatusBadRequest).JSON(
				apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid setting", "Invalid push setting configuration"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to create push setting"),
			)
		}
	}

	response := dto.UserPushSettingResponse{
		ID:         setting.ID,
		UserID:     setting.UserID,
		Provider:   setting.Provider,
		Enabled:    setting.Enabled,
		DeviceID:   setting.DeviceID,
		DeviceName: setting.DeviceName,
		Settings:   setting.Settings,
		CreatedAt:  setting.CreatedAt,
		UpdatedAt:  setting.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetSettings 获取当前用户的推送设置列表
// GET /api/v1/push-settings
func (h *UserPushSettingHandler) GetSettings(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	provider := c.Query("provider")

	var settings []dto.UserPushSettingResponse
	var total int64

	if provider != "" {
		// 获取指定提供商的设置
		userSettings, err := h.userPushSettingService.GetEnabledUserSettingsByProvider(c.Context(), userID, provider)
		if err != nil {
			logger.Error("Failed to get user push settings by provider", 
				zap.Uint("user_id", userID), 
				zap.String("provider", provider),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get push settings"),
			)
		}

		settings = make([]dto.UserPushSettingResponse, len(userSettings))
		for i, setting := range userSettings {
			settings[i] = dto.UserPushSettingResponse{
				ID:         setting.ID,
				UserID:     setting.UserID,
				Provider:   setting.Provider,
				Enabled:    setting.Enabled,
				DeviceID:   setting.DeviceID,
				DeviceName: setting.DeviceName,
				Settings:   setting.Settings,
				CreatedAt:  setting.CreatedAt,
				UpdatedAt:  setting.UpdatedAt,
			}
		}
		total = int64(len(settings))
	} else {
		// 获取分页的设置列表
		userSettings, totalCount, err := h.userPushSettingService.ListSettings(c.Context(), userID, page, limit)
		if err != nil {
			logger.Error("Failed to list user push settings", 
				zap.Uint("user_id", userID), 
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to list push settings"),
			)
		}

		settings = make([]dto.UserPushSettingResponse, len(userSettings))
		for i, setting := range userSettings {
			settings[i] = dto.UserPushSettingResponse{
				ID:         setting.ID,
				UserID:     setting.UserID,
				Provider:   setting.Provider,
				Enabled:    setting.Enabled,
				DeviceID:   setting.DeviceID,
				DeviceName: setting.DeviceName,
				Settings:   setting.Settings,
				CreatedAt:  setting.CreatedAt,
				UpdatedAt:  setting.UpdatedAt,
			}
		}
		total = totalCount
	}

	response := dto.ListResponse[dto.UserPushSettingResponse]{
		Data:  settings,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return c.JSON(response)
}

// GetSetting 获取指定的推送设置
// GET /api/v1/push-settings/:id
func (h *UserPushSettingHandler) GetSetting(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	settingID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid ID", "Invalid setting ID"),
		)
	}

	setting, err := h.userPushSettingService.GetSetting(c.Context(), userID, uint(settingID))
	if err != nil {
		logger.Error("Failed to get user push setting", 
			zap.Uint("user_id", userID), 
			zap.Uint("setting_id", uint(settingID)),
			zap.Error(err))
		
		switch err {
		case service.ErrUserPushSettingNotFound:
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Setting not found", "Push setting not found"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get push setting"),
			)
		}
	}

	response := dto.UserPushSettingResponse{
		ID:         setting.ID,
		UserID:     setting.UserID,
		Provider:   setting.Provider,
		Enabled:    setting.Enabled,
		DeviceID:   setting.DeviceID,
		DeviceName: setting.DeviceName,
		Settings:   setting.Settings,
		CreatedAt:  setting.CreatedAt,
		UpdatedAt:  setting.UpdatedAt,
	}

	return c.JSON(response)
}

// UpdateSetting 更新推送设置
// PUT /api/v1/push-settings/:id
func (h *UserPushSettingHandler) UpdateSetting(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	settingID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid ID", "Invalid setting ID"),
		)
	}

	var req dto.UpdateUserPushSettingRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "Failed to parse request body"),
		)
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Validation failed", err.Error()),
		)
	}

	// 获取现有设置
	existingSetting, err := h.userPushSettingService.GetSetting(c.Context(), userID, uint(settingID))
	if err != nil {
		switch err {
		case service.ErrUserPushSettingNotFound:
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Setting not found", "Push setting not found"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get push setting"),
			)
		}
	}

	// 更新字段
	if req.Enabled != nil {
		existingSetting.Enabled = *req.Enabled
	}
	if req.DeviceName != nil {
		existingSetting.DeviceName = *req.DeviceName
	}
	if req.Settings != nil {
		existingSetting.Settings = req.Settings
	}

	setting, err := h.userPushSettingService.UpdateSetting(c.Context(), userID, existingSetting)
	if err != nil {
		logger.Error("Failed to update user push setting", 
			zap.Uint("user_id", userID), 
			zap.Uint("setting_id", uint(settingID)),
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(
			apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to update push setting"),
		)
	}

	response := dto.UserPushSettingResponse{
		ID:         setting.ID,
		UserID:     setting.UserID,
		Provider:   setting.Provider,
		Enabled:    setting.Enabled,
		DeviceID:   setting.DeviceID,
		DeviceName: setting.DeviceName,
		Settings:   setting.Settings,
		CreatedAt:  setting.CreatedAt,
		UpdatedAt:  setting.UpdatedAt,
	}

	return c.JSON(response)
}

// EnableSetting 启用推送设置
// POST /api/v1/push-settings/:id/enable
func (h *UserPushSettingHandler) EnableSetting(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	settingID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid ID", "Invalid setting ID"),
		)
	}

	err = h.userPushSettingService.EnableSetting(c.Context(), userID, uint(settingID))
	if err != nil {
		logger.Error("Failed to enable user push setting", 
			zap.Uint("user_id", userID), 
			zap.Uint("setting_id", uint(settingID)),
			zap.Error(err))
		
		switch err {
		case service.ErrUserPushSettingNotFound:
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Setting not found", "Push setting not found"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to enable push setting"),
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Push setting enabled successfully",
	})
}

// DisableSetting 禁用推送设置
// POST /api/v1/push-settings/:id/disable
func (h *UserPushSettingHandler) DisableSetting(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	settingID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid ID", "Invalid setting ID"),
		)
	}

	err = h.userPushSettingService.DisableSetting(c.Context(), userID, uint(settingID))
	if err != nil {
		logger.Error("Failed to disable user push setting", 
			zap.Uint("user_id", userID), 
			zap.Uint("setting_id", uint(settingID)),
			zap.Error(err))
		
		switch err {
		case service.ErrUserPushSettingNotFound:
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Setting not found", "Push setting not found"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to disable push setting"),
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Push setting disabled successfully",
	})
}

// DeleteSetting 删除推送设置
// DELETE /api/v1/push-settings/:id
func (h *UserPushSettingHandler) DeleteSetting(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	settingID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid ID", "Invalid setting ID"),
		)
	}

	err = h.userPushSettingService.DeleteSetting(c.Context(), userID, uint(settingID))
	if err != nil {
		logger.Error("Failed to delete user push setting", 
			zap.Uint("user_id", userID), 
			zap.Uint("setting_id", uint(settingID)),
			zap.Error(err))
		
		switch err {
		case service.ErrUserPushSettingNotFound:
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Setting not found", "Push setting not found"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to delete push setting"),
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Push setting deleted successfully",
	})
}

// GetSupportedProviders 获取支持的推送提供商列表
// GET /api/v1/push-settings/providers
func (h *UserPushSettingHandler) GetSupportedProviders(c *fiber.Ctx) error {
	// 返回支持的推送提供商列表
	providers := []fiber.Map{
		{
			"name":        "bark",
			"display_name": "Bark",
			"description":  "iOS Bark push notification service",
			"platform":     "ios",
			"settings": fiber.Map{
				"base_url":  "Custom Bark server URL (optional)",
				"sound":     "Notification sound (optional)",
				"icon":      "Notification icon URL (optional)", 
				"group":     "Notification group (optional)",
				"level":     "Notification level: active, critical, timeSensitive, passive (optional)",
				"auto_copy": "Auto copy message to clipboard (optional)",
				"call":      "Ring for 30 seconds (optional)",
			},
		},
	}

	return c.JSON(fiber.Map{
		"providers": providers,
		"total":     len(providers),
	})
}

// ValidateDevice 验证设备ID是否可用
// POST /api/v1/push-settings/validate-device
func (h *UserPushSettingHandler) ValidateDevice(c *fiber.Ctx) error {
	var req dto.ValidateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "Failed to parse request body"),
		)
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Validation failed", err.Error()),
		)
	}

	err := h.userPushSettingService.ValidateDeviceID(c.Context(), req.Provider, req.DeviceID)
	if err != nil {
		switch err {
		case service.ErrDeviceAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(
				apierrors.NewAPIError(fiber.StatusConflict, "Device already exists", "Device with this ID is already registered"),
			)
		default:
			logger.Error("Failed to validate device ID", 
				zap.String("provider", req.Provider),
				zap.String("device_id", req.DeviceID),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to validate device"),
			)
		}
	}

	return c.JSON(fiber.Map{
		"valid": true,
		"message": "Device ID is available",
	})
}