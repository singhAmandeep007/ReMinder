package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/auth"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/usernamegen"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/db"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/domain"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/repository"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (newAccessToken, newRefreshToken string, error error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, error error)
	GetMe(ctx context.Context, userId string) (*domain.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	log         *logger.Logger
	authManager *auth.AuthManager
}

func NewAuthService(userRepo repository.UserRepository, log *logger.Logger, authManager *auth.AuthManager) AuthService {
	return &authService{
		userRepo:    userRepo,
		log:         log,
		authManager: authManager,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*domain.User, error) {

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil && err == nil {
		return nil, domain.ErrUserAlreadyExist
	}

	if err != nil && err != db.ErrNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user
	err = s.userRepo.Create(ctx, &domain.User{
		Email:    email,
		Password: hashedPassword,
		Username: usernamegen.Generate(),
		ID:       uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	// Validate user credentials
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	isPasswordValid := utils.VerifyPassword(user.Password, password)

	if !isPasswordValid {
		return "", "", domain.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := s.authManager.GenerateTokenPair(user.ID, map[string]interface{}{
		"email":    user.Email,
		"role":     user.Role,
		"username": user.Username,
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token
	newAccessToken, newRefreshToken, err := s.authManager.RefreshTokens(refreshToken)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// GetMe retrieves the user details for the given user ID
func (s *authService) GetMe(ctx context.Context, userId string) (*domain.User, error) {
	// Fetch user details
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
