package handler

import (
	"strconv"

	"nebula-live/internal/domain/service"
	"nebula-live/pkg/auth"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	rbacService service.RBACService
	userService service.UserService
	logger      *zap.Logger
}

// NewRoleHandler 创建角色处理器实例
func NewRoleHandler(rbacService service.RBACService, userService service.UserService, logger *zap.Logger) *RoleHandler {
	return &RoleHandler{
		rbacService: rbacService,
		userService: userService,
		logger:      logger,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID uint `json:"user_id" validate:"required,min=1"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListRolesResponse 角色列表响应
type ListRolesResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse create role request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// TODO: 添加请求验证

	role, err := h.rbacService.CreateRole(c.Context(), req.Name, req.DisplayName, req.Description, false)
	if err != nil {
		h.logger.Error("Failed to create role", zap.Error(err))

		if err == service.ErrRoleAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(errors.NewAPIError(fiber.StatusConflict, "Role already exists", "A role with this name already exists"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to create role"))
	}

	response := RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetRole 获取角色信息
func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	role, err := h.rbacService.GetRoleByID(c.Context(), uint(id))
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}

		h.logger.Error("Failed to get role", zap.Error(err), zap.Uint("role_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role"))
	}

	response := RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(response)
}

// UpdateRole 更新角色信息
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse update role request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	role, err := h.rbacService.UpdateRole(c.Context(), uint(id), req.DisplayName, req.Description)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}

		h.logger.Error("Failed to update role", zap.Error(err), zap.Uint("role_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to update role"))
	}

	response := RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(response)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	if err := h.rbacService.DeleteRole(c.Context(), uint(id)); err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}
		if err == service.ErrSystemRoleCannotDelete {
			return c.Status(fiber.StatusForbidden).JSON(errors.NewAPIError(fiber.StatusForbidden, "Cannot delete system role", "System roles cannot be deleted"))
		}

		h.logger.Error("Failed to delete role", zap.Error(err), zap.Uint("role_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to delete role"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	roles, err := h.rbacService.ListRoles(c.Context(), offset, limit)
	if err != nil {
		h.logger.Error("Failed to list roles", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to list roles"))
	}

	roleResponses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListRolesResponse{
		Roles: roleResponses,
		Total: len(roleResponses),
		Page:  page,
		Limit: limit,
	}

	return c.JSON(response)
}

// AssignRole 为用户分配角色
func (h *RoleHandler) AssignRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	roleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	var req AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse assign role request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// 获取当前用户作为分配者
	currentUser, exists := auth.GetCurrentUser(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Authentication required"))
	}

	// 检查角色是否存在
	role, err := h.rbacService.GetRoleByID(c.Context(), uint(roleID))
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}
		h.logger.Error("Failed to get role for assignment", zap.Error(err), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role"))
	}

	// 使用用户服务分配角色
	if err := h.userService.AssignRole(c.Context(), req.UserID, role.Name, currentUser.UserID); err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}
		if err == service.ErrUserRoleAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(errors.NewAPIError(fiber.StatusConflict, "Role already assigned", "User already has this role"))
		}

		h.logger.Error("Failed to assign role to user", zap.Error(err), zap.Uint("user_id", req.UserID), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to assign role"))
	}

	return c.JSON(fiber.Map{
		"message": "Role assigned successfully",
	})
}

// RemoveRole 移除用户角色
func (h *RoleHandler) RemoveRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	roleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	userIDStr := c.Params("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	// 检查角色是否存在
	role, err := h.rbacService.GetRoleByID(c.Context(), uint(roleID))
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}
		h.logger.Error("Failed to get role for removal", zap.Error(err), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role"))
	}

	// 使用用户服务移除角色
	if err := h.userService.RemoveRole(c.Context(), uint(userID), role.Name); err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}
		if err == service.ErrUserRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User role not found", "User does not have this role"))
		}

		h.logger.Error("Failed to remove role from user", zap.Error(err), zap.Uint("user_id", uint(userID)), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to remove role"))
	}

	return c.JSON(fiber.Map{
		"message": "Role removed successfully",
	})
}

// GetUserRoles 获取用户的所有角色
func (h *RoleHandler) GetUserRoles(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	roles, err := h.userService.GetUserRoles(c.Context(), uint(userID))
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "User not found", "User with the given ID does not exist"))
		}

		h.logger.Error("Failed to get user roles", zap.Error(err), zap.Uint("user_id", uint(userID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get user roles"))
	}

	roleResponses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(fiber.Map{
		"roles": roleResponses,
	})
}