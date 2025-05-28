package router

import (
	"github.com/gofiber/fiber/v2"
	"nebulaLive/internal/api/handler"
)

// UserRouter 结构体实现 Router 接口
type UserRouter struct {
	userHandler *handler.UserHandler
}

// NewUserRouter 创建一个新的 UserRouter
func NewUserRouter(userHandler *handler.UserHandler) *UserRouter {
	return &UserRouter{userHandler: userHandler}
}

// SetupRoutes 设置用户相关的 API 路由
func (r *UserRouter) SetupRoutes(v1 *fiber.Router) {
	user := (*v1).Group("/users")
	user.Get("/", r.userHandler.Test)
	user.Post("/", r.userHandler.CreateUser)
	user.Get("/:id", r.userHandler.GetUserByID)
}
