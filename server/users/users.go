package users

// Package users contains users service object.
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/pkg/bcrypt"
	"github.com/vanamelnik/gophkeeper/server/storage"
)

var (
	ErrAccessTokenExpired   = errors.New("access token expired")
	ErrIncorrectAccessToken = errors.New("incorrect access token")

	ErrRefreshTokenExpired   = errors.New("refresh token expired")
	ErrIncorrectRefreshToken = errors.New("incorrect refresh token")

	ErrIncorrectUserID = errors.New("incorrect user ID")
)

// Service contains contains methods that provide operations with registered users and user sessions.
type Service struct {
	storage storage.Storage
	// secret is the secret key for Access Token and Refresh Token signing.
	secret               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewService(storage storage.Storage, secret string, accessTokenDuration, refreshTokenDuration time.Duration) Service {
	return Service{
		storage:              storage,
		secret:               secret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// CreateUser stores user info in the database and returns user ID.
func (s Service) CreateUser(ctx context.Context, userName, pwHash string) (uuid.UUID, error) {
	user, err := s.storage.CreateUser(ctx, userName, pwHash)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

// CreateSession creates a new user session and returns the access and refresh tokens.
func (s Service) CreateSession(ctx context.Context, userID uuid.UUID) (models.AccessToken, models.RefreshToken, error) {
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return "", "", err
	}
	t, err := s.newRefreshToken(sessionID) // generate the first refreshToken
	if err != nil {
		return "", "", err
	}
	now := time.Now()
	if err = s.storage.CreateSession(ctx, models.Session{
		ID:           sessionID,
		UserID:       userID,
		RefreshToken: t,
		LoginAt:      &now,
		LogoutAt:     nil,
	}); err != nil {
		return "", "", err
	}
	accessToken, refreshToken, err := s.RefreshTheTokens(ctx, t)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshTheTokens checks whether the given refresh token is not expired. If the token is valid,
// a new pair of tokens is generated and the new refresh token is stored in the db.
func (s Service) RefreshTheTokens(ctx context.Context, refreshToken models.RefreshToken) (models.AccessToken, models.RefreshToken, error) {
	t, err := jwt.Parse(string(refreshToken), func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		var ve jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return "", "", fmt.Errorf("users: refreshTheTokens: %w", ErrRefreshTokenExpired)
		}
		return "", "", fmt.Errorf("users: refreshTheTokens: %w", err)
	}
	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", "", fmt.Errorf("users: refreshTheTokens: %w", ErrIncorrectRefreshToken)
	}
	sessionID, err := uuid.Parse(claims.Id)
	if err != nil {
		return "", "", fmt.Errorf("users: refreshTheTokens: %w: userID=%s", ErrIncorrectUserID, claims.Id)
	}

	session, err := s.storage.GetSessionByID(ctx, sessionID)
	if err != nil {
		return "", "", fmt.Errorf("could not find session with id %s: %w", sessionID, err)
	}
	if session.RefreshToken != refreshToken {
		return "", "", errors.New("internal error: refresh token is not correct")
	}

	newAccessToken, err := s.newAccessToken(session.UserID)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := s.newRefreshToken(sessionID)
	if err != nil {
		return "", "", err
	}

	// Store the new refresh token
	session.RefreshToken = newRefreshToken
	if err := s.storage.UpdateSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("could not update refresh token in the database: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// Login checks whether the user with the username provided exists and
// the given credentials are valid. If all is OK a new session is created.
func (s Service) Login(ctx context.Context, userName, password string) (models.AccessToken, models.RefreshToken, error) {
	user, err := s.storage.GetUserByLogin(ctx, userName)
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword(password, user.PasswordHash); err != nil {
		return "", "", err
	}
	return s.CreateSession(ctx, user.ID)
}

// Logout marks the user session in the db as logged out.
func (s Service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.storage.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	now := time.Now()
	session.LogoutAt = &now
	if err := s.storage.UpdateSession(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s Service) GetDataVersion(ctx context.Context, userID uuid.UUID) (uint64, error) {
	return s.storage.GetUserDataVersion(ctx, userID)
}
