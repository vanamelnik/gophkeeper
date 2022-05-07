package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int
	Login        string
	PasswordHash string
}

type Session struct {
	ID           int
	UserID       int
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

	Item struct {
		ItemID    uuid.UUID
		Meta      MetaData
		CreatedAt time.Time
		DeletedAt time.Time
	}

	TextItem struct {
		Item
		Text string
	}

	BinaryItem struct {
		Item
		Binary []byte
	}

	PasswordItem struct {
		Item
		Password string
	}

	CardItem struct {
		Item
		Number         string
		CardholderName string
		Date           string
		CVC            int
	}

	MetaData map[string]string
)
