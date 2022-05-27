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
	ErrVersionUpToDate = errors.New("data version is up to date")
)

func NewGophkeeper(db storage.Storage) Service {
	// dl := NewDataLoader()
	return Service{
		storage: db,
		// dl: dl,
	}
}

// GetUserData returns the new items and the newer versions of existing local items from the database according to the version map provided.
func (s Service) GetUserData(ctx context.Context, userID uuid.UUID, versionMap map[uuid.UUID]uint64) (*models.UserData, error) {
	data, err := s.storage.GetUserData(ctx, userID)
	if err != nil {
		return nil, err
	}
	updates := models.UserData{
		Version: data.Version,
		Items:   make([]models.Item, 0, len(data.Items)),
	}
	for _, item := range data.Items {
		localVersion, ok := versionMap[item.ID]
		if !ok || // we have a new item
			localVersion != item.Version { // we have a newer version of this item
			updates.Items = append(updates.Items, item)
		}
	}

	return &updates, nil
}

// PublishUserData applies local changes of user data to the database.
func (s Service) PublishUserData(ctx context.Context, userID uuid.UUID, events []models.Event) error {
	tx, err := s.storage.NewUserTransaction(ctx, userID)
	if err != nil {
		return err
	}
	defer tx.RollBack()

	for _, event := range events {
		switch event.Operation {
		case models.OpCreate:
			if err := tx.CreateItem(ctx, event.Item); err != nil {
				return err
			}
		case models.OpUpdate:
			if err := tx.UpdateItem(ctx, event.Item); err != nil {
				return err
			}
		}
	}

	return nil
}
