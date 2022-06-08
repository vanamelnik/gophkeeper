package main

import (
	"context"
	"log"
	"time"

	"github.com/vanamelnik/gophkeeper/client"
	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/client/ui"
	pb "github.com/vanamelnik/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddr   = ":3000"
	syncInterval = 5 * time.Second
	sendInterval = 10 * time.Second
)

func main() {
	ctx := context.Background()
	repo := repo.New()
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	must(err)
	defer conn.Close()
	pbClient := pb.NewGophkeeperClient(conn)

	c := client.New(ctx,
		pbClient,
		syncInterval,
		sendInterval,
		repo,
		ui.ResolveConflict)

	userInterface := ui.NewUI(ctx, repo, c, pbClient)

	userInterface.Run()
	if c.IsLoggedIn() {
		must(c.LogOut()) // TODO: когда будет реализовано сохранение сеанса на диск, то не будет необходимости log out при выходе
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
