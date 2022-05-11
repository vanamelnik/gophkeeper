package models

// Package models defines canonical models for representation
// such objects as User, Session, Event and all kind of Items.

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user of the service.
// One User can have multiple sessions from different devices.
type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	DeletedAt    time.Time
}

// Session represents a single client session of the user with given ID.
// For each session a pair of AccessToken + RefreshToken is created.
// If AccessToken expires, RefreshToken will be the key to create a new token pair.
type Session struct {
	ID           uint
	UserID       uuid.UUID
	RefreshToken RefreshToken
	LoginAt      time.Time
	LogoutAt     time.Time
}

// RefreshToken is a random string that used for check if the service may generate
// a new access token to replace an expired one.
type RefreshToken string

// Items are the data types that are contained in the user's storage:
type (
	// TextItem contains a text string.
	TextItem struct {
		ItemHeader
		Text string
	}

	// BinaryItem contains binary data.
	BinaryItem struct {
		ItemHeader
		Binary []byte
	}

	// PasswordItem contains one of user's passwords.
	PasswordItem struct {
		ItemHeader
		Password string
	}

	// CardItem contains user's credit card data.
	CardItem struct {
		ItemHeader
		Number         string
		CardholderName string
		Date           string
		CVC            uint32
	}

	// ItemHeader contains the necessary service fields and metadata.
	ItemHeader struct {
		ItemID    uuid.UUID
		Meta      JSONMetadata
		CreatedAt time.Time
		DeletedAt time.Time
	}

	// JSONMetadata is a JSON string that represents a set of key-value pairs,
	// that may contain different additional data such as login, bank name, kind of notes etc.
	// The types of such data are not strictly defined and must be handled on the client side.
	JSONMetadata string
)

// UserData represents the snapshot of all user's data stored in the storage.
// Each snapshot has an unique version number that is incremented when the data
// is updated on the server.
type UserData struct {
	Version   uint64
	Texts     []TextItem
	Blobs     []BinaryItem
	Passwords []PasswordItem
	Cards     []CardItem
}
