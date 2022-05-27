package main

import (
	"fmt"
	"log"
	"time"
)

const (
	serverAddr   = ":3000"
	syncInterval = 5 * time.Second
	sendInterval = 10 * time.Second
)

func main() {
	// repo := repo.New()
	// conn, err := grpc.Dial(serverAddr)
	// must(err)
	// defer conn.Close()
	// pbClient := pb.NewGophkeeperClient(conn)

	// if no stored active session
	//	 sign in / email

	// c, err := client.New(context.Background(),
	// 	pbClient,
	// 	syncInterval,
	// 	sendInterval,
	// 	repo)
	// user loop

	fmt.Println("GophKeeper client is the cool client for GophKeeper service")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
