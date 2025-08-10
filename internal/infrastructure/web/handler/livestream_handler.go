package handler

import (
	"context"
	"errors"

	"nebula-live/internal/domain/service"
	"nebula-live/internal/pkg/livestream"
	apierrors "nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LiveStreamHandler struct {
	liveStreamService service.LiveStreamService
	logger            *zap.Logger
}

type GetStreamStatusRequest struct {
	Platform string `json:"platform" validate:"required"`
	RoomID   string `json:"room_id" validate:"required"`
}

type StreamStatusResponse struct {
	Platform string `json:"platform" example:"douyu"`
	RoomID   string `json:"room_id" example:"534740"`
	Status   string `json:"status" example:"online"`
}

type SupportedPlatformsResponse struct {
	Platforms []string `json:"platforms" example:"douyu,bilibili"`
}

type RoomInfoResponse struct {
	Platform      string `json:"platform" example:"douyu"`
	RoomID        string `json:"room_id" example:"534740"`
	Status        string `json:"status" example:"online"`
	Title         string `json:"title,omitempty" example:"【六神】游戏室"`
	Description   string `json:"description,omitempty" example:"欢迎来到直播间"`
	Cover         string `json:"cover,omitempty" example:"https://example.com/cover.jpg"`
	Keyframe      string `json:"keyframe,omitempty" example:"https://example.com/keyframe.jpg"`
	OwnerID       string `json:"owner_id,omitempty" example:"28206057"`
	OwnerName     string `json:"owner_name,omitempty" example:"丨马老六丨"`
	OwnerAvatar   string `json:"owner_avatar,omitempty" example:"https://example.com/avatar.jpg"`
	LiveStartTime int64  `json:"live_start_time,omitempty" example:"1609459200"`
	ViewerCount   int64  `json:"viewer_count,omitempty" example:"1234"`
	Category      string `json:"category,omitempty" example:"第五人格"`
}

func NewLiveStreamHandler(liveStreamService service.LiveStreamService, logger *zap.Logger) *LiveStreamHandler {
	return &LiveStreamHandler{
		liveStreamService: liveStreamService,
		logger:            logger,
	}
}

// GetStreamStatus godoc
// @Summary      Get Live Stream Status
// @Description  Get the current status of a live stream room on a specific platform
// @Tags         Live Streaming
// @Accept       json
// @Produce      json
// @Param        platform path string true "Streaming platform" Enums(douyu, bilibili) example(douyu)
// @Param        roomId path string true "Room ID" example(534740)
// @Success      200 {object} StreamStatusResponse "Stream status retrieved successfully"
// @Failure      400 {object} errors.APIError "Invalid request parameters"
// @Failure      404 {object} errors.APIError "Room not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Router       /live-streams/{platform}/rooms/{roomId}/status [get]
func (h *LiveStreamHandler) GetStreamStatus(c *fiber.Ctx) error {
	platform := c.Params("platform")
	roomID := c.Params("roomId")

	if platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "platform is required"),
		)
	}

	if roomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "room_id is required"),
		)
	}

	streamInfo, err := h.liveStreamService.GetStreamStatus(context.Background(), platform, roomID)
	if err != nil {
		h.logger.Error("Failed to get live stream status",
			zap.String("platform", platform),
			zap.String("room_id", roomID),
			zap.Error(err))

		// Handle specific error types
		switch {
		case errors.Is(err, livestream.ErrRoomNotFound):
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Room not found", "The specified live room does not exist"),
			)
		case errors.Is(err, livestream.ErrPlatformNotFound):
			return c.Status(fiber.StatusBadRequest).JSON(
				apierrors.NewAPIError(fiber.StatusBadRequest, "Unsupported platform", "The specified platform is not supported"),
			)
		case errors.Is(err, livestream.ErrInvalidRoomID):
			return c.Status(fiber.StatusBadRequest).JSON(
				apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid room ID", "The provided room ID is invalid"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Failed to get stream status", err.Error()),
			)
		}
	}

	// Create structured response using the defined type
	response := StreamStatusResponse{
		Platform: streamInfo.Platform,
		RoomID:   streamInfo.RoomID,
		Status:   string(streamInfo.Status),
	}

	return c.JSON(response)
}

// GetSupportedPlatforms godoc
// @Summary      Get Supported Streaming Platforms
// @Description  Get a list of all supported live streaming platforms
// @Tags         Live Streaming
// @Accept       json
// @Produce      json
// @Success      200 {object} SupportedPlatformsResponse "List of supported platforms"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Router       /live-streams/platforms [get]
func (h *LiveStreamHandler) GetSupportedPlatforms(c *fiber.Ctx) error {
	platforms := h.liveStreamService.GetSupportedPlatforms()

	// Create structured response using the defined type
	response := SupportedPlatformsResponse{
		Platforms: platforms,
	}

	return c.JSON(response)
}

// GetRoomInfo godoc
// @Summary      Get Live Room Information
// @Description  Get detailed information about a live stream room including title, owner, viewer count, etc.
// @Tags         Live Streaming
// @Accept       json
// @Produce      json
// @Param        platform path string true "Streaming platform" Enums(douyu, bilibili) example(douyu)
// @Param        roomId path string true "Room ID" example(534740)
// @Success      200 {object} RoomInfoResponse "Room information retrieved successfully"
// @Failure      400 {object} errors.APIError "Invalid request parameters"
// @Failure      404 {object} errors.APIError "Room not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Router       /live-streams/{platform}/rooms/{roomId}/info [get]
func (h *LiveStreamHandler) GetRoomInfo(c *fiber.Ctx) error {
	platform := c.Params("platform")
	roomID := c.Params("roomId")

	if platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "platform is required"),
		)
	}

	if roomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "room_id is required"),
		)
	}

	roomInfo, err := h.liveStreamService.GetRoomInfo(context.Background(), platform, roomID)
	if err != nil {
		h.logger.Error("Failed to get room info",
			zap.String("platform", platform),
			zap.String("room_id", roomID),
			zap.Error(err))

		// Handle specific error types
		switch {
		case errors.Is(err, livestream.ErrRoomNotFound):
			return c.Status(fiber.StatusNotFound).JSON(
				apierrors.NewAPIError(fiber.StatusNotFound, "Room not found", "The specified live room does not exist"),
			)
		case errors.Is(err, livestream.ErrPlatformNotFound):
			return c.Status(fiber.StatusBadRequest).JSON(
				apierrors.NewAPIError(fiber.StatusBadRequest, "Unsupported platform", "The specified platform is not supported"),
			)
		case errors.Is(err, livestream.ErrInvalidRoomID):
			return c.Status(fiber.StatusBadRequest).JSON(
				apierrors.NewAPIError(fiber.StatusBadRequest, "Invalid room ID", "The provided room ID is invalid"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				apierrors.NewAPIError(fiber.StatusInternalServerError, "Failed to get room info", err.Error()),
			)
		}
	}

	// Create structured response using the defined type
	response := RoomInfoResponse{
		Platform:      roomInfo.Platform,
		RoomID:        roomInfo.RoomID,
		Status:        string(roomInfo.Status),
		Title:         roomInfo.Title,
		Description:   roomInfo.Description,
		Cover:         roomInfo.Cover,
		Keyframe:      roomInfo.Keyframe,
		OwnerID:       roomInfo.OwnerID,
		OwnerName:     roomInfo.OwnerName,
		OwnerAvatar:   roomInfo.OwnerAvatar,
		LiveStartTime: roomInfo.LiveStartTime,
		ViewerCount:   roomInfo.ViewerCount,
		Category:      roomInfo.Category,
	}

	return c.JSON(response)
}
