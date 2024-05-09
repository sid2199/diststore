package p2p

import (
	"encoding/gob"
	"io"
)


type Decoder interface {
	// study
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct {}

type DefaultDecoder struct {}

func NewGOBDecoder() *GOBDecoder {
	return &GOBDecoder{}
}

func (dec GOBDecoder) Decode(r io.Reader, msg *Message) error {
	// study
	return gob.NewDecoder(r).Decode(msg)
}


func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

func (dec DefaultDecoder) Decode(r io.Reader, msg *Message) error {
	buf := make([]byte, 1024)

	n, err := r.Read(buf)
	if err != nil  {
		return err
	}
	msg.Payload = buf[:n]

	return nil
}
