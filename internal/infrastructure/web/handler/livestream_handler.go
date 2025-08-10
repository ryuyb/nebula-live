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
