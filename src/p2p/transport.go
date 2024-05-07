package p2p


// Peer represents the remote node
type Peer interface{}



// Transport handle the communication b/w the nodes in a network.
// Communication can be of TCP, UDP, webSocket, etc...
type Transport interface{
	ListenAndAccept() error
}


