package users

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Service struct {
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	UserID string
}

func (s Service) UserAuth(accessToken string) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret key"), nil // TODO: add keyfunc
	})
	if err != nil {
		return uuid.Nil, err
	}
	// TODO: write token validation code
	// ...
	claims, ok := t.Claims.(*AccessTokenClaims)
	if !ok {
		// TODO: return incorrect accessToken error
	}
	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		// TODO: return incorrect userID error
	}
	return id, nil
}
