package handler

import (
	"nebula-live/internal/domain/service"
	"nebula-live/internal/infrastructure/config"
	"nebula-live/pkg/auth"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService service.UserService
	jwtManager  *auth.JWTManager
	logger      *zap.Logger
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService service.UserService, config *config.Config, logger *zap.Logger) *AuthHandler {
	// 创建JWT管理器
	tokenConfig := &auth.TokenConfig{
		SecretKey:       config.JWT.Secret,
		AccessTokenTTL:  config.JWT.AccessTokenTTL,
		RefreshTokenTTL: config.JWT.RefreshTokenTTL,
		Issuer:          config.JWT.Issuer,
	}
	
	return &AuthHandler{
		userService: userService,
		jwtManager:  auth.NewJWTManager(tokenConfig),
		logger:      logger,
	}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Nickname string `json:"nickname" validate:"max=100"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresAt    int64        `json:"expires_at"`
	Message      string       `json:"message"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse register request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// TODO: 添加请求验证

	user, err := h.userService.CreateUser(c.Context(), req.Username, req.Email, req.Password, req.Nickname)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		
		if err == service.ErrUserAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(errors.NewAPIError(fiber.StatusConflict, "User already exists", "Username or email already exists"))
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to register user"))
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response := AuthResponse{
		User:    userResponse,
		Message: "User registered successfully",
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Login 用户登录
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// TODO: 添加请求验证

	user, err := h.userService.ValidateUser(c.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to validate user credentials", 
			zap.String("username", req.Username), 
			zap.Error(err))
		
		switch err {
		case service.ErrInvalidCredentials:
			return c.Status(fiber.StatusUnauthorized).JSON(errors.NewAPIError(fiber.StatusUnauthorized, "Invalid credentials", "Username or password is incorrect"))
		case service.ErrUserBanned:
			return c.Status(fiber.StatusForbidden).JSON(errors.NewAPIError(fiber.StatusForbidden, "Account banned", "Your account has been banned"))
		case service.ErrUserInactive:
			return c.Status(fiber.StatusForbidden).JSON(errors.NewAPIError(fiber.StatusForbidden, "Account inactive", "Your account is inactive"))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to authenticate user"))
		}
	}

	// 生成JWT令牌
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate JWT tokens", 
			zap.Uint("user_id", user.ID), 
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to generate authentication tokens"))
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response := AuthResponse{
		User:         userResponse,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresAt:    tokenPair.ExpiresAt,
		Message:      "Login successful",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	// 从上下文中获取当前用户信息
	currentUser, exists := auth.GetCurrentUser(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "No authenticated user found"))
	}

	// 从数据库获取最新用户信息
	user, err := h.userService.GetUserByID(c.Context(), currentUser.UserID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "Current user not found"))
		}
		
		h.logger.Error("Failed to get current user", 
			zap.Uint("user_id", currentUser.UserID), 
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get current user"))
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse refresh token request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// 使用刷新令牌生成新的令牌对
	tokenPair, err := h.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(errors.NewAPIError(fiber.StatusUnauthorized, "Invalid refresh token", "Failed to refresh authentication token"))
	}

	response := map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    tokenPair.TokenType,
		"expires_at":    tokenPair.ExpiresAt,
		"message":       "Token refreshed successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}