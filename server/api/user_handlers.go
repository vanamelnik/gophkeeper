package api

import (
	"context"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/pkg/bcrypt"
	"github.com/vanamelnik/gophkeeper/server/storage"
	"github.com/vanamelnik/gophkeeper/server/users"

	pb "github.com/vanamelnik/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SignUp implements GophkeeperServer interface.
func (s server) SignUp(ctx context.Context, data *pb.SignInData) (*pb.UserAuth, error) {
	pwHash, err := bcrypt.BcryptPassword(data.UserPassword)
	if err != nil {
		log.Printf("api: signUp: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	userID, err := s.users.CreateUser(ctx, data.Email, pwHash)
	if err != nil {
		log.Printf("api: signUp: %s", err)
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	accessToken, refreshToken, err := s.users.CreateSession(ctx, userID)
	if err != nil {
		log.Printf("api: signUp: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("api: signUp: successfully created the new user: %s", data.Email)
	return &pb.UserAuth{
		AccessToken:  &pb.AccessToken{AccessToken: string(accessToken)},
		RefreshToken: &pb.RefreshToken{RefreshToken: string(refreshToken)},
	}, nil
}

// LogIn implements GophkeeperServer interface.
func (s server) LogIn(ctx context.Context, data *pb.SignInData) (*pb.UserAuth, error) {
	accessToken, refreshToken, err := s.users.Login(ctx, data.Email, data.UserPassword)
	if err != nil {
		log.Printf("api: logIn: %s", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Error(codes.Unauthenticated, "wrong password")
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("api: logIn: user %s logged in", data.Email)

	return &pb.UserAuth{
		AccessToken:  &pb.AccessToken{AccessToken: string(accessToken)},
		RefreshToken: &pb.RefreshToken{RefreshToken: string(refreshToken)},
	}, nil
}

// LogOut implements GophkeeperServer interface.
func (s server) LogOut(ctx context.Context, rt *pb.RefreshToken) (*empty.Empty, error) {
	sessionID, err := s.users.GetSessionID(models.RefreshToken(rt.RefreshToken))
	if err != nil {
		log.Printf("api: logOut: %s", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.users.Logout(ctx, sessionID); err != nil {
		log.Printf("api: logOut: %s", err)
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

// GetNewTokens implements GophkeeperServer interface.
func (s server) GetNewTokens(ctx context.Context, rt *pb.RefreshToken) (*pb.UserAuth, error) {
	accessToken, refreshToken, err := s.users.RefreshTheTokens(ctx, models.RefreshToken(rt.RefreshToken))
	if err != nil {
		log.Printf("api: getNewTokens: %s", err)
		// logout if refresh token is expired
		if errors.Is(err, users.ErrRefreshTokenExpired) {
			sessionID, _ := s.users.GetSessionID(models.RefreshToken(rt.RefreshToken))
			// nolint: errcheck
			s.users.Logout(ctx, sessionID)
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.UserAuth{
		AccessToken:  &pb.AccessToken{AccessToken: string(accessToken)},
		RefreshToken: &pb.RefreshToken{RefreshToken: string(refreshToken)},
	}, nil
}
