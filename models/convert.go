package models

import (
	"errors"

	"github.com/google/uuid"
	pb "github.com/vanamelnik/gophkeeper/proto"
)

// PbToItem converts protobuf Item to canonical Item struct.
func PbToItem(item *pb.Item) (Item, error) {
	itemID, err := uuid.Parse(item.ItemId.ItemId)
	if err != nil {
		return Item{}, err
	}
	createdAt := item.CreatedAt.AsTime()
	result := Item{
		ID:        itemID,
		Version:   item.Version,
		CreatedAt: &createdAt,
		DeletedAt: nil,
		Meta:      JSONMetadata(item.Metadata.Metadata),
	}
	if item.DeletedAt != nil {
		if !item.DeletedAt.IsValid() {
			return Item{}, errors.New("incorrect DeletedAt field")
		}
		deletedAt := item.DeletedAt.AsTime()
		result.DeletedAt = &deletedAt
	}

	switch pl := item.Payload.(type) {
	case *pb.Item_Blob:
		result.Payload = BinaryData{pl.Blob.Data}
	case *pb.Item_Text:
		result.Payload = TextData{pl.Text.Text}
	case *pb.Item_Password:
		result.Payload = PasswordData{pl.Password.Password}
	case *pb.Item_Card:
		result.Payload = CardData{
			Number:         pl.Card.Number,
			CardholderName: pl.Card.Name,
			Date:           pl.Card.Date,
			CVC:            pl.Card.Cvc,
		}
	default:
		return Item{}, errors.New("unknown type of the payload")
	}
	return result, nil
}

// ItemToPb converts canonical Item to protobuf Item
func ItemToPb(item Item) *pb.Item {
	pbItem := pb.Item{
		ItemId: &pb.ItemID{
			ItemId: item.ID.String(),
		},
		Metadata: &pb.Metadata{
			Metadata: string(item.Meta),
		},
	}
	switch body := item.Payload.(type) {
	case TextData:
		text := pb.Item_Text{Text: &pb.Text{Text: body.Text}}
		pbItem.Payload = &text
	case BinaryData:
		blob := pb.Item_Blob{Blob: &pb.Blob{Data: body.Binary}}
		pbItem.Payload = &blob
	case PasswordData:
		password := pb.Item_Password{Password: &pb.Password{Password: body.Password}}
		pbItem.Payload = &password
	case CardData:
		card := pb.Item_Card{Card: &pb.Card{
			Number: body.Number,
			Name:   body.CardholderName,
			Date:   body.Date,
			Cvc:    body.CVC,
		}}
		pbItem.Payload = &card
	}
	return &pbItem
}
