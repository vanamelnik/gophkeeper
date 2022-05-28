package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

// CreateSession implements storage.Storage interface.
func (s Storage) CreateSession(ctx context.Context, session models.Session) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO sessions
		(id, user_id, refresh_token, login_at)
		VALUES ($1, $2, $3, $4);`,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.LoginAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSession implements storage.Storage interface.
func (s Storage) UpdateSession(ctx context.Context, session models.Session) error {
	_, err := s.db.ExecContext(
		ctx,
		`UPDATE sessions
		SET refresh_token=$1, logout_at=$2
		WHERE id=$3 AND logout_at IS NULL;`,
		session.RefreshToken,
		session.LogoutAt,
		session.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	return nil
}

// GetSessionByID implements storage.Storage interface.
func (s Storage) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session := models.Session{ID: sessionID}
	err := s.db.QueryRowContext(
		ctx,
		`SELECT user_id, refresh_token, login_at, logout_at
		FROM sessions
		WHERE id=$1 AND logout_at IS NULL`,
		sessionID,
	).Scan(&session.UserID, &session.RefreshToken, &session.LoginAt, &session.LogoutAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, storage.ErrNotFound
		}

		return models.Session{}, err
	}

	return session, nil
}

// GetActiveUserSessions implements storage.Storage interface.
func (s Storage) GetActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	sessions := make([]models.Session, 0)
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT (id, refresh_token, login_at, logout_at)
		FROM sessions
		WHERE user_id = $1 AND logout_at IS NULL`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		session := models.Session{UserID: userID}
		if err := rows.Scan(&session.ID, &session.RefreshToken, &session.LoginAt, &session.LogoutAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
