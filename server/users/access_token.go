package users

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const jwtIssuer = "GophKeeper"

func (s Service) NewAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(s.accessTokenDuration).Unix(),
		Id:        userID,
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		NotBefore: time.Now().Unix(),
	})
	ss, err := token.SignedString([]byte(s.accessTokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (s Service) Authenticate(accessToken string) (uuid.UUID, error) {
	t, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return s.accessTokenSecret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	// TODO: write token validation code
	// ...
	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		// TODO: return incorrect accessToken error
	}
	id, err := uuid.Parse(claims.Id)
	if err != nil {
		// TODO: return incorrect userID error
	}
	return id, nil
}
