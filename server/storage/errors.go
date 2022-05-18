package storage

import "errors"

var (
	ErrAlreadyExists = errors.New("entry already exists")
	ErrNotFound      = errors.New("entry not found")
)
