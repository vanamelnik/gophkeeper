package client

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sessionWorker periodically fetches updates from the server and sends local updates to the server.
func (c *Client) sessionWorker() {
	whatsNew := time.NewTicker(c.syncInterval)
	defer whatsNew.Stop()
	timeToSend := time.NewTicker(c.sendInterval)
	defer timeToSend.Stop()

clientLoop:
	for {
		select {
		case <-c.stopCh:
			break clientLoop
		case <-whatsNew.C:
			if err := c.whatsNew(); err != nil {
				if errors.Is(err, ErrReloginNeeded) {
					log.Println("client: user relogin needed, session stopped")
					c.CloseClientSession() // Relogin needed
					//nolint: errcheck
					c.LogOut()
					continue
				}
			}
		case event := <-c.eventCh:
			c.eventsPool = append(c.eventsPool, event)
		case <-timeToSend.C:
			if err := c.sendEvents(); err != nil {
				log.Println("client: user relogin needed, session stopped")
				c.CloseClientSession()
				continue
			}
			c.eventsPool = c.eventsPool[:0] // if all events are successfully sent, clear the events pool
		}
	}
}

// whatsNew sends the whatsNew request to the server.
// Multiple connection attempts are made. In case of failure, the user session ends.
// If the server presponses "update the data", GetUpdates function invoked.
func (c *Client) whatsNew() (retErr error) {
	request := &pb.WhatsNewRequest{
		Token:       &pb.AccessToken{AccessToken: string(c.repo.GetAccessToken())},
		DataVersion: c.repo.GetDataVersion(),
	}
	firstTime := true
	for i := 0; i < c.maxNumberOfRetries; i++ {
		if !firstTime {
			// timeout grows exponentially from 0 to 625ms.
			timeToWait := time.Millisecond * time.Duration(math.Pow(retrySleepBase, float64(i)))
			log.Printf("client: WhatsNew: internal server error, try again in %v attempt %d", timeToWait, i+1)
			time.Sleep(timeToWait)
		}
		firstTime = false
		_, err := c.pbClient.WhatsNew(c.ctx, request)
		if err == nil { // Server responses "all is up to date"
			return nil
		}
		retErr = err // save the last thrown error
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied { // Server responses "update the data"
			// get updates from the server
			dataVersion, items, err := c.getUpdates()
			if err != nil {
				return err
			}
			// merge updates
			c.MergeItems(dataVersion, items)
			return nil
		}
		if st.Code() == codes.Unavailable {
			continue // Server unavailable - try again, you know it could work...
		}
		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: WhatsNew: %s; trying to renew the token pair", err)
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: WhatsNew: %s; relogin needed", err)
					return err
				}
				log.Println("client: WhatsNew: tokens are refreshed, trying again")
				continue // tokens are renewed - try again
			}
			log.Printf("client: WhatsNew: %s; relogin needed", err)

			return ErrReloginNeeded // unathenticated and the problem isn't in expired access token - relogin needed
		}
		if st.Code() == codes.Internal { // if there is internal error - try again
			continue // internal server error - try again
		}
	}
	return retErr
}

// getUpdates fetches updates from the server.
// Multiple connection attempts are made. In case of failure, the error "relogin nedded" returns.
// Function panics if it could not marshal version map.
func (c *Client) getUpdates() (uint64, []models.Item, error) {
	versionMap := c.repo.BuildItemVersionMap()
	versions, err := json.Marshal(versionMap)
	if err != nil {
		log.Fatalf("client: GetUpdates: could not marshal the map: %s", err)
	}
	request := &pb.DownloadUserDataRequest{
		Token:      &pb.AccessToken{AccessToken: string(c.repo.GetAccessToken())},
		VersionMap: string(versions),
	}
	var retErr error
	firstTime := true
	for i := 0; i < c.maxNumberOfRetries; i++ {
		if !firstTime {
			// timeout grows exponentially from 0 to 625ms.
			timeToWait := time.Millisecond * time.Duration(math.Pow(retrySleepBase, float64(i)))
			log.Printf("client: GetUpdates: internal server error, try again in %v attempt %d", timeToWait, i+1)
			time.Sleep(timeToWait)
		}
		firstTime = false
		userData, err := c.pbClient.DownloadUserData(c.ctx, request)
		retErr = err
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
		if st.Code() == codes.Unavailable {
			continue // Server unavailable - try again, you know it could work...
		}

		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: GetUpdates: %s; trying to renew the token pair", err)
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: GetUpdates: %s; relogin needed", err)
					return 0, nil, ErrReloginNeeded // could not renew the tokens - end the session!
				}
				log.Println("client: GetUpdates: tokens are refreshed, trying again")

				continue // tokens are renewed, try again
			}
			log.Printf("client: GetUpdates: %s; relogin needed", err)

			return 0, nil, ErrReloginNeeded // unauthenticated and problem isn't with expired access token - end the session!
		}
		if st.Code() == codes.Internal {
			continue // internal server error - try again
		}
	}

	return 0, nil, retErr
}

// sendEvents sends all accumulated events to the server.
func (c *Client) sendEvents() (retErr error) {
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
	firstTime := true
	for i := 0; i < c.maxNumberOfRetries; i++ {
		if !firstTime {
			// timeout grows exponentially from 0 to 625ms.
			timeToWait := time.Millisecond * time.Duration(math.Pow(retrySleepBase, float64(i)))
			log.Printf("client: sendEvents: internal server error, try again in %v attempt %d", timeToWait, i+1)
			time.Sleep(timeToWait)
		}
		firstTime = false

		log.Printf("sendEvents: %+v", events)
		_, err := c.pbClient.PublishLocalChanges(c.ctx, &pb.PublishLocalChangesRequest{
			Token:       &pb.AccessToken{AccessToken: string(c.repo.GetAccessToken())},
			DataVersion: c.repo.GetDataVersion(),
			Events:      events,
		})
		log.Println("done", err)
		if err == nil {
			return nil
		}
		retErr = err
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			continue
		}
		if st.Code() == codes.Unauthenticated {
			if strings.Contains(err.Error(), users.ErrAccessTokenExpired.Error()) {
				log.Printf("client: sendEvents: %s; trying to renew the token pair", err)
				if err := c.RenewTokens(); err != nil { // try to renew the tokens
					log.Printf("client: sendEvents: %s; relogin needed", err)
					return err // could not renew the tokens - log out!
				}
				log.Println("client: sendEvents: tokens are refreshed, trying again")
				continue // tokens are renewed, try again
			}
			log.Printf("client: sendEvents: %s; relogin needed", err)
			return err // unauthenticated and problem isn't with expired access token - log out!
		}
		if st.Code() == codes.Internal {
			continue // internal server error - try again
		}
	}
	return retErr
}
