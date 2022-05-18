package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vanamelnik/gophkeeper/models"
)

func PasswordToItem(p Password) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(p)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        p.ID,
		CreatedAt: &now,
		DeletedAt: nil,
		Data:      models.PasswordData{Password: p.Password},
		Meta:      models.JSONMetadata(meta),
	}, nil
}
func CardToItem(c CreditCard) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(c)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        c.ID,
		CreatedAt: &now,
		DeletedAt: nil,
		Data: models.CardData{
			Number:         c.Number,
			CardholderName: c.CardHolder,
			Date:           c.ExpirationDate,
			CVC:            c.CVC,
		},
		Meta: models.JSONMetadata(meta),
	}, nil
}
func BlobToItem(b Blob) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(b)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        b.ID,
		CreatedAt: &now,
		DeletedAt: nil,
		Data:      models.BinaryData{Binary: b.Data},
		Meta:      models.JSONMetadata(meta),
	}, nil

}
func TextToItem(t Text) (models.Item, error) {
	now := time.Now()
	meta, err := json.Marshal(t)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        t.ID,
		CreatedAt: &now,
		DeletedAt: nil,
		Data:      models.TextData{Text: t.Text},
		Meta:      models.JSONMetadata(meta),
	}, nil
}
