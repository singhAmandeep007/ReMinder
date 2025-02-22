package main

import (
	"fmt"
	"gin-server/internal/config"
	"gin-server/internal/db"
	"gin-server/internal/handlers"
	"gin-server/internal/logger"
	"gin-server/internal/models"
	"gin-server/internal/repositories"
	"gin-server/internal/router"
	"gin-server/internal/services"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	appLogger := logger.NewSimpleLogger()

	dbManager, err := db.NewDBManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	dbManager.Migrate(&models.User{})

	defer dbManager.CloseDatabase()

	userRepo := repositories.NewUserRepository(dbManager)

	authService := services.NewAuthService(userRepo, cfg)

	authHandler := handlers.NewAuthHandler(authService, appLogger)

	r := router.SetupRouter(authHandler)

	address := fmt.Sprintf(":%s", cfg.Port)
	appLogger.Infof("Server listening on %s", address)
	if err := r.Run(address); err != nil {
		appLogger.Errorf("Failed to run server: %v", err)
		log.Fatalf("Failed to run server: %v", err)
	}
}
