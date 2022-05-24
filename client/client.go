package client

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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
)

var (
	ErrReloginNeeded = errors.New("relogin needed, session stopped")
)

type (
	Client struct {
		ctx context.Context

		pbClient pb.GophkeeperClient

		syncInterval time.Duration
		sendInterval time.Duration

		repo *repo.Repo

		eventCh chan models.Event
		closeCh chan struct{}

		accessToken  models.AccessToken
		refreshToken models.RefreshToken

		maxNumberOfRetries int

		conflictResolveFn ConflictResolveFn
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
		accessToken:        accessToken,
		refreshToken:       refreshToken,
		maxNumberOfRetries: maxRetries,
		conflictResolveFn:  conflictResolveFn,
	}
	if err := c.WhatsNew(); err != nil {
		log.Println("Could not start the client - problems with connection (see messages above). Relogin needed.")
		return nil, ErrReloginNeeded
	}
	go c.worker()
	return &c, nil
}

func (c Client) Close() {
	if c.closeCh != nil {
		close(c.closeCh)
		c.closeCh = nil
	}
}

// PublishEvent sends the event to the queue of events waiting to be sent to the server.
func (c Client) PublishEvent(event models.Event) {
	c.eventCh <- event
}

// WhatsNew sends the WhatsNew request to the server.
// Multiple connection attempts are made. In case of failure, the user session ends.
func (c Client) WhatsNew() error {
	request := &pb.WhatsNewRequest{
		Token:       &pb.AccessToken{AccessToken: string(c.accessToken)},
		DataVersion: c.repo.GetDataVersion(),
	}
	for tryCount := 0; tryCount < c.maxNumberOfRetries; tryCount++ {
		_, err := c.pbClient.WhatsNew(c.ctx, request)
		if err == nil { // Server responses "all is up to date"
			return nil
		}

		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound { // Server responses "update the data"
			dataVersion, items, err := c.GetUpdates()
			if err != nil {
				break
			}
			c.Merge(dataVersion, items)
			return nil
		}

		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: WhatsNew: %s; trying to renew the token pair")
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: WhatsNew: %s; relogin needed", err)

					break
				}
				continue // tokens are renewed - try again
			}
			log.Printf("client: WhatsNew: %s; relogin needed", err)
			break // unathenticated and the problem isn't in expired access token - relogin needed
		}
		if st.Code() == codes.Internal { // if there is internal error - try again
			log.Printf("client: WhatsNew: could not get WhatsNew information: %s", err)
			continue
		}
	}
	return ErrReloginNeeded
}

// GetUpdates fetches updates from the server.
// Multiple connection attempts are made. In case of failure, the error "relogin nedded" returns.
// Function panics if it could not marshall version map.
func (c Client) GetUpdates() (uint64, []models.Item, error) {
	versionMap := c.repo.BuildItemVersionMap()
	versions, err := json.Marshal(versionMap)
	if err != nil {
		log.Fatalf("client: GetUpdates: could not marshall the map: %s", err)
	}
	request := &pb.DownloadUserDataRequest{
		Token:      &pb.AccessToken{AccessToken: string(c.accessToken)},
		VersionMap: string(versions),
	}
	for tryCount := 0; tryCount < c.maxNumberOfRetries; tryCount++ {
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
				log.Printf("client: GetUpdates: %s; trying to renew the token pair")
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: GetUpdates: %s; relogin needed", err)
					break // could not renew the tokens - end the session!
				}
				continue // tokens are renewed, try again
			}
			log.Printf("client: GetUpdates: %s; relogin needed", err)
			break // unauthenticated and problem isn't with expired access token - end the session!
		}
		if st.Code() == codes.Internal {
			continue // internal server error - try again
		}
	}
	return 0, nil, ErrReloginNeeded
}

func (c Client) sendEvents(event models.Event) {
	// ...
}

// worker periodically fetches updates from the server and sends local updates to the server.
func (c Client) worker() {
	whatsNew := time.NewTicker(c.syncInterval)
	defer whatsNew.Stop()
	timeToSend := time.NewTicker(c.sendInterval)
	defer timeToSend.Stop()

	for {
		select {
		case <-c.closeCh:
			break
		case <-whatsNew.C:
			if err := c.WhatsNew(); err != nil {
				log.Println("client: user relogin needed, session stopped")
				c.Close() // Relogin needed
			}
		case event := <-c.eventCh:
			c.sendEvent(event)
		case <-timeToSend.C:
		}
	}
}
