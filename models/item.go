package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Item contain user data with the header:
	Item struct {
		ID        uuid.UUID
		CreatedAt time.Time
		DeletedAt time.Time

		// Data can be one of these types:
		//	- TextData
		//	- BinaryData
		//	- PasswordData
		//	- CardData
		Data interface{}

		Meta JSONMetadata
	}

	// Data types:

	// TextData contains a text string.
	TextData struct {
		Text string
	}

	// BinaryData contains binary data.
	BinaryData struct {
		Binary []byte
	}

	// PasswordData contains one of user's passwords.
	PasswordData struct {
		Password string
	}

	// CardData contains user's credit card data.
	CardData struct {
		Number         string
		CardholderName string
		Date           string
		CVC            uint32
	}

	// JSONMetadata is a JSON string that represents a set of key-value pairs,
	// that may contain different additional data such as login, bank name, kind of notes etc.
	// The types of such data are not strictly defined and must be handled on the client side.
	JSONMetadata string
)

// IsValidItem checks the type of item.Data field.
func IsValidItem(item Item) error {
	switch item.Data.(type) {
	case TextData, CardData, PasswordData, BinaryData:
		return nil
	default:
		return ErrInvalidItem
	}
}
