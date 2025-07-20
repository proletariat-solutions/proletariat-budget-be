package domain

import (
	"errors"
	"time"
)

// User and role domain errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserHasActiveRoles = errors.New("user has active roles and cannot be deleted")
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleInUse          = errors.New("role is in use and cannot be deleted")
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    string    `json:"userId"`
}
