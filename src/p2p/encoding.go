package p2p

import (
	"encoding/gob"
	"io"
)


type Decoder interface {
	// study
	Decode(io.Reader, any) error
}

type GOBDecoder struct {}

func NewGOBDecoder() *GOBDecoder {
	return &GOBDecoder{}
}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	// study
	return gob.NewDecoder(r).Decode(v)
}