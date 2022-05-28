package api

import (
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/users"
	"google.golang.org/grpc"
)

// server implements GophkeeperServer interface.
type server struct {
	users      users.Service
	gophkeeper gophkeeper.Service

	pb.UnimplementedGophkeeperServer
}

func NewServer(u users.Service, g gophkeeper.Service) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterGophkeeperServer(s, &server{
		users:      u,
		gophkeeper: g,
	})
	return s
}
