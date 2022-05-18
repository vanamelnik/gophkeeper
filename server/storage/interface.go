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
		// If user with such login is already exists, the erroro ErrAlreadyExists returns.
		CreateUser(ctx context.Context, login, passwordHash string) (models.User, error)
		// UpdateUser updates information of the user with ID provided.
		// All fields must be filled.
		UpdateUser(ctx context.Context, user models.User) error
		// DeleteUser removes the user provided from the database.
		DeleteUser(ctx context.Context, userID uuid.UUID) error
		// GetUserByLogin finds the user with given login.
		GetUserByLogin(ctx context.Context, login string) (models.User, error)

		// CreateSession creates a record in the database for the new
		// session. All data should be validated by the caller.
		CreateSession(ctx context.Context, session models.Session) error
		// UpdateSession updates session information for the session with ID provided.
		// Field UserID is ignored.
		UpdateSession(ctx context.Context, session models.Session) error
		// GetSessionByID returns the active session with ID prvided.
		// If the session is logged out, ErrNotFound returns.
		GetSessionByID(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
		// GetActiveUserSessions returns all active sessions of the specified user.
		// If there are no active sessions, an empy list is returned and a nil error.
		GetActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
		// LogoutAll marks all session of the user provided as logged out.
		LogoutAll(ctx context.Context, userID uuid.UUID) error

		// NewUserTransaction starts a new transaction that implements the specified user's events.
		NewUserTransaction(ctx context.Context, UserID uuid.UUID) (UserTransaction, error)

		// GetUserData returns all user's items: passwords, blobs, texts and cards.
		GetUserData(ctx context.Context, UserID uuid.UUID) (*models.UserData, error)
	}

	// UserTransaction is an interface that wraps methods that performs user events committing.
	// Each transaction must be closed by calling Commit or RollBack method.
	UserTransaction interface {
		// CreateText adds a new text record in the database.
		CreateText(ctx context.Context, item models.TextData) error
		// CreateBlob adds a new blob record in the database.
		CreateBlob(ctx context.Context, item models.BinaryData) error
		// CreatePassword adds a new password record in the database.
		CreatePassword(ctx context.Context, item models.PasswordData) error
		// CreateCard adds a new card record in the database.
		CreateCard(ctx context.Context, item models.CardData) error

		// UpdateText updates the text record in the database.
		UpdateText(ctx context.Context, item models.TextData) error
		// UpdateBlob updates the blob record in the database.
		UpdateBlob(ctx context.Context, item models.BinaryData) error
		// UpdatePassword updates the password record in the database.
		UpdatePassword(ctx context.Context, item models.PasswordData) error
		// UpdateCard updates the card record in the database.
		UpdateCard(ctx context.Context, item models.CardData) error

		// DeleteText deletes the text record in the database.
		DeleteText(ctx context.Context, item models.TextData) error
		// DeleteBlob deletes the blob record in the database.
		DeleteBlob(ctx context.Context, item models.BinaryData) error
		// DeletePassword deletes the password record in the database.
		DeletePassword(ctx context.Context, item models.PasswordData) error
		// DeleteCard deletes the card record in the database.
		DeleteCard(ctx context.Context, item models.CardData) error

		// RollBack cancels the transaction if it's not closed yet.
		RollBack() error
		// Commit closes the transaction and commits all changes.
		// Also it increments DataVersion field of the user.
		Commit() error
	}
)
