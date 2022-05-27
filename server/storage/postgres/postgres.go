package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

var _ storage.Storage = (*Storage)(nil)
var _ storage.UserTransaction = (*UserTransaction)(nil)

type (
	// Storage implements storage.Storage interface for Postgresql database engine.
	Storage struct {
		db *sql.DB
	}

	// UserTransaction implements storage.UserTransaction interface.
	UserTransaction struct {
		*sql.Tx
		userID uuid.UUID
	}
)

//go:embed schema.sql
var queryCreate string

// NewStorage creates a new Postgres storage and creates db schema if not exists.
func NewStorage(connStr string) (Storage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return Storage{}, err
	}
	if err := db.Ping(); err != nil {
		return Storage{}, err
	}
	// create database schema if not exists
	if _, err := db.Exec(queryCreate); err != nil {
		return Storage{}, err
	}

	return Storage{db: db}, nil
}

func (s Storage) Close() error {
	return s.db.Close()
}

// NewUserTransaction implements storage.Storage interface.
func (s Storage) NewUserTransaction(ctx context.Context, userID uuid.UUID) (storage.UserTransaction, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	return &UserTransaction{
		Tx:     tx,
		userID: userID,
	}, nil
}
