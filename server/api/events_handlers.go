package server

import (
	"context"
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
	data, err := s.gophkeeper.GetUserData(ctx, userID, r.DataVersion)
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
		Version: r.DataVersion,
		Items:   pbItems,
	}, nil
}

// PublishLocalChanges implements GophkeeperServer interface.
func (s Server) PublishLocalChanges(ctx context.Context, r *pb.PublishLocalChangesRequest) (*emptypb.Empty, error) {
	userID, err := s.users.Authenticate(ctx, models.AccessToken(r.Token.AccessToken))
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Unauthenticated, err.Error())
	}
	dataVersion, err := s.users.GetDataVersion(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, "user not found")
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	if r.DataVersion < dataVersion {
		return &emptypb.Empty{}, status.Error(codes.PermissionDenied, "local data version is out of date")
	}
	// convert events to canonical format
	events := make([]models.Event, 0, len(r.Events))
	for _, e := range r.Events {
		itemID, err := uuid.Parse(e.Item.ItemId.ItemId)
		if err != nil {
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
		}

		if e.Item.CreatedAt != nil || !e.Item.CreatedAt.IsValid() {
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "incorrect CreatedAt field")
		}
		createdAt := e.Item.CreatedAt.AsTime()
		op := models.Operation(e.Operation.String())
		if !op.Valid() {
			return &emptypb.Empty{}, status.Error(codes.Internal,
				fmt.Sprintf("wrong type of operation: %s", op))
		}
		event := models.Event{
			Operation: op,
			Item: models.Item{
				ID:        itemID,
				Version:   e.Item.Version,
				CreatedAt: &createdAt,
				DeletedAt: nil,
				Payload:   parseData(e.Item.Payload),
				Meta:      models.JSONMetadata(e.Item.Metadata.Metadata),
			},
		}
		if e.Item.DeletedAt != nil {
			if !e.Item.DeletedAt.IsValid() {
				return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "incorrect DeletedAt field")
			}
			deletedAt := e.Item.DeletedAt.AsTime()
			event.Item.DeletedAt = &deletedAt
		}
		events = append(events, event)
	}

	// process local events
	if err := s.gophkeeper.PublishUserData(ctx, userID, events); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, status.Error(codes.OK, "accepted")
}

// parseData converts protobuf Item.Data to models.Item.Data format.
func parseData(data interface{}) interface{} {
	switch d := data.(type) {
	case pb.Item_Password:
		return models.PasswordData{
			Password: d.Password.Password,
		}
	case pb.Item_Text:
		return models.TextData{
			Text: d.Text.Text,
		}
	case pb.Item_Blob:
		return models.BinaryData{
			Binary: d.Blob.Data,
		}
	case pb.Item_Card:
		return models.CardData{
			Number:         d.Card.Number,
			CardholderName: d.Card.Name,
			Date:           d.Card.Date,
			CVC:            d.Card.Cvc,
		}
	}
	return nil
}
