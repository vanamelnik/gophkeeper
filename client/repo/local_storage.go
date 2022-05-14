package repo

import (
	"sync"
	"time"

	"github.com/vanamelnik/gophkeeper/models"
)

// Package local contains tools for storing user data locally on the client side.

type (
	// Item wraps user data types described in models package with Pending flag.
	Item struct {
		// Data can be one of following types:
		//	- models.TextItem
		//	- models.BinaryItem
		//	- models.PasswordItem
		//	- models.CardItem
		Data interface{}
		// Pending is the flag indicating that the data was created locally
		// and is awaiting confirmation that it was successfully stored on the server.
		Pending bool
	}

	// Storage represents the local in-memory storage of user data.
	Storage struct {
		// DataVersion is the version of current user data snapshot.
		DataVersion uint64

		sync.Mutex
		// Items contains user data wrapped in Item struct (with Pending flag).
		Items []Item

		// fileName is the name of the file in .gob format that contains
		// current snapshot of local user data.
		fileName string
		// flushInterval is the interval after which the file is updated if isChange flag is set.
		flushInterval time.Duration

		// isChanged indicates if local user data has been changed and needs to be stored in the file.
		isChanged bool
	}
)

func (s *Storage) CreateItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()

	return nil
}

func (s *Storage) UpdateItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteItem(item models.Item) error {
	if err := models.IsValidItem(item); err != nil {
		return err
	}
	return nil
}
