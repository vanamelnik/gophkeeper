package users

import (
	"time"
)

type Service struct {
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
