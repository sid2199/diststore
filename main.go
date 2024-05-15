package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sid2199/diststore/src/config"
	"github.com/sid2199/diststore/src/fileserver"
	"github.com/sid2199/diststore/src/p2p"
	"github.com/sid2199/diststore/src/store"
)

func main() {
	fmt.Println("[DISTRIBUTED STORAGE]")
	// tcpOpts := p2p.NewTCPTransportOpts(":8080", p2p.NOPHandshake, p2p.NewDefaultDecoder(), p2p.DefaultPeerValidation)
	// tr := p2p.NewTCPTransport(*tcpOpts)

	// // start consuming the meggage chan before accepting the connections
	// go func() {
	// 	for {
	// 		fmt.Printf("%+v\n", <-tr.Consume())
	// 	}
	// }()

	// if err := tr.ListenAndAccept(); err != nil {
	// 	log.Fatal(err)
	// }

	cfg := config.Load("")

	fsOpts := fileserver.NewFileServerOpts("8080_nw", store.CASPathTransformer,
		p2p.NewTCPTransport(*p2p.NewTCPTransportOpts(
			cfg.ListenAddr, p2p.NOPHandshake, p2p.NewDefaultDecoder(), nil),
		), []string{":4000", ":5000"})
	fs := fileserver.NewFileServer(*fsOpts)

	go time.AfterFunc(time.Second*5, func() {
		fs.Close()
	})

	if err := fs.Start(); err != nil {
		log.Fatalf("[ERROR] %s\n", err)
	}

	// select {}
	return
}
