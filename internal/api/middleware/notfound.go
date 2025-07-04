package middleware

import (
	"github.com/gofiber/fiber/v2"
	"nebulaLive/internal/entity"
)

// NotFound returns a middleware that handles requests to undefined routes by returning a 404 status.
func NotFound(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(entity.NotFoundResponse())
	})
}
