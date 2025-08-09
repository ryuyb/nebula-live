package auth

import (
	"github.com/gofiber/fiber/v2"
)

const (
	// AuthContextKey 认证上下文键
	AuthContextKey = "auth_user"
	// UserIDContextKey 用户ID上下文键
	UserIDContextKey = "user_id"
)

// GetCurrentUser 从上下文中获取当前用户信息
func GetCurrentUser(c *fiber.Ctx) (*UserClaims, bool) {
	user := c.Locals(AuthContextKey)
	if user == nil {
		return nil, false
	}

	claims, ok := user.(*UserClaims)
	return claims, ok
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *fiber.Ctx) (uint, bool) {
	userID := c.Locals(UserIDContextKey)
	if userID == nil {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// MustGetCurrentUser 从上下文中获取当前用户信息（必须存在，否则panic）
func MustGetCurrentUser(c *fiber.Ctx) *UserClaims {
	user, exists := GetCurrentUser(c)
	if !exists {
		panic("no authenticated user found in context")
	}
	return user
}

// MustGetCurrentUserID 从上下文中获取当前用户ID（必须存在，否则panic）
func MustGetCurrentUserID(c *fiber.Ctx) uint {
	userID, exists := GetCurrentUserID(c)
	if !exists {
		panic("no authenticated user ID found in context")
	}
	return userID
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c *fiber.Ctx) bool {
	_, exists := GetCurrentUser(c)
	return exists
}
