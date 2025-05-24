package main

import (
	"context"
	"fmt"
	"os"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/bq"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/app"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load(utils.ResolvePathFromProjectRoot(".env.dev"))
	if err != nil {
		panic(fmt.Sprintf("could not load config: %v", err.Error()))
	}

	// Initialize logger
	log := logger.New(
		logger.WithServiceName("gin-server"),
		logger.WithDefaultDestinations(logger.FileLogger, logger.ConsoleLogger),
		logger.WithConsoleDestination(),
		logger.WithFileDestination(utils.ResolvePathFromProjectRoot("logs/gin-server.log"), 10, 5, 30, true),
		logger.WithMinLevel(logger.DebugLevel),
	)
	defer log.Close()

	// Initialize BigQuery client
	ctx := context.Background()

	bqConfig := bq.Config{
		ProjectID:           "flowing-castle-447017-h0",
		CredentialsPath:     utils.ResolvePathFromProjectRoot("flowing-castle-447017-h0-72697064ea74.json"),
		Location:            "US",
		QueryTimeoutSeconds: 300,
	}

	bqClient, err := bq.NewClient(ctx, bqConfig)
	if err != nil {
		log.Infof("Failed to create BigQuery client: %v", err)
	}
	defer bqClient.Close()

	// Initialize application
	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	// Ensure cleanup happens even on error
	defer application.Cleanup()

	// Start the server
	if err := application.Run(); err != nil {
		log.Fatalf("Server terminated with error: %v", err)
		os.Exit(1)
	}
}
