package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/vanamelnik/gophkeeper/server/api"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/storage/postgres"
	"github.com/vanamelnik/gophkeeper/server/users"
	"google.golang.org/grpc"
)

const (
	cfgFileName = "server_config.yaml"
)

func main() {
	log.Println("Starting GophKeeper server")
	must(loadConfig())

	s, err := postgres.NewStorage(viper.GetString("databaseDSN"), postgres.WithDesctructiveReset())
	must(err)
	defer s.Close()

	u := users.NewService(
		s,
		viper.GetString("tokens.secretKey"),
		viper.GetDuration("tokens.accessTokenDuration"),
		viper.GetDuration("tokens.refreshTokenDuration"),
	)

	g := gophkeeper.NewService(s)
	defer g.Close()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	server := api.NewServer(u, g)
	go runServer(server)

	<-sigint
	log.Println("Shutting down... ")
	server.GracefulStop()
}

func loadConfig() error {
	viper.SetConfigFile(cfgFileName)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Println("Config loaded")
	return nil
}

func runServer(s *grpc.Server) {
	port := viper.GetString("server.port")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("gRPC server is listening at %s", port)
	if err := s.Serve(listen); err != nil {
		log.Printf("gRPC server: %s", err)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
