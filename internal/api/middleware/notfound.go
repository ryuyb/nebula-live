package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// NotFound returns a middleware that handles requests to undefined routes by returning a 404 status.
func NotFound(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "Not Found",
		})
	})
}
