package gophkeeper

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

type (
	// Service represents the main service that implements the business logic of the server.
	Service struct {
		ctx     context.Context
		storage storage.Storage

		eventCh chan eventsPack
		stopCh  chan struct{}
		wg      *sync.WaitGroup
	}

	// eventsPack is the pack of events received from the client.
	eventsPack struct {
		userID uuid.UUID
		events []models.Event
	}
)

var (
	ErrVersionUpToDate = errors.New("data version is up to date")
)

func NewService(ctx context.Context, db storage.Storage) Service {
	s := Service{
		ctx:     ctx,
		storage: db,
		eventCh: make(chan eventsPack, 1),
		stopCh:  make(chan struct{}),
		wg:      &sync.WaitGroup{},
	}

	go s.processor() // TODO: implement a worker pool to limit DB connections
	return s
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
func (s Service) PublishUserData(ctx context.Context, userID uuid.UUID, events []models.Event) {
	s.wg.Add(1)
	s.eventCh <- eventsPack{
		userID: userID,
		events: events,
	}
}

func (s Service) Close() {
	s.wg.Wait()
	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}
}
