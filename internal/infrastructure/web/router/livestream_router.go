package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

type LiveStreamRouter struct {
	handler        *handler.LiveStreamHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewLiveStreamRouter(
	handler *handler.LiveStreamHandler,
	authMiddleware *middleware.AuthMiddleware,
) Router {
	return &LiveStreamRouter{
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

func (r *LiveStreamRouter) GetPrefix() string {
	return "/api/v1"
}

func (r *LiveStreamRouter) RegisterRoutes(router fiber.Router) {
	liveStreamGroup := router.Group("/live-streams")

	// Get supported platforms (public endpoint)
	liveStreamGroup.Get("/platforms", r.handler.GetSupportedPlatforms)

	// Get stream status (public endpoint)
	liveStreamGroup.Get("/:platform/rooms/:roomId/status", r.handler.GetStreamStatus)
}
