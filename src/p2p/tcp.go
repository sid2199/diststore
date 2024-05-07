package p2p

import (
	"net"
	"sync"
	"fmt"
)


// TCPPeer represents the remote node over a TCP connection
type TCPPeer struct {
	conn net.Conn

	// if we dial a connection to another peer then it is outBound
	// else if a connection is accecpted it is inBound
	outBound bool
}

type TCPTransporter struct {
	listnerAddr	string
	listner		net.Listener

	mu sync.RWMutex
	peers map[net.Addr] Peer
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outBound: outBound,
	}
}

func NewTCPTransporter(listnerAddr string) *TCPTransporter {
	return &TCPTransporter{
		listnerAddr: listnerAddr,
	}
}

func (t *TCPTransporter) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.listnerAddr)
	if err != nil {
		return err
	}
	t.start()
	return nil
}

func (t *TCPTransporter) start() {
	for {
		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("[ERROR] While Accepting: %s\n", err)
		}
		go t.handleConn(conn)
	}
}


func (t *TCPTransporter) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new conection: %+v | %+v\n", conn, peer)
}
