package server

import (
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/users"
)

// Server implements GophkeeperServer interface.
type Server struct {
	users      users.Service
	gophkeeper gophkeeper.Service

	pb.UnimplementedGophkeeperServer
}
