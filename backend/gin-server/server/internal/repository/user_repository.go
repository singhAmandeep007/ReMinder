package repository

import (
	"context"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/db"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetById(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	collection db.Collection
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *db.DBManager) UserRepository {
	// Get the database instance using the singleton pattern

	return &userRepository{
		collection: db.DB.Collection("users"),
	}
}

// Implementation of UserRepository interface
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.collection.Create(ctx, user)
	return err
}

func (r *userRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.collection.GetById(ctx, id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.GetOne(ctx, map[string]interface{}{"email": email}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
