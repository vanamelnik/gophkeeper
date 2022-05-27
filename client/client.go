package client

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	maxRetries = 5

	retrySleepBase = float64(4)
)

var (
	ErrReloginNeeded = errors.New("relogin needed, session stopped")
)

type (
	// Client represents client interaction with the Gophkeeper server and the local storage.
	Client struct {
		ctx context.Context

		pbClient pb.GophkeeperClient

		// syncInterval - updates download interval.
		syncInterval time.Duration
		// sendInterval - interval for sending local updates to the serverю
		sendInterval time.Duration

		// repo is local items repository
		repo *repo.Repo

		eventCh chan models.Event
		closeCh chan struct{}

		maxNumberOfRetries int

		// conflictResolveFn - callback function for resolving merge conflicts by the user.
		conflictResolveFn ConflictResolveFn

		// eventsPool is the place where unsrnt events accumulate.
		eventsPool []models.Event
	}

	// ConflictResolveFn is callback function that invokes for merge conflict resolving.
	// If the user prefers the received item, the function returns true.
	ConflictResolveFn func(recievedItem models.Item, localEntry repo.Entry) (userChooseReceivedItem bool)
)

// New creates new client session. The user should be logged in beforehand.
func New(
	ctx context.Context,
	pbClient pb.GophkeeperClient,
	syncInterval,
	sendInterval time.Duration,
	storage *repo.Repo,
	accessToken models.AccessToken,
	refreshToken models.RefreshToken,
	conflictResolveFn ConflictResolveFn,
) (*Client, error) {

	c := Client{
		ctx:                ctx,
		pbClient:           pbClient,
		syncInterval:       syncInterval,
		sendInterval:       sendInterval,
		repo:               storage,
		eventCh:            make(chan models.Event, 1),
		closeCh:            make(chan struct{}),
		maxNumberOfRetries: maxRetries,
		conflictResolveFn:  conflictResolveFn,
		eventsPool:         make([]models.Event, 0),
	}

	// store auth token pair
	c.repo.StoreAccessToken(accessToken)
	c.repo.StoreRefreshToken(refreshToken)

	if err := c.WhatsNew(); err != nil {
		log.Println("Could not start the client - problems with connection (see messages above). Relogin needed.")
		return nil, ErrReloginNeeded
	}
	go c.worker()
	log.Println("client started")
	return &c, nil
}

func (c *Client) Close() {
	if c.closeCh != nil {
		close(c.closeCh)
		c.closeCh = nil
	}
}

// PublishEvent sends the event to the queue of events waiting to be sent to the server.
func (c *Client) PublishEvent(event models.Event) {
	c.eventCh <- event
}

// WhatsNew sends the WhatsNew request to the server.
// Multiple connection attempts are made. In case of failure, the user session ends.
// If the server presponses "update the data", GetUpdates function invoked.
func (c *Client) WhatsNew() error {
	request := &pb.WhatsNewRequest{
		Token:       &pb.AccessToken{AccessToken: string(c.repo.GetAccessToken())},
		DataVersion: c.repo.GetDataVersion(),
	}
	for i := 0; i < c.maxNumberOfRetries; i++ {
		_, err := c.pbClient.WhatsNew(c.ctx, request)
		if err == nil { // Server responses "all is up to date"
			return nil
		}

		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied { // Server responses "update the data"
			// get updates from the server
			dataVersion, items, err := c.GetUpdates()
			if err != nil {
				break
			}
			// merge updates
			c.MergeItems(dataVersion, items)
			return nil
		}

		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: WhatsNew: %s; trying to renew the token pair", err)
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: WhatsNew: %s; relogin needed", err)
					log.Println("client: WhatsNew: tokens are refreshed, trying again")
					break
				}
				continue // tokens are renewed - try again
			}
			log.Printf("client: WhatsNew: %s; relogin needed", err)
			break // unathenticated and the problem isn't in expired access token - relogin needed
		}
		if st.Code() == codes.Internal { // if there is internal error - try again
			timeToWait := time.Millisecond * time.Duration(math.Pow(retrySleepBase, float64(i)))
			log.Printf("client: WhatsNew: internal server error, try again in %v", timeToWait)
			time.Sleep(timeToWait)
			continue // internal server error - try again
		}
	}
	return ErrReloginNeeded
}

