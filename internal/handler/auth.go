package handler

import (
	"dating-app-backend/internal/auth"
	"dating-app-backend/internal/logger"
	"dating-app-backend/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	storage *storage.DynamoDB
	logger  *logger.Logger
}

func NewAuthHandler(storage *storage.DynamoDB, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{storage: storage, logger: logger}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		h.logger.Error("Failed to parse login input", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.storage.GetUserByEmail(ctx.Context(), input.Email)
	if err != nil {
		h.logger.Error("Failed to get user by email", "error", err, "email", input.Email)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !user.CheckPassword(input.Password) {
		h.logger.Warn("Invalid password attempt", "email", input.Email)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		h.logger.Error("Failed to generate token", "error", err, "userId", user.ID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	h.logger.Info("User logged in successfully", "userId", user.ID)
	return ctx.JSON(fiber.Map{"token": token})
}
