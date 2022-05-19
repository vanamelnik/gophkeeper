package server

import (
	"context"
	"errors"

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
func (s Server) DownloadUserData(ctx context.Context, r *pb.UserDataRequest) (*pb.UserData, error) {
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

	// Convert to awful protobuf format.
	pbItems := make([]*pb.Item, 0, len(data.Items))
	for _, item := range data.Items {
		pbItem := pb.Item{
			ItemId: &pb.ItemID{
				ItemId: item.ID.String(),
			},
			Metadata: &pb.Metadata{
				Metadata: string(item.Meta),
			},
		}
		switch body := item.Data.(type) {
		case models.TextData:
			text := pb.Item_Text{Text: &pb.Text{Text: body.Text}}
			pbItem.Data = &text
		case models.BinaryData:
			blob := pb.Item_Blob{Blob: &pb.Blob{Data: body.Binary}}
			pbItem.Data = &blob
		case models.PasswordData:
			password := pb.Item_Password{Password: &pb.Password{Password: body.Password}}
			pbItem.Data = &password
		case models.CardData:
			card := pb.Item_Card{Card: &pb.Card{
				Number: body.Number,
				Name:   body.CardholderName,
				Date:   body.Date,
				Cvc:    body.CVC,
			}}
			pbItem.Data = &card
		}
		pbItems = append(pbItems, &pbItem)
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
		return &emptypb.Empty{}, status.Error(codes.PermissionDenied, "data version is out of date")
	}
	// convert events to canonical format
	events := make([]models.Event, 0, len(r.Events))
	for _, e := range r.Events {
		itemID, err := uuid.Parse(e.Item.ItemId.ItemId)
		if err != nil {
			return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
		}

		if e.Item.CreatedAt != nil || !e.Item.CreatedAt.IsValid() {
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "incorrect CreatedAt field")
		}
		createdAt := e.Item.CreatedAt.AsTime()
		event := models.Event{
			Operation: models.Operation(e.Operation.String()),
			Item: models.Item{
				ID:        itemID,
				CreatedAt: &createdAt,
				DeletedAt: nil,
				Data:      parseData(e.Item.Data),
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

	return &emptypb.Empty{}, nil
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
