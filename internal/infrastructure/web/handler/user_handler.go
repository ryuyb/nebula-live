package handler

import (
	"strconv"

	"nebula-live/internal/domain/service"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Nickname string `json:"nickname" validate:"max=100"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname" validate:"max=100"`
	Avatar   string `json:"avatar" validate:"max=500"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// CreateUser godoc
// @Summary      Create User
// @Description  Create a new user in the system
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "User creation data"
// @Success      201 {object} UserResponse "User created successfully"
// @Failure      400 {object} errors.APIError "Invalid request parameters"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      409 {object} errors.APIError "User already exists"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse create user request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// TODO: 添加请求验证

	user, err := h.userService.CreateUser(c.Context(), req.Username, req.Email, req.Password, req.Nickname)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))

		if err == service.ErrUserAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(errors.NewAPIError(fiber.StatusConflict, "User already exists", "Username or email already exists"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to create user"))
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetUser godoc
// @Summary      Get User
// @Description  Get user information by ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} UserResponse "User retrieved successfully"
// @Failure      400 {object} errors.APIError "Invalid user ID"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      404 {object} errors.APIError "User not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	user, err := h.userService.GetUserByID(c.Context(), uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to get user", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get user"))
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(response)
}

// UpdateUser godoc
// @Summary      Update User
// @Description  Update user information
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        user body UpdateUserRequest true "User update data"
// @Success      200 {object} UserResponse "User updated successfully"
// @Failure      400 {object} errors.APIError "Invalid request parameters"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      404 {object} errors.APIError "User not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse update user request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// 获取现有用户
	user, err := h.userService.GetUserByID(c.Context(), uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to get user for update", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get user"))
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := h.userService.UpdateUser(c.Context(), user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to update user"))
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(response)
}

// DeleteUser godoc
// @Summary      Delete User
// @Description  Delete a user from the system
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      204 "User deleted successfully"
// @Failure      400 {object} errors.APIError "Invalid user ID"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      404 {object} errors.APIError "User not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	if err := h.userService.DeleteUser(c.Context(), uint(id)); err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to delete user", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to delete user"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ListUsers godoc
// @Summary      List Users
// @Description  Get list of users with pagination
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Success      200 {object} ListUsersResponse "List of users"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// 解析分页参数
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := h.userService.ListUsers(c.Context(), offset, limit)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to list users"))
	}

	// 获取总数
	total, err := h.userService.CountUsers(c.Context())
	if err != nil {
		h.logger.Error("Failed to count users", zap.Error(err))
		// 如果获取总数失败，仍然返回用户列表，但总数设为-1
		total = -1
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status.String(),
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListUsersResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return c.JSON(response)
}

// ActivateUser godoc
// @Summary      Activate User
// @Description  Activate a user account
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string "User activated successfully"
// @Failure      400 {object} errors.APIError "Invalid user ID"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      404 {object} errors.APIError "User not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	if err := h.userService.ActivateUser(c.Context(), uint(id)); err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to activate user", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to activate user"))
	}

	return c.JSON(fiber.Map{
		"message": "User activated successfully",
	})
}

// DeactivateUser godoc
// @Summary      Deactivate User
// @Description  Deactivate a user account
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string "User deactivated successfully"
// @Failure      400 {object} errors.APIError "Invalid user ID"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      404 {object} errors.APIError "User not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	if err := h.userService.DeactivateUser(c.Context(), uint(id)); err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to deactivate user", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to deactivate user"))
	}

	return c.JSON(fiber.Map{
		"message": "User deactivated successfully",
	})
}

// BanUser godoc
// @Summary      Ban User
// @Description  Ban a user account
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string "User banned successfully"
// @Failure      400 {object} errors.APIError "Invalid user ID"
// @Failure      401 {object} errors.APIError "Unauthorized"
// @Failure      404 {object} errors.APIError "User not found"
// @Failure      500 {object} errors.APIError "Internal server error"
// @Security     Bearer
// @Router       /users/{id}/ban [post]
func (h *UserHandler) BanUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	if err := h.userService.BanUser(c.Context(), uint(id)); err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to ban user", zap.Error(err), zap.Uint("user_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to ban user"))
	}

	return c.JSON(fiber.Map{
		"message": "User banned successfully",
	})
}
