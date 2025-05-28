package router

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRouter sets up the API routes using RouterRegistry.
func SetupRouter(app *fiber.App, registry *RouterRegistry) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Setup all registered routes
	registry.SetupAllRoutes(&v1)
}
