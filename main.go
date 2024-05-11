package main

import (
	"fmt"
	"log"

	"github.com/sid2199/diststore/src/p2p"
)

func main() {
	fmt.Println("[DISTRIBUTED STORAGE]")
	tcpOpts := p2p.NewTCPTransportOpts(":8080", p2p.NOPHandshake, p2p.NewDefaultDecoder(), p2p.DefaultPeerValidation)
	tr := p2p.NewTCPTransport(*tcpOpts)

	// start consuming the meggage chan before accepting the connections
	go func() {
		for {
			fmt.Printf("%+v\n", <-tr.Consume())
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
	return
}
