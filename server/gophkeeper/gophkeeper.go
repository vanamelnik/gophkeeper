package gophkeeper

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

// Service represents the main service that implements the business logic of the server.
type Service struct {
	storage storage.Storage
	// dl DataLoader
}

var (
	ErrVersionUpToDate = errors.New("data version is already up to date")
)

func NewGophkeeper(db storage.Storage) Service {
	// dl := NewDataLoader()
	return Service{
		storage: db,
		// dl: dl,
	}
}

// GetUserData retrieves a snapshot of all user's items in storage.
// If the version of the data is equal to the version provided, the ErrVersionUpToDate is thrown.
func (s Service) GetUserData(ctx context.Context, userID uuid.UUID, version uint64) (*models.UserData, error) {
	data, err := s.storage.GetUserData(ctx, userID)
	if err != nil {
		return nil, err
	}
	if version == data.Version {
		return nil, ErrVersionUpToDate
	}
	return data, nil
}

// PublishUserData applies local changes of user data to the database.
func (s Service) PublishUserData(ctx context.Context, userID uuid.UUID, events []models.Event) error {
	tx, err := s.storage.NewUserTransaction(ctx, userID)
	if err != nil {
		return err
	}
	defer tx.RollBack()

	for _, event := range events {
		switch data := event.Item.Data.(type) {
		case models.TextData:
			err := processText(ctx, tx, event.Operation, data)
		}
	}

	return nil
}

func processText(ctx context.Context, tx storage.UserTransaction, op models.Operation, text models.TextData) error {
	switch op {
	case models.OpCreate:
	case models.OpUpdate:
	case models.OpDelete:
	}
	return nil
}
func processPassword(ctx context.Context, tx storage.UserTransaction, op models.Operation) error {
	switch op {
	case models.OpCreate:
	case models.OpUpdate:
	case models.OpDelete:
	}
	return nil
}
func processBinary(ctx context.Context, tx storage.UserTransaction, op models.Operation) error {
	switch op {
	case models.OpCreate:
	case models.OpUpdate:
	case models.OpDelete:
	}
	return nil
}
func processCard(ctx context.Context, tx storage.UserTransaction, op models.Operation) error {
	switch op {
	case models.OpCreate:
	case models.OpUpdate:
	case models.OpDelete:
	}
	return nil
}
