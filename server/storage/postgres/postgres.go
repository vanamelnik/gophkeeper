package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

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
		tx     *sql.Tx
		userID uuid.UUID
	}

	PostgresOption func(s Storage) error
)

//go:embed schema.sql
var queryCreate string

// NewStorage creates a new Postgres storage and creates db schema if not exists.
func NewStorage(connStr string, opts ...PostgresOption) (Storage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return Storage{}, fmt.Errorf("newStorage: %w", err)
	}
	if err := db.Ping(); err != nil {
		return Storage{}, fmt.Errorf("newStorage: %w", err)
	}

	s := Storage{db}
	for _, optFn := range opts {
		if err := optFn(s); err != nil {
			return Storage{}, fmt.Errorf("newStorage: %w", err)
		}
	}

	// create database schema if not exists
	if _, err := db.Exec(queryCreate); err != nil {
		return Storage{}, fmt.Errorf("newStorage: %w", err)
	}

	return s, nil
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
		tx:     tx,
		userID: userID,
	}, nil
}

// WithDestructiveReset erases all tables in the DB!
func WithDesctructiveReset() PostgresOption {
	return func(s Storage) error {
		_, err := s.db.Exec(`DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;`)
		return err
	}
}
