package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// UserPushSettingRouter 用户推送设置路由器
type UserPushSettingRouter struct {
	handler        *handler.UserPushSettingHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewUserPushSettingRouter 创建用户推送设置路由器
func NewUserPushSettingRouter(
	handler *handler.UserPushSettingHandler,
	authMiddleware *middleware.AuthMiddleware,
) Router {
	return &UserPushSettingRouter{
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册用户推送设置路由
func (r *UserPushSettingRouter) RegisterRoutes(router fiber.Router) {
	// 用户推送设置路由组
	pushSettings := router.Group("/push-settings")

	// 所有用户推送设置相关操作都需要认证
	pushSettings.Use(r.authMiddleware.RequireAuth())

	// 公开端点（不需要认证）
	router.Get("/push-settings/providers", r.handler.GetSupportedProviders)     // 获取支持的推送提供商
	router.Post("/push-settings/validate-device", r.handler.ValidateDevice)     // 验证设备ID是否可用
	
	// 用户推送设置管理
	pushSettings.Post("/", r.handler.CreateSetting)      // 创建推送设置
	pushSettings.Get("/", r.handler.GetSettings)         // 获取推送设置列表
	pushSettings.Get("/:id", r.handler.GetSetting)       // 获取指定推送设置
	pushSettings.Put("/:id", r.handler.UpdateSetting)    // 更新推送设置
	pushSettings.Delete("/:id", r.handler.DeleteSetting) // 删除推送设置

	// 推送设置状态管理
	pushSettings.Post("/:id/enable", r.handler.EnableSetting)   // 启用推送设置
	pushSettings.Post("/:id/disable", r.handler.DisableSetting) // 禁用推送设置
}

// GetPrefix 获取路由前缀
func (r *UserPushSettingRouter) GetPrefix() string {
	return "/api/v1"
}
