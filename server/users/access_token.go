package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

const jwtIssuer = "GophKeeper"

// newAccessToken creates a new token for the user provided and signs it with the secret key.
func (s Service) newAccessToken(userID uuid.UUID) (models.AccessToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(s.accessTokenDuration).Unix(),
		Id:        userID.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		NotBefore: time.Now().Unix(),
	})
	ss, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("users: NewAccessToken: %w", err)
	}
	return models.AccessToken(ss), nil
}

// Authenticate checks if the given access token is valid and, if so, returns the user ID.
func (s Service) Authenticate(ctx context.Context, accessToken models.AccessToken) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(string(accessToken), &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		var ve jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return uuid.Nil, fmt.Errorf("users: authenticate: %w", ErrAccessTokenExpired)
		}
		return uuid.Nil, fmt.Errorf("users: authenticate: %w", err)
	}
	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("users: authenticate: %w", ErrIncorrectAccessToken)
	}
	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("users: authenticate: %w: userID=%s", ErrIncorrectUserID, claims.Id)
	}
	return id, nil
}
