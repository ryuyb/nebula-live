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

func NewLiveStreamHandler(liveStreamService service.LiveStreamService, logger *zap.Logger) *LiveStreamHandler {
	return &LiveStreamHandler{
		liveStreamService: liveStreamService,
		logger:            logger,
	}
}

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

	return c.JSON(streamInfo)
}

func (h *LiveStreamHandler) GetSupportedPlatforms(c *fiber.Ctx) error {
	platforms := h.liveStreamService.GetSupportedPlatforms()

	return c.JSON(fiber.Map{
		"platforms": platforms,
	})
}
