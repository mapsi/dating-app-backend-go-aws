package handler

import (
	"dating-app-backend/internal/auth"
	"dating-app-backend/internal/logger"
	"dating-app-backend/internal/storage"
	"strconv"

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

	currentUser, err := h.storage.GetUserByID(ctx.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get current user", "error", err, "userID", userID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get current user"})
	}

	minAge, _ := strconv.Atoi(ctx.Query("minAge", "0"))
	maxAge, _ := strconv.Atoi(ctx.Query("maxAge", "0"))
	gender := ctx.Query("gender", "")
	sortBy := ctx.Query("sortBy", "combined") // Default to combined sorting

	h.logger.Info("Discovering users", "userID", userID, "minAge", minAge, "maxAge", maxAge, "gender", gender, "sortBy", sortBy)
	// TODO: Implement pagination
	discoveredUsers, err := h.storage.DiscoverUsers(ctx.Context(), *currentUser, 10, minAge, maxAge, gender, sortBy)
	if err != nil {
		h.logger.Error("Failed to discover users", "error", err, "userID", userID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to discover users"})
	}

	h.logger.Info("Users discovered successfully", "userID", userID, "count", len(discoveredUsers))
	return ctx.JSON(fiber.Map{"results": discoveredUsers})
}
