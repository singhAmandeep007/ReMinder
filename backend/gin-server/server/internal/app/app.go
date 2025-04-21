package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/db"
)

// App represents the application
type App struct {
	httpServer *http.Server
	router     *gin.Engine
	dbManager  *db.DBManager
	log        *logger.Logger
	cfg        *config.Config
}

// New creates a new App instance
func New(cfg *config.Config, log *logger.Logger) (*App, error) {

	// Initialize database
	dbManager, err := db.NewDBManager(cfg, log)
	if err != nil {
		return nil, err
	}

	// Create dependency container
	container := NewContainer(cfg, log, dbManager)

	// Initialize router with dependency container
	router := NewRouter(container)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}

	return &App{
		httpServer: server,
		router:     router,
		dbManager:  dbManager,
		log:        log,
		cfg:        cfg,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Connect to the database
	if err := a.dbManager.DB.Connect(ctx); err != nil {
		a.log.Errorf("Failed to connect to database: %v", err)
	}

	err := a.dbManager.DB.Migrate(ctx)
	if err != nil {
		a.log.Errorf("Failed to migrate database: %v", err)
	}

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		a.log.Infof("Starting server at port %d", a.cfg.Port)
		serverErrors <- a.httpServer.ListenAndServe()
	}()
	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err

	case <-shutdown:
		a.log.Infof("Starting graceful shutdown...")

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Gracefully shutdown the server by waiting on existing requests
		if err := a.httpServer.Shutdown(ctx); err != nil {
			// If shutdown timed out, force close
			a.log.Errorf("Graceful shutdown timed out error %s", err)
			a.httpServer.Close()
			return err
		}

		a.log.Infof("Server gracefully stopped")
	}

	return nil
}

// Add cleanup method to handle resource cleanup
func (a *App) Cleanup() error {
	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close database connections
	if err := a.dbManager.DB.Close(ctx); err != nil {
		a.log.Errorf("Error closing database: %v", err)
	}

	// Shutdown HTTP server
	return a.httpServer.Shutdown(ctx)
}
