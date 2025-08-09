package handler

import (
	"strconv"

	"nebula-live/internal/domain/service"
	"nebula-live/pkg/auth"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	rbacService service.RBACService
	logger      *zap.Logger
}

// NewPermissionHandler 创建权限处理器实例
func NewPermissionHandler(rbacService service.RBACService, logger *zap.Logger) *PermissionHandler {
	return &PermissionHandler{
		rbacService: rbacService,
		logger:      logger,
	}
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
	Resource    string `json:"resource" validate:"required,min=2,max=50"`
	Action      string `json:"action" validate:"required,min=2,max=50"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// AssignPermissionToRoleRequest 为角色分配权限请求
type AssignPermissionToRoleRequest struct {
	RoleID uint `json:"role_id" validate:"required,min=1"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	IsSystem    bool   `json:"is_system"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListPermissionsResponse 权限列表响应
type ListPermissionsResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var req CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse create permission request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// TODO: 添加请求验证

	permission, err := h.rbacService.CreatePermission(c.Context(), req.Name, req.DisplayName, req.Description, req.Resource, req.Action, false)
	if err != nil {
		h.logger.Error("Failed to create permission", zap.Error(err))

		if err == service.ErrPermissionAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(errors.NewAPIError(fiber.StatusConflict, "Permission already exists", "A permission with this name already exists"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to create permission"))
	}

	response := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		DisplayName: permission.DisplayName,
		Description: permission.Description,
		Resource:    permission.Resource,
		Action:      permission.Action,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetPermission 获取权限信息
func (h *PermissionHandler) GetPermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid permission ID", "Permission ID must be a valid number"))
	}

	permission, err := h.rbacService.GetPermissionByID(c.Context(), uint(id))
	if err != nil {
		if err == service.ErrPermissionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Permission not found", "Permission with the given ID does not exist"))
		}

		h.logger.Error("Failed to get permission", zap.Error(err), zap.Uint("permission_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get permission"))
	}

	response := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		DisplayName: permission.DisplayName,
		Description: permission.Description,
		Resource:    permission.Resource,
		Action:      permission.Action,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(response)
}

// UpdatePermission 更新权限信息
func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid permission ID", "Permission ID must be a valid number"))
	}

	var req UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse update permission request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	permission, err := h.rbacService.UpdatePermission(c.Context(), uint(id), req.DisplayName, req.Description)
	if err != nil {
		if err == service.ErrPermissionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Permission not found", "Permission with the given ID does not exist"))
		}

		h.logger.Error("Failed to update permission", zap.Error(err), zap.Uint("permission_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to update permission"))
	}

	response := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		DisplayName: permission.DisplayName,
		Description: permission.Description,
		Resource:    permission.Resource,
		Action:      permission.Action,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(response)
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid permission ID", "Permission ID must be a valid number"))
	}

	if err := h.rbacService.DeletePermission(c.Context(), uint(id)); err != nil {
		if err == service.ErrPermissionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Permission not found", "Permission with the given ID does not exist"))
		}
		if err == service.ErrSystemPermissionCannotDelete {
			return c.Status(fiber.StatusForbidden).JSON(errors.NewAPIError(fiber.StatusForbidden, "Cannot delete system permission", "System permissions cannot be deleted"))
		}

		h.logger.Error("Failed to delete permission", zap.Error(err), zap.Uint("permission_id", uint(id)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to delete permission"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	permissions, err := h.rbacService.ListPermissions(c.Context(), offset, limit)
	if err != nil {
		h.logger.Error("Failed to list permissions", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to list permissions"))
	}

	permissionResponses := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResponses[i] = PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			DisplayName: permission.DisplayName,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
			IsSystem:    permission.IsSystem,
			CreatedAt:   permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListPermissionsResponse{
		Permissions: permissionResponses,
		Total:       len(permissionResponses),
		Page:        page,
		Limit:       limit,
	}

	return c.JSON(response)
}

// AssignPermissionToRole 为角色分配权限
func (h *PermissionHandler) AssignPermissionToRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	permissionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid permission ID", "Permission ID must be a valid number"))
	}

	var req AssignPermissionToRoleRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse assign permission request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid request body", err.Error()))
	}

	// 获取当前用户作为分配者
	currentUser, exists := auth.GetCurrentUser(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.NewAPIError(fiber.StatusUnauthorized, "Unauthorized", "Authentication required"))
	}

	// 检查权限是否存在
	_, err = h.rbacService.GetPermissionByID(c.Context(), uint(permissionID))
	if err != nil {
		if err == service.ErrPermissionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Permission not found", "Permission with the given ID does not exist"))
		}
		h.logger.Error("Failed to get permission for assignment", zap.Error(err), zap.Uint("permission_id", uint(permissionID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get permission"))
	}

	// 检查角色是否存在
	_, err = h.rbacService.GetRoleByID(c.Context(), req.RoleID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}
		h.logger.Error("Failed to get role for permission assignment", zap.Error(err), zap.Uint("role_id", req.RoleID))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role"))
	}

	// 分配权限到角色
	if err := h.rbacService.AssignPermissionToRole(c.Context(), req.RoleID, uint(permissionID), currentUser.UserID); err != nil {
		if err == service.ErrRolePermissionAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(errors.NewAPIError(fiber.StatusConflict, "Permission already assigned", "Role already has this permission"))
		}

		h.logger.Error("Failed to assign permission to role", zap.Error(err), zap.Uint("role_id", req.RoleID), zap.Uint("permission_id", uint(permissionID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to assign permission"))
	}

	return c.JSON(fiber.Map{
		"message": "Permission assigned to role successfully",
	})
}

// RemovePermissionFromRole 移除角色的权限
func (h *PermissionHandler) RemovePermissionFromRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	permissionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid permission ID", "Permission ID must be a valid number"))
	}

	roleIDStr := c.Params("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	// 检查权限是否存在
	_, err = h.rbacService.GetPermissionByID(c.Context(), uint(permissionID))
	if err != nil {
		if err == service.ErrPermissionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Permission not found", "Permission with the given ID does not exist"))
		}
		h.logger.Error("Failed to get permission for removal", zap.Error(err), zap.Uint("permission_id", uint(permissionID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get permission"))
	}

	// 检查角色是否存在
	_, err = h.rbacService.GetRoleByID(c.Context(), uint(roleID))
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}
		h.logger.Error("Failed to get role for permission removal", zap.Error(err), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role"))
	}

	// 移除角色的权限
	if err := h.rbacService.RemovePermissionFromRole(c.Context(), uint(roleID), uint(permissionID)); err != nil {
		if err == service.ErrRolePermissionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role permission not found", "Role does not have this permission"))
		}

		h.logger.Error("Failed to remove permission from role", zap.Error(err), zap.Uint("role_id", uint(roleID)), zap.Uint("permission_id", uint(permissionID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to remove permission"))
	}

	return c.JSON(fiber.Map{
		"message": "Permission removed from role successfully",
	})
}

// GetRolePermissions 获取角色的所有权限
func (h *PermissionHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleIDStr := c.Params("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid role ID", "Role ID must be a valid number"))
	}

	// 检查角色是否存在
	_, err = h.rbacService.GetRoleByID(c.Context(), uint(roleID))
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAPIError(fiber.StatusNotFound, "Role not found", "Role with the given ID does not exist"))
		}
		h.logger.Error("Failed to get role for permissions", zap.Error(err), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role"))
	}

	permissions, err := h.rbacService.GetRolePermissions(c.Context(), uint(roleID))
	if err != nil {
		h.logger.Error("Failed to get role permissions", zap.Error(err), zap.Uint("role_id", uint(roleID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get role permissions"))
	}

	permissionResponses := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResponses[i] = PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			DisplayName: permission.DisplayName,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
			IsSystem:    permission.IsSystem,
			CreatedAt:   permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(fiber.Map{
		"permissions": permissionResponses,
	})
}

// GetUserPermissions 获取用户的所有权限
func (h *PermissionHandler) GetUserPermissions(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewAPIError(fiber.StatusBadRequest, "Invalid user ID", "User ID must be a valid number"))
	}

	permissions, err := h.rbacService.GetUserPermissions(c.Context(), uint(userID))
	if err != nil {
		h.logger.Error("Failed to get user permissions", zap.Error(err), zap.Uint("user_id", uint(userID)))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewAPIError(fiber.StatusInternalServerError, "Internal server error", "Failed to get user permissions"))
	}

	permissionResponses := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResponses[i] = PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			DisplayName: permission.DisplayName,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
			IsSystem:    permission.IsSystem,
			CreatedAt:   permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(fiber.Map{
		"permissions": permissionResponses,
	})
}