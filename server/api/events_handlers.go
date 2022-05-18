package server

import (
	"context"
	"errors"

	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DownloadUserData implements GophkeeperServer interface.
func (s Server) DownloadUserData(ctx context.Context, r *pb.UpdateDataRequest) (*pb.UserData, error) {
	userID, err := s.users.Authenticate(models.AccessToken(r.Token.AccessToken))
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
