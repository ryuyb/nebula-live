package middleware

import (
	"nebula-live/internal/domain/service"
	"nebula-live/pkg/auth"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RBACMiddleware RBAC权限验证中间件
type RBACMiddleware struct {
	rbacService service.RBACService
	logger      *zap.Logger
}

// NewRBACMiddleware 创建RBAC中间件
func NewRBACMiddleware(rbacService service.RBACService, logger *zap.Logger) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
		logger:      logger,
	}
}

// RequirePermission 要求指定权限的中间件
func (m *RBACMiddleware) RequirePermission(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文获取当前用户
		currentUser, exists := auth.GetCurrentUser(c)
		if !exists {
			m.logger.Debug("No authenticated user found for permission check",
				zap.String("resource", resource),
				zap.String("action", action))
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Authentication required"),
			)
		}

		// 检查用户权限
		hasPermission, err := m.rbacService.HasPermission(c.Context(), currentUser.UserID, resource, action)
		if err != nil {
			m.logger.Error("Failed to check user permission",
				zap.Uint("user_id", currentUser.UserID),
				zap.String("resource", resource),
				zap.String("action", action),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(
				errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to verify permissions"),
			)
		}

		if !hasPermission {
			m.logger.Debug("User lacks required permission",
				zap.Uint("user_id", currentUser.UserID),
				zap.String("username", currentUser.Username),
				zap.String("resource", resource),
				zap.String("action", action))
			return c.Status(fiber.StatusForbidden).JSON(
				errors.NewAPIError(fiber.StatusForbidden, "Forbidden", "Insufficient permissions"),
			)
		}

		m.logger.Debug("Permission check passed",
			zap.Uint("user_id", currentUser.UserID),
			zap.String("username", currentUser.Username),
			zap.String("resource", resource),
			zap.String("action", action))

		return c.Next()
	}
}

// RequireRole 要求指定角色的中间件
func (m *RBACMiddleware) RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文获取当前用户
		currentUser, exists := auth.GetCurrentUser(c)
		if !exists {
			m.logger.Debug("No authenticated user found for role check",
				zap.String("role", roleName))
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Authentication required"),
			)
		}

		// 检查用户角色
		hasRole, err := m.rbacService.HasRole(c.Context(), currentUser.UserID, roleName)
		if err != nil {
			m.logger.Error("Failed to check user role",
				zap.Uint("user_id", currentUser.UserID),
				zap.String("role", roleName),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(
				errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to verify role"),
			)
		}

		if !hasRole {
			m.logger.Debug("User lacks required role",
				zap.Uint("user_id", currentUser.UserID),
				zap.String("username", currentUser.Username),
				zap.String("role", roleName))
			return c.Status(fiber.StatusForbidden).JSON(
				errors.NewAPIError(fiber.StatusForbidden, "Forbidden", "Required role not found"),
			)
		}

		m.logger.Debug("Role check passed",
			zap.Uint("user_id", currentUser.UserID),
			zap.String("username", currentUser.Username),
			zap.String("role", roleName))

		return c.Next()
	}
}

// RequireAdmin 要求管理员角色的中间件
func (m *RBACMiddleware) RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 从上下文获取当前用户
		currentUser, exists := auth.GetCurrentUser(c)
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Authentication required"),
			)
		}

		// 检查是否为管理员
		isAdmin, err := m.rbacService.HasRole(c.Context(), currentUser.UserID, "admin")
		if err != nil {
			m.logger.Error("Failed to check admin role",
				zap.Uint("user_id", currentUser.UserID),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(
				errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to verify admin role"),
			)
		}

		if !isAdmin {
			m.logger.Debug("User is not an admin",
				zap.Uint("user_id", currentUser.UserID),
				zap.String("username", currentUser.Username))
			return c.Status(fiber.StatusForbidden).JSON(
				errors.NewAPIError(fiber.StatusForbidden, "Forbidden", "Administrator privileges required"),
			)
		}

		m.logger.Debug("Admin check passed",
			zap.Uint("user_id", currentUser.UserID),
			zap.String("username", currentUser.Username))

		return c.Next()
	}
}
