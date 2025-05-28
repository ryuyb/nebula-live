package handler

import (
	"github.com/gofiber/fiber/v2"
	"nebulaLive/internal/service"
	"strconv"
)

// UserHandler handles HTTP requests for User entity.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Test(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"code": 200,
	})
}

// CreateUser handles the creation of a new user.
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	type CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userService.CreateUser(c.Context(), req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID handles getting a user by ID.
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}
