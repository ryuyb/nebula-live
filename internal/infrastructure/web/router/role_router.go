package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// RoleRouter 角色路由器
type RoleRouter struct {
	roleHandler    *handler.RoleHandler
	authMiddleware *middleware.AuthMiddleware
	rbacMiddleware *middleware.RBACMiddleware
}

// NewRoleRouter 创建角色路由器
func NewRoleRouter(roleHandler *handler.RoleHandler, authMiddleware *middleware.AuthMiddleware, rbacMiddleware *middleware.RBACMiddleware) Router {
	return &RoleRouter{
		roleHandler:    roleHandler,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

// RegisterRoutes 注册角色相关路由
func (r *RoleRouter) RegisterRoutes(router fiber.Router) {
	// 角色管理路由组 - 需要认证和admin角色
	roles := router.Group("/roles").Use(
		r.authMiddleware.RequireAuth(),
		r.rbacMiddleware.RequireAdmin(),
	)
	{
		// 基础CRUD操作
		roles.Post("/", r.roleHandler.CreateRole)      // 创建角色
		roles.Get("/:id", r.roleHandler.GetRole)       // 获取角色信息
		roles.Put("/:id", r.roleHandler.UpdateRole)    // 更新角色信息
		roles.Delete("/:id", r.roleHandler.DeleteRole) // 删除角色
		roles.Get("/", r.roleHandler.ListRoles)        // 获取角色列表

		// 角色分配管理
		roles.Post("/:id/assign", r.roleHandler.AssignRole)          // 为用户分配角色
		roles.Delete("/:id/users/:userId", r.roleHandler.RemoveRole) // 移除用户角色
		roles.Get("/users/:userId", r.roleHandler.GetUserRoles)      // 获取用户的所有角色
	}
}

// GetPrefix 获取路由前缀
func (r *RoleRouter) GetPrefix() string {
	return "/api/v1"
}
