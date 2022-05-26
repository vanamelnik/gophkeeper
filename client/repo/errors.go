package repo

import "errors"

var (
	ErrNotFound      = errors.New("item not found or deleted")
	ErrAlreadyExists = errors.New("item already exists")
)
