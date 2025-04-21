package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/constants"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
)

var (
	// ErrDatabaseConnection is returned when the database connection fails
	ErrDatabaseConnection = errors.New("database connection error")
	// ErrNotFound is returned when a requested entity is not found
	ErrNotFound = errors.New("entity not found")
	// ErrDuplicate is returned when an entity already exists
	ErrDuplicate = errors.New("entity already exists")
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
	// ErrInternal is returned for internal database errors
	ErrInternal = errors.New("internal database error")
	// ErrNotImplemented is returned when a method is not implemented
	ErrNotImplemented = errors.New("method not implemented")
)

// Database defines the interface for database operations
type Database interface {
	// Connect establishes a connection to the database
	Connect(ctx context.Context) error

	// Close closes the database connection
	Close(ctx context.Context) error

	// Ping checks if the database is accessible
	Ping(ctx context.Context) error

	// Migrate runs database migrations
	Migrate(ctx context.Context) error

	// Seed populates the database with initial data
	Seed(ctx context.Context) error

	// Collection returns a collection/table handler for the given name
	Collection(name string) Collection
}

// Collection defines the interface for collection/table operations
type Collection interface {
	// Inserts a new document/record
	Create(ctx context.Context, data interface{}) (string, error)

	// Retrieves a document/record by Id
	GetById(ctx context.Context, id string, result interface{}) error

	GetOne(ctx context.Context, filter map[string]interface{}, result interface{}) error

	// Retrieves documents/records matching the filter
	GetAllByCondition(ctx context.Context, filter map[string]interface{}, results interface{}) error

	// Updates a document/record by ID
	UpdateById(ctx context.Context, id string, data interface{}) error

	// Removes a document/record by ID
	DeleteById(ctx context.Context, id string) error

	// Returns the number of documents/records matching the filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// DBManager will hold the active database connection.
type DBManager struct {
	DB Database
}

// Factory creates a specific database implementation based on the config
func NewDBManager(cfg *config.Config, logger *logger.Logger) (*DBManager, error) {
	var db Database
	var err error

	// Import database drivers dynamically to avoid direct dependencies
	switch cfg.DBType {
	case constants.SQLite:
		db, err = NewSQLiteDatabase(cfg, logger)
	case constants.Firestore:
		db, err = NewFirestoreDatabase(cfg, logger)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	if err != nil {
		return nil, err
	}

	// if err = db.Connect(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to connect to database: %w", err)
	// }

	return &DBManager{DB: db}, nil
}
