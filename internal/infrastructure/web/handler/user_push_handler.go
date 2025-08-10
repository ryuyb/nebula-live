package handler

import (
	"nebula-live/internal/domain/service"
	"nebula-live/internal/infrastructure/web/dto"
	"nebula-live/internal/pkg/push"
	"nebula-live/pkg/auth"
	apierrors "nebula-live/pkg/errors"
	"nebula-live/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserPushHandler 用户推送处理器
type UserPushHandler struct {
	pushService service.PushService
}

// NewUserPushHandler 创建用户推送处理器
func NewUserPushHandler(pushService service.PushService) *UserPushHandler {
	return &UserPushHandler{
		pushService: pushService,
	}
}

// SendToMyDevices godoc
// @Summary      Send Push to My Devices
// @Description  Send push notification to current user's all enabled devices
// @Tags         Push Notifications
// @Accept       json
// @Produce      json
// @Param        notification body dto.UserPushRequest true "Push notification data"
// @Success      200 {object} dto.UserPushResult "Push notification sent successfully"
// @Failure      400 {object} errors.APIError "Invalid request parameters or validation failed"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /push/my-devices [post]
func (h *UserPushHandler) SendToMyDevices(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	var req dto.UserPushRequest
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

	// 创建推送消息
	message := &push.PushMessage{
		Title:    req.Title,
		Body:     req.Body,
		URL:      req.URL,
		Sound:    req.Sound,
		Icon:     req.Icon,
		Group:    req.Group,
		Level:    push.PushLevel(req.Level),
		AutoCopy: req.AutoCopy,
		Call:     req.Call,
	}

	// 发送到用户的所有设备
	responses, err := h.pushService.SendToUserDevices(c.Context(), userID, message)
	if err != nil {
		logger.Error("Failed to send push notification to user devices", 
			zap.Uint("user_id", userID), 
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(
			apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to send push notifications"),
		)
	}

	// 转换响应
	responseData := make([]dto.PushResponse, len(responses))
	successCount := 0
	
	for i, resp := range responses {
		responseData[i] = dto.PushResponse{
			Success:   resp.Success,
			MessageID: resp.MessageID,
			Provider:  resp.Provider,
			Error:     resp.Error,
		}
		if resp.Success {
			successCount++
		}
	}

	result := dto.UserPushResult{
		UserID:       userID,
		TotalDevices: len(responses),
		SuccessCount: successCount,
		FailedCount:  len(responses) - successCount,
		Responses:    responseData,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// SendToMyDevicesByProvider godoc
// @Summary      Send Push to My Devices by Provider
// @Description  Send push notification to current user's devices for specific provider
// @Tags         Push Notifications
// @Accept       json
// @Produce      json
// @Param        provider path string true "Push provider name" Enums(bark) example(bark)
// @Param        notification body dto.UserPushRequest true "Push notification data"
// @Success      200 {object} dto.UserPushResult "Push notification sent successfully"
// @Failure      400 {object} errors.APIError "Invalid request parameters or validation failed"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /push/my-devices/{provider} [post]
func (h *UserPushHandler) SendToMyDevicesByProvider(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	provider := c.Params("provider")
	if provider == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid provider", "Provider is required"),
		)
	}

	var req dto.UserPushRequest
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

	// 创建推送消息
	message := &push.PushMessage{
		Title:    req.Title,
		Body:     req.Body,
		URL:      req.URL,
		Sound:    req.Sound,
		Icon:     req.Icon,
		Group:    req.Group,
		Level:    push.PushLevel(req.Level),
		AutoCopy: req.AutoCopy,
		Call:     req.Call,
	}

	// 发送到用户指定提供商的设备
	responses, err := h.pushService.SendToUserDevicesByProvider(c.Context(), userID, provider, message)
	if err != nil {
		logger.Error("Failed to send push notification to user devices by provider", 
			zap.Uint("user_id", userID), 
			zap.String("provider", provider),
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(
			apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to send push notifications"),
		)
	}

	// 转换响应
	responseData := make([]dto.PushResponse, len(responses))
	successCount := 0
	
	for i, resp := range responses {
		responseData[i] = dto.PushResponse{
			Success:   resp.Success,
			MessageID: resp.MessageID,
			Provider:  resp.Provider,
			Error:     resp.Error,
		}
		if resp.Success {
			successCount++
		}
	}

	result := dto.UserPushResult{
		UserID:       userID,
		Provider:     provider,
		TotalDevices: len(responses),
		SuccessCount: successCount,
		FailedCount:  len(responses) - successCount,
		Responses:    responseData,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// TestMyPushSettings godoc
// @Summary      Test Push Settings
// @Description  Send a test push notification to current user's all enabled devices
// @Tags         Push Notifications
// @Accept       json
// @Produce      json
// @Success      200 {object} dto.UserPushResult "Test notification sent successfully"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /push/test [post]
func (h *UserPushHandler) TestMyPushSettings(c *fiber.Ctx) error {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(
			apierrors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "User not authenticated"),
		)
	}

	// 创建测试消息
	message := &push.PushMessage{
		Title: "推送测试",
		Body:  "这是一条测试消息，用于验证您的推送设置是否正常工作。",
	}

	// 发送到用户的所有设备
	responses, err := h.pushService.SendToUserDevices(c.Context(), userID, message)
	if err != nil {
		logger.Error("Failed to send test push notification", 
			zap.Uint("user_id", userID), 
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(
			apierrors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to send test notification"),
		)
	}

	// 转换响应
	responseData := make([]dto.PushResponse, len(responses))
	successCount := 0
	
	for i, resp := range responses {
		responseData[i] = dto.PushResponse{
			Success:   resp.Success,
			MessageID: resp.MessageID,
			Provider:  resp.Provider,
			Error:     resp.Error,
		}
		if resp.Success {
			successCount++
		}
	}

	result := dto.UserPushResult{
		UserID:       userID,
		TotalDevices: len(responses),
		SuccessCount: successCount,
		FailedCount:  len(responses) - successCount,
		Responses:    responseData,
		Message:      "Test notification sent",
	}

	return c.Status(fiber.StatusOK).JSON(result)
}