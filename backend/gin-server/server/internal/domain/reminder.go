package domain

import "time"

// ReminderState represents the current state of a reminder
type ReminderState string

const (
	ReminderStateInactive   ReminderState = "INACTIVE"
	ReminderStateInProgress ReminderState = "INPROGRESS"
	ReminderStateCompleted  ReminderState = "COMPLETED"
	ReminderStateArchived   ReminderState = "ARCHIVED"
)

// ReminderFocusSession represents a focus session for a reminder
type ReminderFocusSession struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// Reminder represents a reminder/task item
// Fields are aligned with OpenAPI specification in /openapi/api.yaml
type Reminder struct {
	ID                  string                 `json:"id" db:"id"`
	Title               string                 `json:"title" db:"title"`
	Description         string                 `json:"description,omitempty" db:"description"`
	State               ReminderState          `json:"state" db:"state"`
	IsPinned            bool                   `json:"isPinned" db:"is_pinned"`
	UserID              string                 `json:"userId" db:"user_id"`
	ReminderGroupID     *string                `json:"reminderGroupId,omitempty" db:"reminder_group_id"`
	DueDate             *time.Time             `json:"dueDate,omitempty" db:"due_date"`
	RepeatTimes         *int                   `json:"repeatTimes,omitempty" db:"repeat_times"`
	RepeatInterval      *int64                 `json:"repeatInterval,omitempty" db:"repeat_interval"`
	CurrentFocusSession *ReminderFocusSession  `json:"currentFocusSession,omitempty" db:"-"`
	FocusSessions       []ReminderFocusSession `json:"focusSessions" db:"-"`
	CreatedAt           time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time              `json:"updatedAt" db:"updated_at"`
}

// ReminderWithGroup represents a reminder with its associated group
type ReminderWithGroup struct {
	Reminder
	Group *ReminderGroup `json:"group,omitempty"`
}
