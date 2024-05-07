package main

import (
	"fmt"
	"log"

	"github.com/sid2199/diststore/src/p2p"
)

func main() {
	fmt.Println("[DISTRIBUTED STORAGE]")
	tr := p2p.NewTCPTransporter("127.0.0.1:8080")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	// select{}
	return
}
