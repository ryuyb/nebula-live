package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// AuthRouter 认证路由器
type AuthRouter struct {
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewAuthRouter 创建认证路由器
func NewAuthRouter(authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) Router {
	return &AuthRouter{
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册认证相关路由
func (r *AuthRouter) RegisterRoutes(router fiber.Router) {
	// 认证路由组
	auth := router.Group("/auth")
	
	// 公开认证路由（不需要token）
	{
		auth.Post("/register", r.authHandler.Register)     // 用户注册
		auth.Post("/login", r.authHandler.Login)           // 用户登录
		auth.Post("/refresh", r.authHandler.RefreshToken)  // 刷新令牌
	}

	// 需要认证的路由
	authenticated := auth.Use(r.authMiddleware.RequireAuth())
	{
		authenticated.Get("/me", r.authHandler.GetCurrentUser) // 获取当前用户信息
	}
}

// GetPrefix 获取路由前缀
func (r *AuthRouter) GetPrefix() string {
	return "/api/v1"
}