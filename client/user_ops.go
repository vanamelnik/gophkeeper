package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/proto"
	pb "github.com/vanamelnik/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const minPasswordLength = 6

// RenewTokens sends the refresh token to the server and renews client's accessToken and refreshToken.
func (c *Client) RenewTokens() error {
	userAuth, err := c.pbClient.GetNewTokens(c.ctx, &proto.RefreshToken{RefreshToken: string(c.repo.GetRefreshToken())})
	if err != nil {
		log.Printf("client: could not get new pair of tokens, relogin needed: %s", err)
		return ErrReloginNeeded
	}
	c.repo.StoreAccessToken(models.AccessToken(userAuth.AccessToken.String()))
	c.repo.StoreRefreshToken(models.RefreshToken(userAuth.RefreshToken.String()))

	return nil
}

// SignUp sends user's email and password to the server to register a new user.
// If the registration is successfull, the new user session is created and auth token pair is returned.
func SignUp(ctx context.Context, pbClient pb.GophkeeperClient, email, password string) (models.AccessToken, models.RefreshToken, error) {
	if err := validatePassword(password); err != nil {
		return "", "", err
	}
	userAuth, err := pbClient.SignUp(ctx, &pb.SignInData{
		Email:        email,
		UserPassword: password,
	})
	if err == nil {
		return models.AccessToken(userAuth.AccessToken.AccessToken), models.RefreshToken(userAuth.RefreshToken.RefreshToken), nil
	}
	se, _ := status.FromError(err)
	var errMsg string
	switch se.Code() {
	case codes.Internal:
		errMsg = fmt.Sprintf("signUp: internal server error: %s", se.Message())
	case codes.AlreadyExists:
		errMsg = fmt.Sprintf("signUp: user with email %s already exists: %s", email, se.Message())
	}
	log.Println(errMsg)

	return "", "", errors.New(errMsg)
}

// LogIn sends user's email and password to the server to authenticate the user and to create the new user session.
func LogIn(ctx context.Context, pbClient pb.GophkeeperClient, email, password string) (models.AccessToken, models.RefreshToken, error) {
	userAuth, err := pbClient.LogIn(ctx, &pb.SignInData{
		Email:        email,
		UserPassword: password,
	})
	if err == nil {
		return models.AccessToken(userAuth.AccessToken.AccessToken), models.RefreshToken(userAuth.RefreshToken.RefreshToken), nil
	}
	se, _ := status.FromError(err)
	var errMsg string
	switch se.Code() {
	case codes.Internal:
		errMsg = fmt.Sprintf("logIn: internal server error: %s", se.Message())
	case codes.NotFound:
		errMsg = fmt.Sprintf("logIn: user with email %s is not found: %s", email, se.Message())
	case codes.Unauthenticated:
		errMsg = fmt.Sprintf("logIn: could not authenticate the user with email %s: %s", email, se.Message())
	}
	log.Println(errMsg)

	return "", "", errors.New(errMsg)
}

func LogOut(ctx context.Context, pbClient pb.GophkeeperClient, r *repo.Repo) error {
	_, err := pbClient.LogOut(ctx, &pb.RefreshToken{RefreshToken: string(r.GetRefreshToken())})
	//regardless of the success of the operation, delete tokens from the repository
	r.StoreAccessToken("")
	r.StoreRefreshToken("")
	return err
}

func validatePassword(p string) error {
	// TODO: add more password complexity checks
	if len([]rune(p)) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}

	return nil
}
