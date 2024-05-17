package p2p

import (
	"errors"
	"log"
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
	ListenAddr     string
	Handshake      Handshake
	Decoder        Decoder
	PeerValidation func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	// study
	listner net.Listener
	msgChan chan Message

	tearDownChan chan int

	// mu    sync.RWMutex
	// peers map[net.Addr]Peer
}

func NewTCPTransportOpts(listnerAddr string, handshake Handshake, decoder Decoder, peerValidation func(Peer) error) *TCPTransportOpts {
	return &TCPTransportOpts{
		ListenAddr:     listnerAddr,
		Handshake:      handshake,
		Decoder:        decoder,
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
		msgChan:          make(chan Message),
		tearDownChan:     make(chan int, 1),
	}
}

func (t *TCPTransport) Consume() <-chan Message {
	return t.msgChan
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.start()
	log.Println("[Started listening and accepting connections]")
	return nil
}

func (t *TCPTransport) start() {
	for {
		// study
		conn, err := t.listner.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			log.Printf("[ERROR] While Accepting: %s\n", err)
		}
		go t.handleConn(conn, false)
	}
}

// Dial implements Transport Interface
func (t *TCPTransport) Dial(addr string) error {
	// TODO: study
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("[ERROR] While Dial TCP Transport, err: %s\n", err)
		return err
	}
	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, outBound bool) {
	var err error
	defer func() {
		log.Println("Closing peer connection")
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outBound)
	log.Printf("new conection: %+v | %+v\n", conn, peer)

	if err := t.Handshake(peer); err != nil {
		conn.Close()
		log.Printf("[ERROR] TCP handshake error: %s\n", err)
		return
	}
	log.Println("Handshake completed")

	if t.PeerValidation != nil {
		if err = t.PeerValidation(peer); err != nil {
			log.Printf("Peer validation failed: %v\n", err)
			return
		}
	}
	log.Println("Peer Validated Successfully")

	msg := NewMessage(conn.RemoteAddr())
	for {
		conn.Read(msg.Payload)
		// TODO: only to drop conn if the connectin is closed by the foreign entity
		if err = t.Decoder.Decode(conn, msg); err != nil {
			log.Printf("[ERROR] TCP error: %s\n", err)
			return
		}
		// TODO: could also send pointer to the message
		t.msgChan <- *msg
	}
}

func (t *TCPTransport) Close() error {
	log.Println("Closing TCP Tranport")
	t.tearDownChan <- 1
	return t.listner.Close()
}
