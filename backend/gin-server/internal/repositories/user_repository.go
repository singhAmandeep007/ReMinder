package repositories

import (
	"context"
	"fmt"

	"gin-server/internal/db"
	"gin-server/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	dbManager *db.DBManager
}

func NewUserRepository(dbManager *db.DBManager) UserRepository {
	return &userRepository{dbManager: dbManager}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	err := r.dbManager.Create(ctx, user)
	return err
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	condition := map[string]interface{}{"username": username}
	err := r.dbManager.FindByCondition(ctx, &user, condition)
	if err != nil {
		return nil, fmt.Errorf("user not found with username: %s, error: %w", username, err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.dbManager.FindByID(ctx, &user, id)
	if err != nil {
		return nil, fmt.Errorf("user not found with ID: %s, error: %w", id, err)
	}
	return &user, nil
}
