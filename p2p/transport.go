package p2p

import "net"

// Peer to represent a node
type Peer interface {
	RemoteAddr()	net.Addr
	Close() error
}

// Transport to handle communications between nodes. I am using TCP for my project.
type Transport interface {
	ListenAndAccept() 		error
	Dial(string)			error
	Consume() 				<- chan RPC
	Close() 				error
}