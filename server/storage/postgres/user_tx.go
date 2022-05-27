package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

// CreateItem implements storage.UserTransaction interface.
func (tx *UserTransaction) CreateItem(ctx context.Context, item models.Item) error {
	switch data := item.Payload.(type) {
	case models.TextData:
		return tx.createText(ctx, item, data)
	case models.BinaryData:
		return tx.createBlob(ctx, item, data)
	case models.PasswordData:
		return tx.createPassword(ctx, item, data)
	case models.CardData:
		return tx.createCard(ctx, item, data)
	}

	return errors.New("unreachable error: wrong item payload type")
}

// UpdateItem implements storage.UserTransaction interface.
func (tx *UserTransaction) UpdateItem(ctx context.Context, item models.Item) error {
	switch data := item.Payload.(type) {
	case models.TextData:
		return tx.updateText(ctx, item, data)
	case models.BinaryData:
		return tx.updateBlob(ctx, item, data)
	case models.PasswordData:
		return tx.updatePassword(ctx, item, data)
	case models.CardData:
		return tx.updateCard(ctx, item, data)
	}

	return errors.New("unreachable error: wrong item payload type")
}

// RollBack implements storage.UserTransaction interface.
func (tx *UserTransaction) RollBack() error {
	return tx.RollBack()
}

// Commit implements storage.UserTransaction interface.
func (tx *UserTransaction) Commit() error {
	_, err := tx.Exec(`UPDATE users SET data_version = data_version+1 WHERE id=$1;`, tx.userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	return tx.Commit()
}

// createText adds a new text item into the texts table. Item version is set to 1.
func (tx *UserTransaction) createText(ctx context.Context, item models.Item, data models.TextData) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO texts (id, user_id, version, meta, created_at,text_string)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		item.ID, tx.userID, 1, item.Meta, item.CreatedAt, data.Text,
	)
	if err != nil {
		return err
	}

	return nil
}

// createCard adds a new card item into the cards table. Item version is set to 1.
func (tx *UserTransaction) createCard(ctx context.Context, item models.Item, data models.CardData) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO cards (id, user_id, version, meta, created_at, card_number, cardholder_name, expiration_date, cvc)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		item.ID, tx.userID, 1, item.Meta, item.CreatedAt, data.Number, data.CardholderName, data.Date, data.CVC,
	)
	if err != nil {
		return err
	}

	return nil
}

// createBlob adds a new blob item into the blobs table. Item version is set to 1.
func (tx *UserTransaction) createBlob(ctx context.Context, item models.Item, data models.BinaryData) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO blobs (id, user_id, version, meta, created_at, blob)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		item.ID, tx.userID, 1, item.Meta, item.CreatedAt, data.Binary,
	)
	if err != nil {
		return err
	}

	return nil
}

// createPassword adds a new password item into the passwords table. Item version is set to 1.
func (tx *UserTransaction) createPassword(ctx context.Context, item models.Item, data models.PasswordData) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO passwords (id, user_id, version, meta, created_at, password)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		item.ID, tx.userID, 1, item.Meta, item.CreatedAt, data.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

// updateText updates an existing text item in the texts table.
func (tx *UserTransaction) updateText(ctx context.Context, item models.Item, data models.TextData) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE texts SET version=$1, meta=$2, deleted_at=$3, text_string=$4 WHERE id=$5;`,
		item.Version+1, // This increments the item version!
		item.Meta,
		item.DeletedAt,
		data.Text,
		item.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	return nil
}

// updatePassword updates an existing password item in the passwords table.
func (tx *UserTransaction) updatePassword(ctx context.Context, item models.Item, data models.PasswordData) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE passwords SET version=$1, meta=$2, deleted_at=$3, ###=$4 WHERE id=$5;`,
		item.Version+1, // This increments the item version!
		item.Meta,
		item.DeletedAt,
		data.Password,
		item.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	return nil
}

// updateCard updates an existing card item in the cards table.
func (tx *UserTransaction) updateCard(ctx context.Context, item models.Item, data models.CardData) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE cards
		SET version=$1, meta=$2, deleted_at=$3, card_number=$4, cardholder_name=$5, expiration_date=$6, cvc=$7
		WHERE id=$8;`,
		item.Version+1, // This increments the item version!
		item.Meta,
		item.DeletedAt,
		data.Number,
		data.CardholderName,
		data.Date,
		data.CVC,
		item.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	return nil
}

// updateBlob updates an existing blob item in the blobs table.
func (tx *UserTransaction) updateBlob(ctx context.Context, item models.Item, data models.BinaryData) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE blobs SET version=$1, meta=$2, deleted_at=$3, blob=$4 WHERE id=$5;`,
		item.Version+1, // This increments the item version!
		item.Meta,
		item.DeletedAt,
		data.Binary,
		item.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	return nil
}
