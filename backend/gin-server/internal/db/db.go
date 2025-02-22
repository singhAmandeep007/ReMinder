package db

import (
	"context"
	"errors"
	"fmt"

	"gin-server/internal/config"
)

// Custom errors
var (
	ErrDatabaseConnection = errors.New("database connection error")
	ErrRecordNotFound     = errors.New("record not found")
)

// Interface-Based Design:

// Database interface defines fundamental database operations common to all database types.
type Database interface {
	Connect() error
	Ping() error
	Close() error

	// SQL Database Operations
	Migrate(models ...interface{}) error

	// CRUD operations
	Create(ctx context.Context, model interface{}) error
	FindByID(ctx context.Context, model interface{}, id string) error
	FindByCondition(ctx context.Context, model interface{}, condition map[string]interface{}) error
	FindAll(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, model interface{}, id string) error
}

// DBManager will hold the active database connection.
type DBManager struct {
	DB Database
}

// NewDBManager initializes and returns a new DBManager based on the config.
func NewDBManager(cfg config.Config) (*DBManager, error) {
	var db Database
	var err error

	switch cfg.DBType {
	case "sqlite":
		db, err = NewSQLiteDatabase(cfg.SqliteFile)

	// case "postgres":
	// 	db, err = NewPostgresDatabase(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName)
	// case "mongodb":
	// 	db, err = NewMongoDBDatabase(cfg.MongoDBUri, cfg.MongoDBName)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err = db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DBManager{DB: db}, nil
}

// PingDatabase checks the database connection.
func (dm *DBManager) PingDatabase() error {
	return dm.DB.Ping()
}

// CloseDatabase closes the database connection.
func (dm *DBManager) CloseDatabase() error {
	return dm.DB.Close()
}

// Migrate runs the database migrations.
func (dm *DBManager) Migrate(models ...interface{}) error {
	return dm.DB.Migrate(models...)
}

// --- CRUD OPERATIONS ---

func (dm *DBManager) Create(ctx context.Context, model interface{}) error {
	return dm.DB.Create(ctx, model)
}

func (dm *DBManager) FindByID(ctx context.Context, model interface{}, id string) error {
	return dm.DB.FindByID(ctx, model, id)
}

func (dm *DBManager) FindByCondition(ctx context.Context, model interface{}, condition map[string]interface{}) error {
	return dm.DB.FindByCondition(ctx, model, condition)
}

func (dm *DBManager) FindAll(ctx context.Context, model interface{}) error {
	return dm.DB.FindAll(ctx, model)
}

func (dm *DBManager) Update(ctx context.Context, model interface{}, id string, updates map[string]interface{}) error {
	return dm.DB.Update(ctx, model, id, updates)
}

func (dm *DBManager) Delete(ctx context.Context, model interface{}, id string) error {
	return dm.DB.Delete(ctx, model, id)
}