// GetUpdates fetches updates from the server.
// Multiple connection attempts are made. In case of failure, the error "relogin nedded" returns.
// Function panics if it could not marshall version map.
func (c *Client) GetUpdates() (uint64, []models.Item, error) {
	versionMap := c.repo.BuildItemVersionMap()
	versions, err := json.Marshal(versionMap)
	if err != nil {
		log.Fatalf("client: GetUpdates: could not marshall the map: %s", err)
	}
	request := &pb.DownloadUserDataRequest{
		Token:      &pb.AccessToken{AccessToken: string(c.repo.GetAccessToken())},
		VersionMap: string(versions),
	}
	for i := 0; i < c.maxNumberOfRetries; i++ {
		userData, err := c.pbClient.DownloadUserData(c.ctx, request)
		if err == nil {
			items := make([]models.Item, 0, len(userData.Items))
			for _, pbItem := range userData.Items {
				item, err := models.PbToItem(pbItem)
				if err != nil {
					log.Fatal(err)
				}
				items = append(items, item)
			}

			return userData.DataVersion, items, nil
		}

		st, _ := status.FromError(err)
		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: GetUpdates: %s; trying to renew the token pair", err)
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: GetUpdates: %s; relogin needed", err)
					break // could not renew the tokens - end the session!
				}
				log.Println("client: GetUpdates: tokens are refreshed, trying again")
				continue // tokens are renewed, try again
			}
			log.Printf("client: GetUpdates: %s; relogin needed", err)
			break // unauthenticated and problem isn't with expired access token - end the session!
		}
		if st.Code() == codes.Internal {
			timeToWait := time.Millisecond * time.Duration(math.Pow(retrySleepBase, float64(i)))
			log.Printf("client: GetUpdates: internal server error, try again in %v", timeToWait)
			time.Sleep(timeToWait)
			continue // internal server error - try again
		}
	}
	return 0, nil, ErrReloginNeeded
}

// sendEvents sends all accumulated events to the server.
func (c *Client) sendEvents() error {
	if len(c.eventsPool) == 0 {
		return nil
	}
	events := make([]*pb.Event, 0, len(c.eventsPool))
	for _, e := range c.eventsPool {
		events = append(events, &pb.Event{
			Operation: pb.Event_Operation(pb.Event_Operation_value[string(e.Operation)]),
			Item:      models.ItemToPb(e.Item),
		})
	}
	for i := 0; i < c.maxNumberOfRetries; i++ {
		_, err := c.pbClient.PublishLocalChanges(c.ctx, &pb.PublishLocalChangesRequest{
			Token:       &pb.AccessToken{AccessToken: string(c.repo.GetAccessToken())},
			DataVersion: c.repo.GetDataVersion(),
			Events:      events,
		})
		if err == nil {
			return nil
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: sendEvents: %s; trying to renew the token pair", err)
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: sendEvents: %s; relogin needed", err)
					break // could not renew the tokens - log out!
				}
				log.Println("client: sendEvents: tokens are refreshed, trying again")
				continue // tokens are renewed, try again
			}
			log.Printf("client: sendEvents: %s; relogin needed", err)
			break // unauthenticated and problem isn't with expired access token - log out!
		}
		if st.Code() == codes.Internal {
			timeToWait := time.Millisecond * time.Duration(math.Pow(retrySleepBase, float64(i)))
			log.Printf("client: sendEvents: internal server error, try again in %v", timeToWait)
			time.Sleep(timeToWait)
			continue // internal server error - try again
		}
	}
	return ErrReloginNeeded
}

// worker periodically fetches updates from the server and sends local updates to the server.
func (c *Client) worker() {
	whatsNew := time.NewTicker(c.syncInterval)
	defer whatsNew.Stop()
	timeToSend := time.NewTicker(c.sendInterval)
	defer timeToSend.Stop()

clientLoop:
	for {
		select {
		case <-c.closeCh:
			break clientLoop
		case <-whatsNew.C:
			if err := c.WhatsNew(); err != nil {
				log.Println("client: user relogin needed, session stopped")
				c.Close() // Relogin needed
				continue
			}
		case event := <-c.eventCh:
			c.eventsPool = append(c.eventsPool, event)
		case <-timeToSend.C:
			if err := c.sendEvents(); err != nil {
				log.Println("client: user relogin needed, session stopped")
				c.Close()
				continue
			}
			c.eventsPool = c.eventsPool[:0] // if all events are successfully sent, clear the events pool
		}
	}
	if c.eventCh != nil {
		close(c.eventCh)
		c.eventCh = nil
	}
	log.Println("client is stopped")
}
