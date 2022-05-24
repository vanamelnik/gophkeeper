package client

import (
	"log"

	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/proto"
)

// RenewTokens sends the refresh token to the server and renews client's accessToken and refreshToken.
func (c Client) RenewTokens() error {
	userAuth, err := c.pbClient.GetNewTokens(c.ctx, &proto.RefreshToken{RefreshToken: string(c.refreshToken)})
	if err != nil {
		log.Printf("client: could not get new pair of tokens, relogin needed: %s", err)
		return ErrReloginNeeded
	}
	c.accessToken = models.AccessToken(userAuth.AccessToken.String())
	c.refreshToken = models.RefreshToken(userAuth.RefreshToken.String())

	return nil
}
