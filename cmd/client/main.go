package main

import (
	"fmt"

	"github.com/vanamelnik/gophkeeper/client/repo"
)

func main() {
	_ = repo.Storage{}
	fmt.Println("GophKeeper client is the cool client for GophKeeper service")
}
