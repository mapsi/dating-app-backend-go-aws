package handler

import (
	"dating-app-backend/internal/auth"
	"dating-app-backend/internal/logger"
	"dating-app-backend/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type DiscoverHandler struct {
	storage *storage.DynamoDB
	logger  *logger.Logger
}

func NewDiscoverHandler(storage *storage.DynamoDB, logger *logger.Logger) *DiscoverHandler {
	return &DiscoverHandler{storage: storage, logger: logger}
}

func (h *DiscoverHandler) DiscoverUsers(ctx *fiber.Ctx) error {
	userID, err := auth.GetUserIDFromToken(ctx)
	if err != nil {
		h.logger.Error("Failed to get user ID from token", "error", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	h.logger.Info("Discovering users", "userID", userID)
	// TODO: Implement pagination
	discoveredUsers, err := h.storage.DiscoverUsers(ctx.Context(), userID, 10) // Limit to 10 users
	if err != nil {
		h.logger.Error("Failed to discover users", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to discover users"})
	}

	h.logger.Info("Users discovered successfully", "userID", userID, "count", len(discoveredUsers))
	return ctx.JSON(fiber.Map{"results": discoveredUsers})
}
