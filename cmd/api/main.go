package main

import (
	"dating-app-backend/internal/app"
	"dating-app-backend/internal/config"
	"dating-app-backend/internal/logger"
	"os"
)

func main() {
	log := logger.NewLogger()
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	log.Info("Config loaded successfully", "config", cfg)

	app, err := app.New(cfg, log)
	if err != nil {
		log.Error("Failed to create app", "error", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		log.Error("Failed to run app: %v", err)
		os.Exit(1)
	}

}
