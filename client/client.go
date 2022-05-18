package client

import (
	"time"

	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
)

type Client struct {
	pbClient     pb.GophkeeperClient
	syncInterval time.Duration
	storage      *repo.Storage
	eventCh      chan models.Event
	closeCh      chan struct{}
}

func New(pbClient pb.GophkeeperClient, syncInterval time.Duration, storage *repo.Storage) Client {
	c := Client{
		pbClient:     pbClient,
		syncInterval: syncInterval,
		storage:      storage,
		eventCh:      make(chan models.Event, 1),
		closeCh:      make(chan struct{}),
	}
	c.WhatsNew()
	go c.worker()
	return c
}

func (c Client) Close() {
	if c.closeCh != nil {
		close(c.closeCh)
		c.closeCh = nil
	}
}

func (c Client) PublishEvent(event models.Event) {
	c.eventCh <- event
}

func (c Client) WhatsNew() {
	// Get updates from the server
	// ...
}

func (c Client) sendEvent(event models.Event) {
	// Add pending item to the local repo

	// Try to send the item to the server

	// If repo is out of date, refresh it

	// 		If there is a conflict, resolve it

	// 		Send the item again
}

// worker periodically fetches updates from the server and sends local updates to the server.
func (c Client) worker() {
	whatsNew := time.NewTicker(c.syncInterval)
	for {
		select {
		case <-c.closeCh:
			break
		case <-whatsNew.C:
			c.WhatsNew()
		case event := <-c.eventCh:
			c.sendEvent(event)
		}
	}
}
