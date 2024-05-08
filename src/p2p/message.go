package p2p

type Message struct {
	Payload []byte
}

func NewMessage() *Message {
	return &Message{}
}