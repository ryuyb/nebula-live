package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// UserRouter 用户路由器
type UserRouter struct {
	userHandler    *handler.UserHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewUserRouter 创建用户路由器
func NewUserRouter(userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) Router {
	return &UserRouter{
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册用户相关路由
func (r *UserRouter) RegisterRoutes(router fiber.Router) {
	// 用户路由组 - 所有路由都需要认证
	users := router.Group("/users").Use(r.authMiddleware.RequireAuth())
	{
		users.Post("/", r.userHandler.CreateUser)           // 创建用户
		users.Get("/:id", r.userHandler.GetUser)            // 获取用户信息
		users.Put("/:id", r.userHandler.UpdateUser)         // 更新用户信息
		users.Delete("/:id", r.userHandler.DeleteUser)      // 删除用户
		users.Get("/", r.userHandler.ListUsers)             // 获取用户列表
		
		// 用户状态管理
		users.Post("/:id/activate", r.userHandler.ActivateUser)     // 激活用户
		users.Post("/:id/deactivate", r.userHandler.DeactivateUser) // 停用用户
		users.Post("/:id/ban", r.userHandler.BanUser)               // 禁用用户
	}
}

// GetPrefix 获取路由前缀
func (r *UserRouter) GetPrefix() string {
	return "/api/v1"
}