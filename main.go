package main

import (
	"fmt"
	"log"

	"github.com/sid2199/diststore/src/p2p"
)

func main() {
	fmt.Println("[DISTRIBUTED STORAGE]")
	tcpOpts := p2p.NewTCPTransportOpts(":8080", p2p.NOPHandshake, p2p.NewGOBDecoder())
	tr := p2p.NewTCPTransport(*tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	// select{}
	return
}
