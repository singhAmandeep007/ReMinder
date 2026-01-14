package request

import "github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/domain"

// CreateReminderRequest represents the request body for creating a reminder
type CreateReminderRequest struct {
	Title   string  `json:"title" binding:"required,min=1"`
	GroupID *string `json:"groupId,omitempty"`
}

// UpdateReminderRequest represents the request body for updating a reminder
type UpdateReminderRequest struct {
	Title               *string                        `json:"title,omitempty" binding:"omitempty,min=1"`
	State               *domain.ReminderState          `json:"state,omitempty"`
	IsPinned            *bool                          `json:"isPinned,omitempty"`
	DueDate             *string                        `json:"dueDate,omitempty"`
	RepeatTimes         *int                           `json:"repeatTimes,omitempty"`
	RepeatInterval      *int64                         `json:"repeatInterval,omitempty"`
	CurrentFocusSession *domain.ReminderFocusSession   `json:"currentFocusSession,omitempty"`
	FocusSessions       *[]domain.ReminderFocusSession `json:"focusSessions,omitempty"`
	GroupID             *string                        `json:"groupId,omitempty"`
}

// CreateReminderGroupRequest represents the request body for creating a reminder group
type CreateReminderGroupRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

// UpdateReminderGroupRequest represents the request body for updating a reminder group
type UpdateReminderGroupRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1"`
}
