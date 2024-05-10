package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listneAddr := ":4000"
	ops := NewTCPTransportOpts(listneAddr, NOPHandshake, NewDefaultDecoder())
	tr := NewTCPTransport(*ops)

	assert.Equal(t, tr.ListenAddr, listneAddr)
	assert.Nil(t, tr.ListenAndAccept())
}
