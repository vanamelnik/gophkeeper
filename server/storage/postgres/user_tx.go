package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

var _ storage.UserTransaction = (*UserTransaction)(nil)

type UserTransaction struct {
	tx *sql.Tx
}

// NewUserTransaction implements storage.Storage interface.
func (s Storage) NewUserTransaction(ctx context.Context, userID uuid.UUID) (UserTransaction, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return UserTransaction{}, err
	}

	return UserTransaction{tx: tx}, nil
}
