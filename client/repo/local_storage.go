package repo

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

// Package local contains tools for storing user data locally on the client side.

type (
	// Item wraps user data types described in models package with Pending flag.
	Item struct {
		Data models.Item
		// Pending is the flag indicating that the data was created locally
		// and is awaiting confirmation that it was successfully stored on the server.
		Pending bool
	}

	// Storage represents the local in-memory storage of user data.
	Storage struct {
		// DataVersion is the version of current user data snapshot.
		DataVersion uint64

		// TODO: написать сюда, зачем нужен мьютекс)))
		sync.Mutex
		// Items contains user data wrapped in Item struct (with Pending flag).
		items []Item

		// fileName is the name of the file in .gob format that contains
		// current snapshot of local user data.
		fileName string
		// flushInterval is the interval after which the file is updated if isChange flag is set.
		flushInterval time.Duration

		// isChanged indicates if local user data has been changed and needs to be stored in the file.
		isChanged bool
	}
)

// CreateItem stores user data to the local storage and marks it as 'pending'.
func (s *Storage) CreateItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	if _, err := s.GetItemByID(item.ID); err == nil {
		return ErrAlreadyExists
	}
	s.items = append(s.items, Item{
		Data:    item,
		Pending: true,
	})

	return nil
}

// UpdateItem updates the item in local storage and marks it as 'pending'.
func (s *Storage) UpdateItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	for i, storedItem := range s.items {
		if storedItem.Data.ID == item.ID && storedItem.Data.DeletedAt == nil { // Undelete is not possible
			s.items[i] = Item{
				Data:    item,
				Pending: true,
			}
			return nil
		}
	}
	return ErrNotFound
}

// DeleteItem marks the item in local storage as 'deleted' and 'pending'
func (s *Storage) DeleteItem(itemID uuid.UUID) error {
	s.Lock()
	defer s.Unlock()
	for i, storedItem := range s.items {
		if storedItem.Data.ID == itemID && storedItem.Data.DeletedAt == nil {
			now := time.Now()
			s.items[i] = Item{
				Data: models.Item{
					ID:        itemID,
					CreatedAt: storedItem.Data.CreatedAt,
					DeletedAt: &now,
					Data:      nil, // We erase all user data in deleted items,
					Meta:      "",  // because that's private data.
				},
				Pending: true,
			}
			return nil
		}
	}

	return nil
}

// GetItemByID fetches the non deleted item with given ID from local storage.
func (s *Storage) GetItemByID(itemID uuid.UUID) (Item, error) {
	s.Lock()
	defer s.Unlock()
	for _, item := range s.items {
		if item.Data.ID == itemID && item.Data.DeletedAt != nil {
			return item, nil
		}
	}
	return Item{}, ErrNotFound
}

// GetDataSnapshot retrieves the snapshot of all user data in local storage.
func (s *Storage) GetDataSnapshot() []Item {
	return s.items
}
