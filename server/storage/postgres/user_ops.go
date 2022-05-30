package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

// CreateUser implements storage.Storage interface.
func (s Storage) CreateUser(ctx context.Context, email, passwordHash string) (models.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return models.User{}, err
	}
	now := time.Now()
	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO users
		(id, email, password_hash, data_version, created_at)
		VALUES ($1, $2, $3, $4, $5);`,
		id,
		email,
		passwordHash,
		0,
		now,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return models.User{}, storage.ErrAlreadyExists
		}
		return models.User{}, err
	}

	return models.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    &now,
		DeletedAt:    nil,
	}, nil
}

// GetUserByEmail implements storage.Storage interface.
func (s Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	u := models.User{Email: email}
	err := s.db.QueryRowContext(ctx, `SELECT id, password_hash, created_at FROM users WHERE email=$1 AND deleted_at IS NULL;`, email).
		Scan(&u.ID, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

// GetUserDataVersion implements storage.Storage interface.
func (s Storage) GetUserDataVersion(ctx context.Context, userID uuid.UUID) (uint64, error) {
	var dataVersion uint64
	err := s.db.QueryRowContext(ctx, `SELECT data_version FROM users WHERE id=$1 AND deleted_at IS NULL;`, userID).Scan(&dataVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrNotFound
		}
		return 0, err
	}
	return dataVersion, nil
}
