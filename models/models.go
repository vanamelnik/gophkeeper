package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uint
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	DeletedAt    time.Time
}

type RefreshToken string

type Session struct {
	ID           uuid.UUID
	UserID       uint
	RefreshToken string
	LoginAt      time.Time
	LogoutAt     time.Time
}

type (
	UserData struct {
		Version   uint64
		Texts     []TextItem
		Blobs     []BinaryItem
		Passwords []PasswordItem
		Cards     []CardItem
	}

	ItemHeader struct {
		ItemID    uuid.UUID
		Meta      MetaData
		CreatedAt time.Time
		DeletedAt time.Time
	}

	TextItem struct {
		ItemHeader
		Text string
	}

	BinaryItem struct {
		ItemHeader
		Binary []byte
	}

	PasswordItem struct {
		ItemHeader
		Password string
	}

	CardItem struct {
		ItemHeader
		Number         string
		CardholderName string
		Date           string
		CVC            uint32
	}

	MetaData map[string]string
)
