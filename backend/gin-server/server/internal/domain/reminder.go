package domain

import "time"

type Reminder struct {
	ID          string `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	IsPinned    bool   `json:"isPinned" db:"is_pinned"`
	UserID      string `json:"userId" db:"user_id"` // Foreign Key to User
	// the omitempty tag is used in struct field tags to control how the field is handled during JSON serialization (when converting a struct to JSON using the encoding/json package). Specifically, it tells the JSON encoder to omit the field from the resulting JSON if the field has its zero value. For strings: "" (empty string) For integers: 0 For booleans: false For pointers, slices, maps, and interfaces: nil
	ReminderGroupID string    `json:"reminderGroupId,omitempty" db:"reminder_group_id"` // Foreign Key to ReminderGroup
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}
