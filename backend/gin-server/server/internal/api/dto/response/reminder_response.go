package response

import "github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/domain"

// ReminderResponse wraps a single reminder in a data field
type ReminderResponse struct {
	Data *domain.ReminderWithGroup `json:"data"`
}

// RemindersResponse wraps an array of reminders in a data field
type RemindersResponse struct {
	Data []domain.ReminderWithGroup `json:"data"`
}

// ReminderGroupResponse wraps a single reminder group in a data field
type ReminderGroupResponse struct {
	Data *domain.ReminderGroup `json:"data"`
}

// ReminderGroupsResponse wraps an array of reminder groups in a data field
type ReminderGroupsResponse struct {
	Data []domain.ReminderGroup `json:"data"`
}

// MessageResponse represents a success message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}
