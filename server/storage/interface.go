package storage

// Package storage contains tools for interacting with the database.

import (
	"context"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

type (
	// Storage is an interface that wraps methods for interacting with the database.
	Storage interface {
		// CreateUser creates a new record with user's data in the database.
		// It returns models.User object with generated user ID and CreatedAt field.
		// If user with such email is already exists, the erroro ErrAlreadyExists returns.
		CreateUser(ctx context.Context, email, passwordHash string) (models.User, error)

		// TODO:
		// // UpdateUser updates information of the user with ID provided.
		// // All fields must be filled.
		// UpdateUser(ctx context.Context, user models.User) error
		// // DeleteUser removes the user provided from the database.
		// DeleteUser(ctx context.Context, userID uuid.UUID) error

		// GetUserByEmail finds the user with given email.
		GetUserByEmail(ctx context.Context, email string) (models.User, error)
		// GetUserDataVersion returns current user data version.
		GetUserDataVersion(ctx context.Context, userID uuid.UUID) (uint64, error)

		// CreateSession creates a record in the database for the new
		// session. All data should be validated by the caller.
		CreateSession(ctx context.Context, session models.Session) error
		// UpdateSession updates session information for the session with ID provided.
		// Fields UserID and LoginAt is ignored.
		UpdateSession(ctx context.Context, session models.Session) error
		// GetSessionByID returns the active session with ID prvided.
		// If the session is logged out, ErrNotFound returns.
		GetSessionByID(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
		// GetActiveUserSessions returns all active sessions of the specified user.
		// If there are no active sessions, an empy list is returned and a nil error.
		GetActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

		// TODO:
		// // LogoutAll marks all session of the user provided as logged out.
		// LogoutAll(ctx context.Context, userID uuid.UUID) error

		// NewUserTransaction starts a new transaction that implements the specified user's events.
		NewUserTransaction(ctx context.Context, userID uuid.UUID) (UserTransaction, error)

		// GetUserData returns all user's items: passwords, blobs, texts and cards.
		GetUserData(ctx context.Context, userID uuid.UUID) (*models.UserData, error)
	}

	// UserTransaction is an interface that wraps methods that performs user events committing.
	// Each transaction must be closed by calling Commit or RollBack method.
	UserTransaction interface {
		// CreateItem adds a new record in the database.
		CreateItem(ctx context.Context, item models.Item) error

		// UpdateItem updates the record in the database.
		UpdateItem(ctx context.Context, item models.Item) error

		// RollBack cancels the transaction if it's not closed yet.
		RollBack() error

		// Commit closes the transaction and commits all changes.
		// Also it increments DataVersion field of the user.
		Commit() error
	}
)
