package users

// Package users contains users service object.
import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vanamelnik/gophkeeper/models"
)

// Service contains contains methods that provide operations with registered users and user sessions.
type Service struct {
	// accessTokenSecret is the secret key for Access Token signing.
	accessTokenSecret    string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewService(accessTokenSecret string, accessTokenDuration, refreshTokenDuration time.Duration) Service {
	return Service{
		accessTokenSecret:    accessTokenSecret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// RefreshTokens checks whether the given refresh token is not expired. If the token is valid,
// a new pair of tokens is generated.
func (s Service) RefreshTokens(refreshToken models.RefreshToken) (jwt.Token, models.RefreshToken, error) {
	return jwt.Token{}, "", nil
}
