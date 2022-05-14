package models

// Package models defines canonical models for representation
// such objects as User, Session, Event and all kind of Items.

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User represents the user of the service.
// One User can have multiple sessions from different devices.
type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	DeletedAt    time.Time
}

// Session represents a single client session of the user with given ID.
// For each session a pair of AccessToken + RefreshToken is created.
// If AccessToken expires, RefreshToken will be the key to create a new token pair.
type Session struct {
	ID           uint
	UserID       uuid.UUID
	RefreshToken RefreshToken
	LoginAt      time.Time
	LogoutAt     time.Time
}

// RefreshToken is a random string that used for check if the service may generate
// a new access token to replace an expired one.
type RefreshToken string

// UserData represents the snapshot of all user's data stored in the storage.
// Each snapshot has an unique version number that is incremented when the data
// is updated on the server.
type UserData struct {
	Version uint64
	Items   []Item
}

// Errors
var (
	ErrInvalidItem = errors.New("invalid item type")
)
