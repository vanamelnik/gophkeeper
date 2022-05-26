package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DownloadUserData implements GophkeeperServer interface.
func (s Server) DownloadUserData(ctx context.Context, r *pb.DownloadUserDataRequest) (*pb.UserData, error) {
	userID, err := s.users.Authenticate(ctx, models.AccessToken(r.Token.AccessToken))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	versionMap := make(map[uuid.UUID]uint64)
	if err := json.Unmarshal([]byte(r.VersionMap), &versionMap); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	data, err := s.gophkeeper.GetUserData(ctx, userID, versionMap)
	if err != nil {
		if errors.Is(err, gophkeeper.ErrVersionUpToDate) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbItems := make([]*pb.Item, 0, len(data.Items))
	for _, item := range data.Items {
		pbItems = append(pbItems, models.ItemToPb(item))
	}

	return &pb.UserData{
		DataVersion: data.Version,
		Items:       pbItems,
	}, nil
}

// PublishLocalChanges implements GophkeeperServer interface.
func (s Server) PublishLocalChanges(ctx context.Context, r *pb.PublishLocalChangesRequest) (*emptypb.Empty, error) {
	userID, err := s.users.Authenticate(ctx, models.AccessToken(r.Token.AccessToken))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	dataVersion, err := s.users.GetDataVersion(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if r.DataVersion < dataVersion {
		return nil, status.Error(codes.PermissionDenied, "local data version is out of date")
	}
	// convert events to canonical format
	events := make([]models.Event, 0, len(r.Events))
	for _, e := range r.Events {
		op := models.Operation(e.Operation.String())
		if !op.Valid() {
			return nil, status.Error(codes.InvalidArgument,
				fmt.Sprintf("wrong type of operation: %s", op))
		}
		item, err := models.PbToItem(e.Item)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		event := models.Event{
			Operation: op,
			Item:      item,
		}
		events = append(events, event)
	}

	// process local events
	if err := s.gophkeeper.PublishUserData(ctx, userID, events); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// WhatsNew implements GophkeeperServer interface.
func (s Server) WhatsNew(ctx context.Context, r *pb.WhatsNewRequest) (*emptypb.Empty, error) {
	userID, err := s.users.Authenticate(ctx, models.AccessToken(r.Token.AccessToken))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	dataVersion, err := s.users.GetDataVersion(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if r.DataVersion != dataVersion {
		return nil, status.Error(codes.PermissionDenied, "out of date")
	}

	return &emptypb.Empty{}, nil
}
