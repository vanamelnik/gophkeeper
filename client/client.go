package client

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
)

const (
	maxRetries = 5

	retrySleepBase = float64(5)
)

// Errors:
var (
	ErrReloginNeeded    = errors.New("relogin needed, session stopped")
	ErrWrongPayloadType = errors.New("wrong payload type")
	ErrUnavailable      = errors.New("the server is unavailable")
	ErrSessionInactive  = errors.New("client session is inactive")
)

type (
	// Client represents client interaction with the Gophkeeper server and the local storage.
	Client struct {
		ctx context.Context

		activeSession bool

		pbClient pb.GophkeeperClient

		// syncInterval - updates download interval.
		syncInterval time.Duration
		// sendInterval - interval for sending local updates to the serverю
		sendInterval time.Duration

		// repo is local items repository
		repo *repo.Repo

		eventCh chan models.Event
		stopCh  chan struct{}

		maxNumberOfRetries int

		// conflictResolveFn - callback function for resolving merge conflicts by the user.
		conflictResolveFn ConflictResolveFn

		// eventsPool is the place where unsrnt events accumulate.
		eventsPool []models.Event
	}

	// ConflictResolveFn is callback function that invokes for merge conflict resolving.
	// If the user prefers the received item, the function returns true.
	ConflictResolveFn func(recievedItem models.Item, localEntry repo.Entry) (userChoseReceivedItem bool)
)

// New creates a new client. For starting a new client session use (*Client).NewSession().
func New(
	ctx context.Context,
	pbClient pb.GophkeeperClient,
	syncInterval,
	sendInterval time.Duration,
	storage *repo.Repo,
	conflictResolveFn ConflictResolveFn,
) *Client {

	c := Client{
		ctx:                ctx,
		activeSession:      false,
		pbClient:           pbClient,
		syncInterval:       syncInterval,
		sendInterval:       sendInterval,
		repo:               storage,
		eventCh:            make(chan models.Event, 1),
		stopCh:             nil, // stop chan will be created within the client session
		maxNumberOfRetries: maxRetries,
		conflictResolveFn:  conflictResolveFn,
		eventsPool:         make([]models.Event, 0),
	}

	return &c
}

// NewSession starts a new client session.
// User must be logged in and the token pair is stored in the local repository.
func (c *Client) NewSession() error {
	if !c.IsLoggedIn() {
		return ErrReloginNeeded
	}
	c.activeSession = true
	c.stopCh = make(chan struct{})
	if err := c.whatsNew(); err != nil {
		log.Printf("client: NewSession: whatsNew: %s", err)
	}
	go c.sessionWorker()
	log.Println("client started")
	return nil
}

func (c *Client) CloseClientSession() {
	if c.stopCh != nil {
		close(c.stopCh)
		c.stopCh = nil
	}
	c.activeSession = false
}

// PublishEvent sends the event to the queue of events waiting to be sent to the server.
func (c *Client) PublishEvent(event models.Event) {
	if c.activeSession {
		c.eventCh <- event
		return
	}

}
