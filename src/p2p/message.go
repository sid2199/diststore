package p2p

import "net"

type Message struct {
	From net.Addr
	Payload []byte
}

func NewMessage(from net.Addr) *Message {
	return &Message{From:from}
}