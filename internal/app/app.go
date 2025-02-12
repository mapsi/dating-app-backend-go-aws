package app

import (
	"dating-app-backend/internal/config"
	"dating-app-backend/internal/handler"
	"dating-app-backend/internal/logger"
	"dating-app-backend/internal/middleware"
	"dating-app-backend/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	config  *config.Config
	storage *storage.DynamoDB
	fiber   *fiber.App
	logger  *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) (*App, error) {
	db, err := storage.NewDynamoDB(cfg, logger)
	if err != nil {
		return nil, err
	}

	app := &App{
		config:  cfg,
		storage: db,
		fiber:   fiber.New(),
		logger:  logger,
	}

	app.setupRoutes()

	return app, nil
}

func (a *App) setupRoutes() {
	userHandler := handler.NewUserHandler(a.storage, a.logger)
	authHandler := handler.NewAuthHandler(a.storage, a.logger)
	discoverHandler := handler.NewDiscoverHandler(a.storage, a.logger)
	swipeHandler := handler.NewSwipeHandler(a.storage, a.logger)

	authMiddleware := middleware.NewAuthMiddleware(a.config)

	a.fiber.Post("/user/create", userHandler.CreateRandomUser)
	a.fiber.Post("/login", authHandler.Login)

	// Protected routes
	a.fiber.Get("/discover", authMiddleware, discoverHandler.DiscoverUsers)
	a.fiber.Post("/swipe", authMiddleware, swipeHandler.RecordSwipe)

	a.logger.Info("Routes set up successfully")
}

func (a *App) Run() error {
	a.logger.Info("Starting application", "port", a.config.Port)
	return a.fiber.Listen(":" + a.config.Port)
}
