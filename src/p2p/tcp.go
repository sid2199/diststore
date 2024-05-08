package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP connection
type TCPPeer struct {
	// study
	conn net.Conn

	// if we dial a connection to another peer then it is outBound
	// else if a connection is accecpted it is inBound
	outBound bool
}

type TCPTransportOpts struct {
	ListenAddr string
	Handshake  Handshake
	Decoder    Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	// study
	listner     net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransportOpts(listnerAddr string, handshake Handshake, decoder Decoder) *TCPTransportOpts {
	return &TCPTransportOpts{
		ListenAddr: listnerAddr,
		Handshake: handshake,
		Decoder: decoder,
	}
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outBound: outBound,
	}
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	t.start()
	return nil
}

func (t *TCPTransport) start() {
	for {
		// study
		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("[ERROR] While Accepting: %s\n", err)
		}
		go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new conection: %+v | %+v\n", conn, peer)

	if err := t.Handshake(peer); err != nil {
		conn.Close()
		fmt.Printf("[ERROR] TCP handshake error: %s\n", err)
		return
	}
	fmt.Println("Handshake completed")

	msg := NewMessage()
	for {
		conn.Read(msg.Payload)
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("[ERROR] TCP error: %s\n", err)
		}
		fmt.Println("[INFO] message: %+v\n", msg)
	}
}
