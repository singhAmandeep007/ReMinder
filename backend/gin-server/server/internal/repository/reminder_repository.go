package repository

import (
	"context"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/db"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/domain"
)

type ReminderRepository interface {
	Create(ctx context.Context, user *domain.Reminder) error
	GetById(ctx context.Context, id string) (*domain.Reminder, error)
	GetAllByUserId(ctx context.Context, userId string) ([]domain.Reminder, error)
}

type reminderRepository struct {
	collection db.Collection
}

// NewReminderRepository creates a new instance of ReminderRepository
func NewReminderRepository(db *db.DBManager) ReminderRepository {
	// Get the database instance using the singleton pattern

	return &reminderRepository{
		collection: db.DB.Collection("reminders"),
	}
}

// Implementation of ReminderRepository interface
func (r *reminderRepository) Create(ctx context.Context, reminder *domain.Reminder) error {
	_, err := r.collection.Create(ctx, reminder)
	return err
}

func (r *reminderRepository) GetById(ctx context.Context, id string) (*domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.collection.GetById(ctx, id, &reminder)
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) GetAllByUserId(ctx context.Context, userId string) ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.collection.GetAllByCondition(ctx, map[string]interface{}{"user_id": userId}, &reminders)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}
