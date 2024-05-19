package p2p

import (
	"log"
	"net"
)

// Peer represents the remote node
type Peer interface{
	Send([]byte) error
	// RemoteAddr() net.Addr
	// Close() error

	// Implemented Interface embadding
	// instead of making a method for every method or variable
	net.Conn
}

func DefaultPeerValidation(peer Peer) error {
	log.Printf("Validating peer: %v\n", peer)
	return nil
}


// Transport handle the communication b/w the nodes in a network.
// Communication can be of TCP, UDP, webSocket, etc...
type Transport interface{
	ListenAndAccept() error
	Consume() <-chan Message
	Close () error
	Dial(string) error
	Broadcast(Payload) error
}


