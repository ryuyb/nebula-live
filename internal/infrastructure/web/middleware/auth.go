package middleware

import (
	"strings"

	"nebula-live/internal/infrastructure/config"
	"nebula-live/pkg/auth"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	// AuthContextKey 认证上下文键
	AuthContextKey = "auth_user"
	// UserIDContextKey 用户ID上下文键
	UserIDContextKey = "user_id"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(config *config.Config, logger *zap.Logger) *AuthMiddleware {
	tokenConfig := &auth.TokenConfig{
		SecretKey:       config.JWT.Secret,
		AccessTokenTTL:  config.JWT.AccessTokenTTL,
		RefreshTokenTTL: config.JWT.RefreshTokenTTL,
		Issuer:          config.JWT.Issuer,
	}

	return &AuthMiddleware{
		jwtManager: auth.NewJWTManager(tokenConfig),
		logger:     logger,
	}
}

// RequireAuth 要求认证的中间件
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取Authorization头
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			m.logger.Debug("Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Missing authorization header"),
			)
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Debug("Invalid authorization header format", zap.String("header", authHeader))
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Invalid authorization header format"),
			)
		}

		token := parts[1]
		if token == "" {
			m.logger.Debug("Empty token")
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Empty token"),
			)
		}

		// 验证token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			m.logger.Debug("Token validation failed", 
				zap.Error(err),
				zap.String("token", token[:min(len(token), 50)]+"..."))

			switch err {
			case auth.ErrExpiredToken:
				return c.Status(fiber.StatusUnauthorized).JSON(
					errors.NewAPIError(fiber.StatusUnauthorized, "Token expired", "Your session has expired, please login again"),
				)
			case auth.ErrInvalidToken:
				return c.Status(fiber.StatusUnauthorized).JSON(
					errors.NewAPIError(fiber.StatusUnauthorized, "Invalid token", "Invalid authentication token"),
				)
			case auth.ErrTokenClaims:
				return c.Status(fiber.StatusUnauthorized).JSON(
					errors.NewAPIError(fiber.StatusUnauthorized, "Invalid token claims", "Invalid token claims"),
				)
			default:
				return c.Status(fiber.StatusUnauthorized).JSON(
					errors.NewAPIError(fiber.StatusUnauthorized, "Authentication failed", "Token validation failed"),
				)
			}
		}

		// 将用户信息存储到上下文中
		c.Locals(AuthContextKey, claims)
		c.Locals(UserIDContextKey, claims.UserID)

		m.logger.Debug("User authenticated successfully", 
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username))

		return c.Next()
	}
}

// OptionalAuth 可选认证的中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取Authorization头
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// 没有认证头，直接继续
			return c.Next()
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 格式不正确，直接继续（不返回错误）
			return c.Next()
		}

		token := parts[1]
		if token == "" {
			// 空token，直接继续
			return c.Next()
		}

		// 验证token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			// token无效，记录日志但不返回错误
			m.logger.Debug("Optional auth token validation failed", 
				zap.Error(err),
				zap.String("token", token[:min(len(token), 50)]+"..."))
			return c.Next()
		}

		// token有效，将用户信息存储到上下文中
		c.Locals(AuthContextKey, claims)
		c.Locals(UserIDContextKey, claims.UserID)

		m.logger.Debug("User optionally authenticated", 
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username))

		return c.Next()
	}
}

// min 获取两个数的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}