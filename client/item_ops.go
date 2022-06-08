package client

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

// Local user data types
type (
	Password struct {
		ID        uuid.UUID `json:"-"`
		Version   uint64    `json:"-"`
		CreatedAt time.Time `json:"-"`
		Login     string    `json:"login"`
		Password  string    `json:"-"`
		Notes     string    `json:"notes,omitempty"`
	}

	Blob struct {
		ID        uuid.UUID `json:"-"`
		Version   uint64    `json:"-"`
		CreatedAt time.Time `json:"-"`
		Data      []byte    `json:"-"`
		Notes     string    `json:"notes,omitempty"`
	}

	Text struct {
		ID        uuid.UUID `json:"-"`
		Version   uint64    `json:"-"`
		CreatedAt time.Time `json:"-"`
		Text      string    `json:"-"`
		Note      string    `json:"notes,omitempty"`
	}

	CreditCard struct {
		ID             uuid.UUID `json:"-"`
		Version        uint64    `json:"-"`
		CreatedAt      time.Time `json:"-"`
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
func (c *Client) CreatePassword(p Password) error {
	if !c.activeSession {
		return ErrSessionInactive
	}
	password, err := PasswordToItem(p)
	if err != nil {
		return err
	}
	// log.Printf("CreatePassword: %+v", password) // FIXME: remove it!
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
func (c *Client) CreateCard(cc CreditCard) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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
func (c *Client) CreateBlob(b Blob) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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
func (c *Client) CreateText(t Text) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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
func (c *Client) UpdatePassword(p Password) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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
func (c *Client) UpdateText(t Text) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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
func (c *Client) UpdateBlob(b Blob) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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
func (c *Client) UpdateCard(card CreditCard) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

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

// DeletePassword deletes the password from the local repository
// and queues an event to publish the changes.
func (c *Client) DeletePassword(p Password) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

	return c.deleteItem(p.ID)
}

// DeleteCard deletes the card from the local repository
// and queues an event to publish the changes.
func (c *Client) DeleteCard(card CreditCard) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

	return c.deleteItem(card.ID)
}

// DeleteText deletes the text from the local repository
// and queues an event to publish the changes.
func (c *Client) DeleteText(t Text) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

	return c.deleteItem(t.ID)
}

// DeleteBlob deletes the blob from the local repository
// and queues an event to publish the changes.
func (c *Client) DeleteBlob(b Blob) error {
	if !c.activeSession {
		return ErrSessionInactive
	}

	return c.deleteItem(b.ID)
}

// GetPasswords retrieves all password items from the local storage and converts them
// to local Password format.
func (c *Client) GetPasswords() ([]Password, error) {
	if !c.activeSession {
		return nil, ErrSessionInactive
	}

	passwords := make([]Password, 0, 0)

	for _, entry := range c.repo.GetDataSnapshot() {
		if _, ok := entry.Item.Payload.(models.PasswordData); ok {
			password, err := ItemToPassword(entry.Item)
			if err != nil {
				panic(fmt.Sprintf("unreachable: could not convert item %+v payload type %T to Password: %s", entry.Item, entry.Item.Payload, err))
			}
			passwords = append(passwords, password)
		}
	}

	return passwords, nil
}

// GetBlobs retrieves all blob items from the local storage and converts them
// to local Blob format.
func (c *Client) GetBlobs() ([]Blob, error) {
	if !c.activeSession {
		return nil, ErrSessionInactive
	}

	blobs := make([]Blob, 0, 0)

	for _, entry := range c.repo.GetDataSnapshot() {
		if _, ok := entry.Item.Payload.(models.BinaryData); ok {
			blob, err := ItemToBlob(entry.Item)
			if err != nil {
				panic(fmt.Sprintf("unreachable: could not convert item %+v payload type %T to Blob: %s", entry.Item, entry.Item.Payload, err))
			}
			blobs = append(blobs, blob)
		}
	}

	return blobs, nil
}

// GetCards retrieves all card items from the local storage and converts them
// to local CreditCard format.
func (c *Client) GetCards() ([]CreditCard, error) {
	if !c.activeSession {
		return nil, ErrSessionInactive
	}

	cards := make([]CreditCard, 0, 0)

	for _, entry := range c.repo.GetDataSnapshot() {
		if _, ok := entry.Item.Payload.(models.PasswordData); ok {
			card, err := ItemToCreditCard(entry.Item)
			if err != nil {
				panic(fmt.Sprintf("unreachable: could not convert item %+v payload type %T to CreditCard: %s", entry.Item, entry.Item.Payload, err))
			}
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// GetTexts retrieves all text items from the local storage and converts them
// to local TExt format.
func (c *Client) GetTexts() ([]Text, error) {
	texts := make([]Text, 0, 0)

	for _, entry := range c.repo.GetDataSnapshot() {
		if _, ok := entry.Item.Payload.(models.TextData); ok {
			text, err := ItemToText(entry.Item)
			if err != nil {
				panic(fmt.Sprintf("unreachable: could not convert item %+v payload type %T to Text: %s", entry.Item, entry.Item.Payload, err))
			}
			texts = append(texts, text)
		}
	}

	return texts, nil
}

// deleteItem invokes local repository method DeleteItem and creates an event to publish the changes.
func (c *Client) deleteItem(id uuid.UUID) error {
	deleted, err := c.repo.DeleteItem(id)
	if err != nil {
		return fmt.Errorf("could not perform changes in the local repository: %w", err)
	}
	c.PublishEvent(models.Event{
		Operation: models.OpUpdate,
		Item:      deleted,
	})
	return nil
}
