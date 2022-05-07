package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID           uint
	Login        string
	PasswordHash string
}

type Session struct {
	ID           uint
	UserID       uint
	RefreshToken string
	AccessToken  jwt.Token
	LoginAt      time.Time
	LogoutAt     time.Time
	DataVersion  uint64
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
		CVC            int
	}

	MetaData map[string]string
)
