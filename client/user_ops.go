package client

import (
	"errors"
	"fmt"
	"log"

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
func (c *Client) SignUp(email, password string) error {
	if c.IsLoggedIn() {
		err := c.LogOut()
		if err != nil {
			log.Println(err)
		}
	}
	// TODO: validate email
	if err := validatePassword(password); err != nil {
		return err
	}
	userAuth, err := c.pbClient.SignUp(c.ctx, &pb.SignInData{
		Email:        email,
		UserPassword: password,
	})
	if err == nil {
		c.repo.StoreAccessToken(models.AccessToken(userAuth.AccessToken.AccessToken))
		c.repo.StoreRefreshToken(models.RefreshToken(userAuth.RefreshToken.RefreshToken))
		return nil
	}
	se, _ := status.FromError(err)
	var errMsg string
	switch se.Code() {
	case codes.Unavailable:
		return ErrUnavailable
	case codes.Internal:
		errMsg = fmt.Sprintf("signUp: internal server error: %s", se.Message())
	case codes.AlreadyExists:
		errMsg = fmt.Sprintf("signUp: user with email %s already exists", email)
	default:
		return err
	}

	return errors.New(errMsg)
}

// LogIn sends user's email and password to the server to authenticate the user and to create the new user session.
// The token pair is stored in the local repository.
func (c *Client) LogIn(email, password string) error {
	if c.IsLoggedIn() {
		err := c.LogOut()
		if err != nil {
			log.Println(err)
		}
	}
	userAuth, err := c.pbClient.LogIn(c.ctx, &pb.SignInData{
		Email:        email,
		UserPassword: password,
	})
	if err == nil {
		c.repo.StoreAccessToken(models.AccessToken(userAuth.AccessToken.AccessToken))
		c.repo.StoreRefreshToken(models.RefreshToken(userAuth.RefreshToken.RefreshToken))
		return nil
	}
	se, _ := status.FromError(err)
	var errMsg string
	switch se.Code() {
	case codes.Unavailable:
		return ErrUnavailable
	case codes.Internal:
		errMsg = fmt.Sprintf("internal server error: %s", se.Message())
	case codes.NotFound, codes.Unauthenticated:
		errMsg = fmt.Sprintf("user with email %s is not found or the password is wrong", email)
	default:
		return err
	}

	return errors.New(errMsg)
}

func (c *Client) LogOut() error {
	_, err := c.pbClient.LogOut(c.ctx, &pb.RefreshToken{RefreshToken: string(c.repo.GetRefreshToken())})
	//regardless of the success of the operation, delete tokens from the repository
	c.repo.StoreAccessToken("")
	c.repo.StoreRefreshToken("")
	return err
}

func (c *Client) IsLoggedIn() bool {
	return c.repo.GetAccessToken() != "" && c.repo.GetRefreshToken() != ""
}

func validatePassword(p string) error {
	// TODO: add more password complexity checks
	if len([]rune(p)) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}

	return nil
}
