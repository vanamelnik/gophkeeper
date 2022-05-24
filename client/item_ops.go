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

// CreatePassword creates a new Password item, stores it in local repository
// and and queues an event to publish it to the server.
func (c Client) CreatePassword(p Password) error {
	password, err := PasswordToItem(p)
	if err != nil {
		return err
	}
	if err := c.repo.CreateItem(password); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      password,
	})
	return nil
}

// CreateCard creates a new Card item, stores it in local repository
// and and queues an event to publish it to the server.
func (c Client) CreateCard(cc CreditCard) error {
	card, err := CardToItem(cc)
	if err != nil {
		return err
	}
	if err := c.repo.CreateItem(card); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      card,
	})
	return nil
}

// CreateBlob creates a new Blob item, stores it in local repository
// and and queues an event to publish it to the server.
func (c Client) CreateBlob(b Blob) error {
	blob, err := BlobToItem(b)
	if err != nil {
		return err
	}
	if err := c.repo.CreateItem(blob); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      blob,
	})
	return nil
}

// CreateText creates a new Text item, stores it in local repository
// and queues an event to publish it to the server.
func (c Client) CreateText(t Text) error {
	text, err := TextToItem(t)
	if err != nil {
		return err
	}
	if err := c.repo.CreateItem(text); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpCreate,
		Item:      text,
	})
	return nil
}

// UpdatePassword updates a password in the local repository
// and queues an event to publish the changes.
func (c Client) UpdatePassword(p Password) error {
	password, err := PasswordToItem(p)
	if err != nil {
		return err
	}
	if err := c.repo.UpdateItem(password); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpUpdate,
		Item:      password,
	})
	return nil
}

// UpdateText updates a text item in the local repository
// and queues an event to publish the changes.
func (c Client) UpdateText(t Text) error {
	text, err := TextToItem(t)
	if err != nil {
		return err
	}
	if err := c.repo.UpdateItem(text); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpUpdate,
		Item:      text,
	})
	return nil
}

// UpdateBlob updates a binary item in the local repository
// and queues an event to publish the changes.
func (c Client) UpdateBlob(b Blob) error {
	blob, err := BlobToItem(b)
	if err != nil {
		return err
	}
	if err := c.repo.UpdateItem(blob); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpUpdate,
		Item:      blob,
	})
	return nil
}

// UpdateCard updates a card in the local repository
// and queues an event to publish the changes.
func (c Client) UpdateCard(card CreditCard) error {
	creditCard, err := CardToItem(card)
	if err != nil {
		return err
	}
	if err := c.repo.UpdateItem(creditCard); err != nil {
		return fmt.Errorf("could not store the item to local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpUpdate,
		Item:      creditCard,
	})
	return nil
}
