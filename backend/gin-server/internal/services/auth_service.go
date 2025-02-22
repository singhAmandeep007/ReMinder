package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"gin-server/internal/config"
	"gin-server/internal/models"
	"gin-server/internal/repositories"
	"gin-server/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      config.Config
}

func NewAuthService(userRepo repositories.UserRepository, cfg config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, email, username, password string) (*models.User, error) {
	existingUser, _ := s.userRepo.GetUserByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) GetUserFromToken(token *jwt.Token) (*models.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	userIDStr, ok := claims["userId"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format in token")
	}
	user, err := s.userRepo.GetUserByID(context.Background(), userID.String())
	if err != nil {
		return nil, errors.New("user not found from token")
	}
	return user, nil
}
