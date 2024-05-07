package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listnerAddr := ":4000"
	tr := NewTCPTransporter(listnerAddr)

	assert.Equal(t, tr.listnerAddr, listnerAddr)
	assert.Nil(t, tr.ListenAndAccept())
}
