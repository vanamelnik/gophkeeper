package gophkeeper

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/storage"
)

type (
	Service struct {
		db storage.Storage
		dl DataLoader
	}
)

var (
	ErrVersionUpToDate = errors.New("data version is already up to date")
)

func NewGophkeeper(db storage.Storage) Service {
	dl := NewDataLoader()
	return Service{
		db: db,
		dl: dl,
	}
}

func (s Service) GetUserData(ctx context.Context, userID uuid.UUID, version uint64) (*models.UserData, error) {
	data, err := s.db.GetUserData(ctx, userID)
	if err != nil {
		return nil, err
	}
	if version == data.Version {
		return nil, ErrVersionUpToDate
	}
	return data, nil
}
