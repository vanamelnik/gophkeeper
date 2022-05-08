package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

type (
	Storage interface {
		CreateUser(ctx context.Context, login, passwordHash string) (models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		DeleteUser(ctx context.Context, UserID uuid.UUID) error

		LoginNewSession(ctx context.Context, UserID uuid.UUID) (models.Session, error)
		UpdateSession(ctx context.Context, session models.Session) error
		LogoutSession(ctx context.Context, session models.Session) error

		NewUserTx(ctx context.Context, UserID uuid.UUID) (UserTransaction, error)
		GetUserData(ctx context.Context, UserID uuid.UUID) (*models.UserData, error)
	}

	UserTransaction interface {
		CreateText(ctx context.Context, item models.TextItem) error
		CreateBlob(ctx context.Context, item models.BinaryItem) error
		CreatePassword(ctx context.Context, item models.PasswordItem) error
		CreateCard(ctx context.Context, item models.CardItem) error

		UpdateText(ctx context.Context, item models.TextItem) error
		UpdateBlob(ctx context.Context, item models.BinaryItem) error
		UpdatePassword(ctx context.Context, item models.PasswordItem) error
		UpdateCard(ctx context.Context, item models.CardItem) error

		DeleteText(ctx context.Context, item models.TextItem) error
		DeleteBlob(ctx context.Context, item models.BinaryItem) error
		DeletePassword(ctx context.Context, item models.PasswordItem) error
		DeleteCard(ctx context.Context, item models.CardItem) error

		Close() error
	}
)
