package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"uservault/internal/models"
	"uservault/internal/service"
)

type UserHandler struct {
	svc    service.UserService
	logger *zap.Logger
}

func NewUserHandler(svc service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.svc.CreateUser(c.Context(), req)
	if err != nil {
		h.logger.Warn("failed to create user", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// For create/update responses, spec does not require age field.
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":   user.ID,
		"name": user.Name,
		"dob":  user.DOB,
	})
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	user, err := h.svc.GetUser(c.Context(), id)
	if err != nil {
		h.logger.Warn("failed to get user", zap.Error(err))
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.svc.UpdateUser(c.Context(), id, req)
	if err != nil {
		h.logger.Warn("failed to update user", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"id":   user.ID,
		"name": user.Name,
		"dob":  user.DOB,
	})
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		h.logger.Warn("failed to delete user", zap.Error(err))
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.SendStatus(http.StatusNoContent)
}

// ListUsers handles GET /users with optional pagination.
// Query params: ?limit=10&offset=0
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit64, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit64 <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid limit"})
	}
	offset64, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil || offset64 < 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid offset"})
	}

	users, err := h.svc.ListUsers(c.Context(), int32(limit64), int32(offset64))
	if err != nil {
		h.logger.Warn("failed to list users", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list users"})
	}

	return c.JSON(users)
}


