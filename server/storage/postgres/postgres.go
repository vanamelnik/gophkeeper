package postgres

import (
	"database/sql"
	_ "embed"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

var _ storage.Storage = (*Storage)(nil)

// Storage implements storage.Storage interface for Postgresql database engine.
type Storage struct {
	db *sql.DB
}

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
