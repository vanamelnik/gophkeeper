package users

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

func (s Service) GetSessionID(rt models.RefreshToken) (uuid.UUID, error) {
	t, _ := jwt.ParseWithClaims(string(rt), &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return uuid.Nil, ErrIncorrectRefreshToken
	}
	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("users: getSessionID: %w", err)
	}
	return id, nil
}

// newRefreshToken creates a new signed refresh token.
func (s Service) newRefreshToken(sessionID uuid.UUID) (models.RefreshToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(s.refreshTokenDuration).Unix(),
		Id:        sessionID.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		NotBefore: time.Now().Unix(),
	})
	ss, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("users: NewRefreshToken: %w", err)
	}
	return models.RefreshToken(ss), nil
}
