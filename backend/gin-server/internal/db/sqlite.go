package db

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDatabase struct {
	DB       *gorm.DB
	FilePath string
}

func NewSQLiteDatabase(filePath string) (*SQLiteDatabase, error) {
	return &SQLiteDatabase{FilePath: filePath}, nil
}

func (s *SQLiteDatabase) Connect() error {
	db, err := gorm.Open(sqlite.Open(s.FilePath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to sqlite: %w", ErrDatabaseConnection)
	}
	s.DB = db
	return nil
}

func (s *SQLiteDatabase) Ping() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db from gorm db: %w", err)
	}
	return sqlDB.Ping()
}

func (s *SQLiteDatabase) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db from gorm db: %w", err)
	}
	return sqlDB.Close()
}

func (s *SQLiteDatabase) Migrate(models ...interface{}) error {
	return s.DB.AutoMigrate(models...)
}

func (s *SQLiteDatabase) Create(ctx context.Context, model interface{}) error {
	result := s.DB.WithContext(ctx).Create(model)
	return result.Error
}

func (s *SQLiteDatabase) FindByID(ctx context.Context, model interface{}, id string) error {
	result := s.DB.WithContext(ctx).First(model, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return result.Error
}

func (s *SQLiteDatabase) FindByCondition(ctx context.Context, model interface{}, condition map[string]interface{}) error {
	result := s.DB.WithContext(ctx).Where(condition).First(model)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return result.Error
}

func (s *SQLiteDatabase) FindAll(ctx context.Context, model interface{}) error {
	result := s.DB.WithContext(ctx).Find(model)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return result.Error
}
func (s *SQLiteDatabase) Update(ctx context.Context, model interface{}, id string, updates map[string]interface{}) error {
	result := s.DB.WithContext(ctx).Model(model).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (s *SQLiteDatabase) Delete(ctx context.Context, model interface{}, id string) error {
	result := s.DB.WithContext(ctx).Delete(model, "id = ?", id)
	return result.Error
}
