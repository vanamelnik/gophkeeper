package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vanamelnik/gophkeeper/models"
)

// PasswordToItem converts local Password struct to the models.Item object.
func PasswordToItem(p Password) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(p)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        p.ID,
		Version:   0,
		CreatedAt: &now,
		DeletedAt: nil,
		Payload:   models.PasswordData{Password: p.Password},
		Meta:      models.JSONMetadata(meta),
	}, nil
}

// CardToItem converts local Card struct to the models.Item object.
func CardToItem(c CreditCard) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(c)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        c.ID,
		Version:   0,
		CreatedAt: &now,
		DeletedAt: nil,
		Payload: models.CardData{
			Number:         c.Number,
			CardholderName: c.CardHolder,
			Date:           c.ExpirationDate,
			CVC:            c.CVC,
		},
		Meta: models.JSONMetadata(meta),
	}, nil
}

// BlobToItem converts local Blob struct to the models.Item object.
func BlobToItem(b Blob) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(b)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        b.ID,
		Version:   0,
		CreatedAt: &now,
		DeletedAt: nil,
		Payload:   models.BinaryData{Binary: b.Data},
		Meta:      models.JSONMetadata(meta),
	}, nil

}

// TextToItem converts local Text struct to the models.Item object.
func TextToItem(t Text) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(t)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        t.ID,
		Version:   0,
		CreatedAt: &now,
		DeletedAt: nil,
		Payload:   models.TextData{Text: t.Text},
		Meta:      models.JSONMetadata(meta),
	}, nil
}
