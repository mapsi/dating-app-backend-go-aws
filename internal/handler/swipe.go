package handler

import (
	"dating-app-backend/internal/auth"
	"dating-app-backend/internal/logger"
	"dating-app-backend/internal/model"
	"dating-app-backend/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type SwipeHandler struct {
	storage *storage.DynamoDB
	logger  *logger.Logger
}

func NewSwipeHandler(storage *storage.DynamoDB, logger *logger.Logger) *SwipeHandler {
	return &SwipeHandler{storage: storage, logger: logger}
}

func (h *SwipeHandler) RecordSwipe(c *fiber.Ctx) error {
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		h.logger.Error("Failed to get user ID from token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	var input struct {
		SwipedId   string                `json:"swipedId"`
		Preference model.SwipePreference `json:"preference"`
	}

	if err := c.BodyParser(&input); err != nil {
		h.logger.Error("Failed to parse swipe input", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	swipe := model.Swipe{
		SwiperId:   userID,
		SwipedId:   input.SwipedId,
		Preference: input.Preference,
	}

	matched, matchID, err := h.storage.RecordSwipe(c.Context(), swipe)
	if err != nil {
		h.logger.Error("Failed to record swipe", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record swipe"})
	}

	result := fiber.Map{"matched": matched}
	if matched {
		result["matchID"] = matchID
	}

	h.logger.Info("Swipe recorded successfully", "swiperId", userID, "swipedId", input.SwipedId, "matched", matched)
	return c.JSON(fiber.Map{"results": result})
}
