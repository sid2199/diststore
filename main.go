package main

import (
	"fmt"
	"log"

	"github.com/sid2199/diststore/src/fileserver"
)

func main() {
	fmt.Println("[DISTRIBUTED STORAGE]")

	// cfg := config.Load("")

	fs1 := fileserver.MakeServer(":3000")
	fs2 := fileserver.MakeServer(":4000", ":3000")

	go func() {
		log.Fatalf("[ERROR] %s\n", fs1.Start())
	}()
	log.Fatalf("[ERROR] %s\n", fs2.Start())

	// select {}
	return
}
