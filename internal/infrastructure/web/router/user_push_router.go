package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// UserPushRouter 用户推送路由器
type UserPushRouter struct {
	handler        *handler.UserPushHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewUserPushRouter 创建用户推送路由器
func NewUserPushRouter(
	handler *handler.UserPushHandler,
	authMiddleware *middleware.AuthMiddleware,
) Router {
	return &UserPushRouter{
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册用户推送路由
func (r *UserPushRouter) RegisterRoutes(router fiber.Router) {
	// 用户推送路由组
	userPush := router.Group("/push")
	
	// 所有用户推送操作都需要认证
	userPush.Use(r.authMiddleware.RequireAuth())
	
	// 用户推送功能
	userPush.Post("/my-devices", r.handler.SendToMyDevices)                    // 发送到我的所有设备
	userPush.Post("/my-devices/:provider", r.handler.SendToMyDevicesByProvider) // 发送到我指定提供商的设备
	userPush.Post("/test", r.handler.TestMyPushSettings)                       // 测试我的推送设置
}

// GetPrefix 获取路由前缀
func (r *UserPushRouter) GetPrefix() string {
	return "/api/v1"
}