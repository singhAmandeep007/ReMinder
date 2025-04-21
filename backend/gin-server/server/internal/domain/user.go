package domain

import (
	"time"

	"errors"
)

const (
	// UserRoleAdmin represents an admin user
	UserRoleAdmin = "admin"
	// UserRoleUser represents a regular user
	UserRoleUser = "user"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" db:"id" firestore:"id"`
	Username  string    `json:"username" db:"username" firestore:"username"`
	Password  string    `json:"password" db:"password" firestore:"password"`
	Role      string    `json:"role" db:"role" firestore:"role"`
	Email     string    `json:"email" db:"email" firestore:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at" firestore:"updated_at"`
}

var ErrUserAlreadyExist = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
