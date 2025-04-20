package service

import (
	"context"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/domain"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/repository"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, userID string, data domain.Reminder) (*domain.Reminder, error)
	ListRemindersByUserID(ctx context.Context, userID string) ([]domain.Reminder, error)
}

type reminderService struct {
	reminderRepo repository.ReminderRepository
	log          *logger.Logger
}

func NewReminderService(reminderRepo repository.ReminderRepository, log *logger.Logger) ReminderService {
	return &reminderService{reminderRepo: reminderRepo, log: log}
}

func (s *reminderService) CreateReminder(ctx context.Context, userID string, data domain.Reminder) (*domain.Reminder, error) {
	reminder := &domain.Reminder{
		UserID: userID,
		// data
	}
	if err := s.reminderRepo.Create(ctx, reminder); err != nil {
		return nil, err
	}
	return reminder, nil
}

func (s *reminderService) ListRemindersByUserID(ctx context.Context, userID string) ([]domain.Reminder, error) {
	reminders, err := s.reminderRepo.GetAllByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}
