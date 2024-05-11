package p2p

import (
	"fmt"
	"net"
	// "sync"
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
	PeerValidation func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	// study
	listner net.Listener
	msgChan chan Message

	// mu    sync.RWMutex
	// peers map[net.Addr]Peer
}

func NewTCPTransportOpts(listnerAddr string, handshake Handshake, decoder Decoder, peerValidation func(Peer) error) *TCPTransportOpts {
	return &TCPTransportOpts{
		ListenAddr: listnerAddr,
		Handshake:  handshake,
		Decoder:    decoder,
		PeerValidation: peerValidation,
	}
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outBound: outBound,
	}
}

func (peer *TCPPeer) Close() error {
	return peer.conn.Close()
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		msgChan: make(chan Message),
	}
}

func (t *TCPTransport) Consume() <-chan Message{
	return t.msgChan
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.start()
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
	var err error
	defer func() {
		fmt.Println("Closing peer connection")
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	fmt.Printf("new conection: %+v | %+v\n", conn, peer)

	if err := t.Handshake(peer); err != nil {
		conn.Close()
		fmt.Printf("[ERROR] TCP handshake error: %s\n", err)
		return
	}
	fmt.Println("Handshake completed")

	if t.PeerValidation != nil {
		if err = t.PeerValidation(peer); err != nil {
			fmt.Printf("Peer validation failed: %v\n", err)
			return
		}
	}
	fmt.Println("Peer Validated Successfully")

	msg := NewMessage(conn.RemoteAddr())
	for {
		conn.Read(msg.Payload)
		// TODO: only to drop conn if the connectin is closed by the foreign entity
		if err = t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("[ERROR] TCP error: %s\n", err)
			return
		}
		// TODO: could also send pointer to the message
		t.msgChan <- *msg
	}
}
