package server

import (
	"context"

	pb "github.com/vanamelnik/gophkeeper/proto"
)

// UpdateData implements GophkeeperServer interface.
func (s Server) UpdateData(ctx context.Context, r *pb.UpdateDataRequest) (*pb.Data, error) {
	userID, err := s.u.Authenticate(r.Token.AccessToken)
	if err != nil {
		//TODO: handle errors
	}
	data, err := s.g.GetUserData(ctx, userID, r.DataVersion)
	if err != nil {
		//TODO: handle errors
	}
	passwords := make([]*pb.Password, 0, len(data.Passwords))
	for _, p := range data.Passwords {
		passwords = append(passwords, &pb.Password{
			ItemId: &pb.ItemID{
				ItemId: p.ItemID.String(),
			},
			Meta: &pb.Metadata{
				Metadata: string(p.Meta),
			},
			Password: p.Password,
		})
	}
	blobs := make([]*pb.Blob, 0, len(data.Blobs))
	for _, b := range data.Blobs {
		blobs = append(blobs, &pb.Blob{
			ItemId: &pb.ItemID{
				ItemId: b.ItemID.String(),
			},
			Meta: &pb.Metadata{
				Metadata: string(b.Meta),
			},
			Data: b.Binary,
		})
	}
	texts := make([]*pb.Text, 0, len(data.Texts))
	for _, t := range data.Texts {
		texts = append(texts, &pb.Text{
			ItemId: &pb.ItemID{
				ItemId: t.ItemID.String(),
			},
			Meta: &pb.Metadata{
				Metadata: string(t.Meta),
			},
			Text: t.Text,
		})
	}
	cards := make([]*pb.Card, 0, len(data.Cards))
	for _, c := range data.Cards {
		cards = append(cards, &pb.Card{
			ItemId: &pb.ItemID{
				ItemId: c.ItemID.String(),
			},
			Meta: &pb.Metadata{
				Metadata: string(c.Meta),
			},
			Number: c.Number,
			Name:   c.CardholderName,
			Date:   c.Date,
			Cvc:    c.CVC,
		})
	}

	return &pb.Data{
		Passwords: passwords,
		Blobs:     blobs,
		Texts:     texts,
		Cards:     cards,
		Version:   data.Version,
	}, nil
}
