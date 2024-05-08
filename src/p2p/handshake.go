package p2p

import "errors"

// HandShake
type Handshake func(Peer) error

func NOPHandshake(Peer) error { return nil}


// Errors

// ErrInvalidHandshake is returned if connection b/w local and
// remote node was not established
var ErrInvalidHandshake = errors.New("invalid handshake")