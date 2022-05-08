package server

import (
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/users"
)

type Server struct {
	u users.Service
	g gophkeeper.Service
}
