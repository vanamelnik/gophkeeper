package repo

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

// Package repo contains tools for storing user data locally on the client side.

type (
	// Entry wraps user data types described in models package with Pending flag.
	Entry struct {
		Item models.Item
		// Pending is the flag indicating that the data was created or updated locally
		// and is awaiting confirmation that the changes were successfully stored on the server.
		Pending bool
	}

	// Repo represents the local in-memory storage of user data.
	Repo struct {
		// dataVersion is the version of current user data snapshot.
		dataVersion uint64

		// TODO: написать сюда, зачем нужен мьютекс)))
		sync.RWMutex
		// Items contains user data wrapped in Item struct (with Pending flag).
		entries []Entry

		// Auth token pair could be stored in the .gob file with with the rest of the data.
		// If these fields are empty, this is interpreted as not having an active session.
		accessToken  models.AccessToken
		refreshToken models.RefreshToken

		// TODO: сделать сервис хранения репозитория в gob-файле:
		// // fileName is the name of the file in .gob format that contains
		// // current snapshot of local user data.
		// fileName string
		// // flushInterval is the interval after which the file is updated if isChange flag is set.
		// flushInterval time.Duration

		// isChanged indicates if local user data has been changed and needs to be stored in the file.
		isChanged bool
	}
)

func New() *Repo {
	return &Repo{
		dataVersion:  0,
		entries:      make([]Entry, 0),
		accessToken:  "",
		refreshToken: "",
		isChanged:    false,
	}
}

// CreateItem stores user data to the local repository and marks it as 'pending'.
func (r *Repo) CreateItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	r.RLock()
	defer r.RUnlock()
	if _, err := r.GetItemByID(item.ID); err == nil {
		return ErrAlreadyExists
	}
	r.entries = append(r.entries, Entry{
		Item:    item,
		Pending: true,
	})
	r.isChanged = true

	return nil
}

// UpdateItem updates the item in local repository and marks it as 'pending'.
func (r *Repo) UpdateItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	r.RLock()
	defer r.RUnlock()
	for i, storedItem := range r.entries {
		if storedItem.Item.ID == item.ID && storedItem.Item.DeletedAt == nil { // Undelete is not possible by this method.
			r.entries[i] = Entry{
				Item:    item,
				Pending: true,
			}

			r.isChanged = true
			return nil
		}
	}

	return ErrNotFound
}

// DeleteItem marks the item in local repository as 'deleted' and 'pending'
func (r *Repo) DeleteItem(itemID uuid.UUID) error {
	r.RLock()
	defer r.RUnlock()
	for i, storedItem := range r.entries {
		if storedItem.Item.ID == itemID && storedItem.Item.DeletedAt == nil {
			now := time.Now()
			r.entries[i] = Entry{
				Item: models.Item{
					ID:        itemID,
					CreatedAt: storedItem.Item.CreatedAt,
					DeletedAt: &now,
					Payload:   nil, // We erase all user data in deleted items,
					Meta:      "",  // because that's private data.
				},
				Pending: true,
			}
			r.isChanged = true
			return nil
		}
	}

	return nil
}

// GetItemByID fetches the non deleted item with given ID from local repository.
func (r *Repo) GetItemByID(itemID uuid.UUID) (Entry, error) {
	r.RLock()
	defer r.RUnlock()
	for _, item := range r.entries {
		if item.Item.ID == itemID && item.Item.DeletedAt != nil {
			return item, nil
		}
	}
	return Entry{}, ErrNotFound
}

// GetDataSnapshot retrieves the snapshot of all user data in local repository.
func (r *Repo) GetDataSnapshot() []Entry {
	return r.entries
}

// GetDataVersion returns current DataVersion of local user data.
func (r *Repo) GetDataVersion() uint64 {
	r.Lock()
	defer r.Unlock()
	return r.dataVersion
}

// BuildVersionMap returns a map of te version of each item.
func (r *Repo) BuildItemVersionMap() map[uuid.UUID]uint64 {
	versions := make(map[uuid.UUID]uint64)
	r.Lock()
	defer r.Unlock()
	for _, e := range r.entries {
		versions[e.Item.ID] = e.Item.Version
	}
	return versions
}

// StoreDataVersion changes dataVersion field in Repo object.
func (r *Repo) StoreDataVersion(dataVersion uint64) {
	r.dataVersion = dataVersion
	r.isChanged = true
}

func (r *Repo) StoreAccessToken(token models.AccessToken) {
	r.Lock()
	defer r.Unlock()
	r.accessToken = token
	r.isChanged = true
}

func (r *Repo) GetAccessToken() models.AccessToken {
	return r.accessToken
}

func (r *Repo) StoreRefreshToken(token models.RefreshToken) {
	r.Lock()
	defer r.Unlock()
	r.refreshToken = token
	r.isChanged = true
}

func (r *Repo) GetRefreshToken() models.RefreshToken {
	return r.refreshToken
}
