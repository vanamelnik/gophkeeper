package client

import (
	"encoding/json"
	"fmt"

	"github.com/vanamelnik/gophkeeper/models"
)

//
// TODO: implement encrypting of user data
//

// PasswordToItem converts local Password struct to the models.Item object.
func PasswordToItem(p Password) (models.Item, error) {
	meta, err := json.Marshal(p)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        p.ID,
		Version:   p.Version,
		CreatedAt: &p.CreatedAt,
		DeletedAt: nil,
		Payload:   models.PasswordData{Password: p.Password},
		Meta:      models.JSONMetadata(meta),
	}, nil
}

// CardToItem converts local Card struct to the models.Item object.
func CardToItem(c CreditCard) (models.Item, error) {
	meta, err := json.Marshal(c)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        c.ID,
		Version:   c.Version,
		CreatedAt: &c.CreatedAt,
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
	meta, err := json.Marshal(b)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        b.ID,
		Version:   b.Version,
		CreatedAt: &b.CreatedAt,
		DeletedAt: nil,
		Payload:   models.BinaryData{Binary: b.Data},
		Meta:      models.JSONMetadata(meta),
	}, nil

}

// TextToItem converts local Text struct to the models.Item object.
func TextToItem(t Text) (models.Item, error) {
	meta, err := json.Marshal(t)
	if err != nil {
		return models.Item{}, fmt.Errorf("could not encode metadata: %w", err)
	}
	return models.Item{
		ID:        t.ID,
		Version:   t.Version,
		CreatedAt: &t.CreatedAt,
		DeletedAt: nil,
		Payload:   models.TextData{Text: t.Text},
		Meta:      models.JSONMetadata(meta),
	}, nil
}

// ItemToPassword converts models.Item to local Password format.
func ItemToPassword(item models.Item) (Password, error) {
	pwd, ok := item.Payload.(models.PasswordData)
	if !ok {
		return Password{}, ErrWrongPayloadType
	}
	p := Password{
		ID:        item.ID,
		Version:   item.Version,
		CreatedAt: *item.CreatedAt,
		Password:  pwd.Password,
	}
	if err := json.Unmarshal([]byte(item.Meta), &p); err != nil {
		return Password{}, err
	}

	return p, nil
}

// ItemToText converts models.Item to local Text format.
func ItemToText(item models.Item) (Text, error) {
	txt, ok := item.Payload.(models.TextData)
	if !ok {
		return Text{}, ErrWrongPayloadType
	}
	t := Text{
		ID:        item.ID,
		Version:   item.Version,
		CreatedAt: *item.CreatedAt,
		Text:      txt.Text,
	}
	if err := json.Unmarshal([]byte(item.Meta), &t); err != nil {
		return Text{}, err
	}

	return t, nil
}

// ItemToBlob converts models.Item to local Blob format.
func ItemToBlob(item models.Item) (Blob, error) {
	bin, ok := item.Payload.(models.BinaryData)
	if !ok {
		return Blob{}, ErrWrongPayloadType
	}
	b := Blob{
		ID:        item.ID,
		Version:   item.Version,
		CreatedAt: *item.CreatedAt,
		Data:      bin.Binary,
	}
	if err := json.Unmarshal([]byte(item.Meta), &b); err != nil {
		return Blob{}, err
	}

	return b, nil
}

// ItemToBlob converts models.Item to local Blob format.
func ItemToCreditCard(item models.Item) (CreditCard, error) {
	crd, ok := item.Payload.(models.CardData)
	if !ok {
		return CreditCard{}, ErrWrongPayloadType
	}
	card := CreditCard{
		ID:             item.ID,
		Version:        item.Version,
		CreatedAt:      *item.CreatedAt,
		Number:         crd.Number,
		ExpirationDate: crd.Date,
		CardHolder:     crd.CardholderName,
		CVC:            crd.CVC,
	}
	if err := json.Unmarshal([]byte(item.Meta), &card); err != nil {
		return CreditCard{}, err
	}

	return card, nil
}
