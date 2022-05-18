package client

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

// Local user data types
type (
	Password struct {
		ID       uuid.UUID `json:"-"`
		Login    string    `json:"login"`
		Password string    `json:"-"`
		Notes    string    `json:"notes,omitempty"`
	}

	Blob struct {
		ID    uuid.UUID `json:"-"`
		Data  []byte    `json:"-"`
		Notes string    `json:"notes,omitempty"`
	}

	Text struct {
		ID   uuid.UUID `json:"-"`
		Text string    `json:"-"`
		Note string    `json:"notes,omitempty"`
	}

	CreditCard struct {
		ID             uuid.UUID `json:"-"`
		BankName       string    `json:"bank_name"`
		Number         string    `json:"-"`
		ExpirationDate string    `json:"-"`
		CardHolder     string    `json:"-"`
		CVC            uint32    `json:"-"`
		Notes          string    `json:"notes,omitempty"`
	}
)

// CreatePassword creates a new Password item, stores it in local storage
// and and queues an event to publish it to the server.
func (c Client) CreatePassword(p Password) error {
	password, err := PasswordToItem(p)
	if err != nil {
		return err
	}
	if err := c.storage.CreateItem(password); err != nil {
		return fmt.Errorf("could not store the item to local storage: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      password,
	})
	return nil
}

// CreateCard creates a new Card item, stores it in local storage
// and and queues an event to publish it to the server.
func (c Client) CreateCard(cc CreditCard) error {
	card, err := CardToItem(cc)
	if err != nil {
		return err
	}
	if err := c.storage.CreateItem(card); err != nil {
		return fmt.Errorf("could not store the item to local storage: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      card,
	})
	return nil
}

// CreateBlob creates a new Blob item, stores it in local storage
// and and queues an event to publish it to the server.
func (c Client) CreateBlob(b Blob) error {
	blob, err := BlobToItem(b)
	if err != nil {
		return err
	}
	if err := c.storage.CreateItem(blob); err != nil {
		return fmt.Errorf("could not store the item to local storage: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      blob,
	})
	return nil
}

// CreateText creates a new Text item, stores it in local storage
// and and queues an event to publish it to the server.
func (c Client) CreateText(t Text) error {
	text, err := TextToItem(t)
	if err != nil {
		return err
	}
	if err := c.storage.CreateItem(text); err != nil {
		return fmt.Errorf("could not store the item to local storage: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      text,
	})
	return nil
}

func (c Client) UpdatePassword(p Password) error {
	password, err := PasswordToItem(p)
	if err != nil {
		return err
	}
	if err := c.storage.UpdateItem(password); err != nil {
		return fmt.Errorf("could not store the item to local storage: %w", err)

	}
}
