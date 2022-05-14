package server

import (
	pb "github.com/vanamelnik/gophkeeper/proto"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/users"
)

// Server implements GophkeeperServer interface.
type Server struct {
	u users.Service
	g gophkeeper.Service
	pb.UnimplementedGophkeeperServer
}
