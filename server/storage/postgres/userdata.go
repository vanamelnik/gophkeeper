package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

// GetUserData implements storage.Storage interface.
func (s Storage) GetUserData(ctx context.Context, userID uuid.UUID) (*models.UserData, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	// nolint: errcheck
	defer tx.Rollback()

	userData := models.UserData{}
	// get DataVersion
	if err := tx.QueryRowContext(ctx, `SELECT data_version FROM users WHERE id=$1 AND deleted_at NOT NULL;`,
		userID).Scan(&userData.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	// get user data of all types
	texts, err := s.getTexts(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	passwords, err := s.getPasswords(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	cards, err := s.getCards(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	blobs, err := s.getBlobs(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	userData.Items = append(userData.Items, texts...)
	userData.Items = append(userData.Items, passwords...)
	userData.Items = append(userData.Items, blobs...)
	userData.Items = append(userData.Items, cards...)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &userData, nil
}

// getTexts retrieves from the database all texts of the user provided.
func (s Storage) getTexts(ctx context.Context, tx *sql.Tx, userID uuid.UUID) ([]models.Item, error) {
	items := make([]models.Item, 0)
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, text_string, meta, created_at, deleted_at FROM texts WHERE user_id=$1;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		data := models.TextData{}
		item := models.Item{}
		if err := rows.Scan(&item.ID, &data.Text, &item.Meta, &item.CreatedAt, &item.DeletedAt); err != nil {
			return nil, err
		}
		item.Payload = data
		items = append(items, item)
	}
	return items, nil
}

// getPasswords retrieves from the database all passwords of the user provided.
func (s Storage) getPasswords(ctx context.Context, tx *sql.Tx, userID uuid.UUID) ([]models.Item, error) {
	items := make([]models.Item, 0)
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, password, meta, created_at, deleted_at FROM passwords WHERE user_id=$1;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		data := models.PasswordData{}
		item := models.Item{}
		if err := rows.Scan(&item.ID, &data.Password, &item.Meta, &item.CreatedAt, &item.DeletedAt); err != nil {
			return nil, err
		}
		item.Payload = data
		items = append(items, item)
	}
	return items, nil
}

// getCards retrieves from the database all credit cards of the user provided.
func (s Storage) getCards(ctx context.Context, tx *sql.Tx, userID uuid.UUID) ([]models.Item, error) {
	items := make([]models.Item, 0)
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, card_number, cardholder_name, expiration_date, cvc, meta, created_at, deleted_at
		FROM cards WHERE user_id=$1;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		data := models.CardData{}
		item := models.Item{}
		if err := rows.Scan(&item.ID, &data.Number, &data.CardholderName, &data.Date, &data.CVC,
			&item.Meta, &item.CreatedAt, &item.DeletedAt); err != nil {
			return nil, err
		}
		item.Payload = data
		items = append(items, item)
	}
	return items, nil
}

// getBlobs retrieves from the database all blobs of the user provided.
func (s Storage) getBlobs(ctx context.Context, tx *sql.Tx, userID uuid.UUID) ([]models.Item, error) {
	items := make([]models.Item, 0)
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, blob, meta, created_at, deleted_at FROM blobs WHERE user_id=$1;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		data := models.BinaryData{}
		item := models.Item{}
		if err := rows.Scan(&item.ID, &data.Binary, &item.Meta, &item.CreatedAt, &item.DeletedAt); err != nil {
			return nil, err
		}
		item.Payload = data
		items = append(items, item)
	}
	return items, nil
}
