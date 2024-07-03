package handler

import (
	"dating-app-backend/internal/logger"
	"dating-app-backend/internal/model"
	"dating-app-backend/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	storage *storage.DynamoDB
	logger  *logger.Logger
}

func NewUserHandler(storage *storage.DynamoDB, logger *logger.Logger) *UserHandler {
	return &UserHandler{storage: storage, logger: logger}
}

func (h *UserHandler) CreateRandomUser(ctx *fiber.Ctx) error {
	user := model.GenerateRandomUser()

	h.logger.Info("Created user", "userId", user.ID)

	err := h.storage.CreateUser(ctx.Context(), user)
	if err != nil {
		msg := "Failed to store user"
		h.logger.Error(msg, "error", err, "userId", user.ID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": msg,
		})
	}

	h.logger.Info("Stored user", "userId", user.ID)
	return ctx.JSON(fiber.Map{
		"result": user,
	})
}
