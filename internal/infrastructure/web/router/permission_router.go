package router

import (
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// PermissionRouter 权限路由器
type PermissionRouter struct {
	permissionHandler *handler.PermissionHandler
	authMiddleware    *middleware.AuthMiddleware
	rbacMiddleware    *middleware.RBACMiddleware
}

// NewPermissionRouter 创建权限路由器
func NewPermissionRouter(permissionHandler *handler.PermissionHandler, authMiddleware *middleware.AuthMiddleware, rbacMiddleware *middleware.RBACMiddleware) Router {
	return &PermissionRouter{
		permissionHandler: permissionHandler,
		authMiddleware:    authMiddleware,
		rbacMiddleware:    rbacMiddleware,
	}
}

// RegisterRoutes 注册权限相关路由
func (r *PermissionRouter) RegisterRoutes(router fiber.Router) {
	// 权限管理路由组 - 需要认证和admin角色
	permissions := router.Group("/permissions").Use(
		r.authMiddleware.RequireAuth(),
		r.rbacMiddleware.RequireAdmin(),
	)
	{
		// 基础CRUD操作
		permissions.Post("/", r.permissionHandler.CreatePermission)      // 创建权限
		permissions.Get("/:id", r.permissionHandler.GetPermission)       // 获取权限信息
		permissions.Put("/:id", r.permissionHandler.UpdatePermission)    // 更新权限信息
		permissions.Delete("/:id", r.permissionHandler.DeletePermission) // 删除权限
		permissions.Get("/", r.permissionHandler.ListPermissions)        // 获取权限列表

		// 权限分配管理
		permissions.Post("/:id/assign", r.permissionHandler.AssignPermissionToRole)            // 为角色分配权限
		permissions.Delete("/:id/roles/:roleId", r.permissionHandler.RemovePermissionFromRole) // 移除角色权限
		permissions.Get("/roles/:roleId", r.permissionHandler.GetRolePermissions)              // 获取角色的所有权限
		permissions.Get("/users/:userId", r.permissionHandler.GetUserPermissions)              // 获取用户的所有权限
	}
}

// GetPrefix 获取路由前缀
func (r *PermissionRouter) GetPrefix() string {
	return "/api/v1"
}
