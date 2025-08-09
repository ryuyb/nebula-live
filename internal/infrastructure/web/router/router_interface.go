package router

import "github.com/gofiber/fiber/v2"

// Router 路由接口
type Router interface {
	RegisterRoutes(router fiber.Router)
	GetPrefix() string
}
