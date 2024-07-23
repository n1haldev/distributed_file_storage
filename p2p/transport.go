package p2p

// Peer to represent a node
type Peer interface {
	Close() error
}

// Transport to handle communications between nodes. I am using TCP for my project.
type Transport interface {
	ListenAndAccept() 		error
	Dial(string)			error
	Consume() 				<- chan RPC
	Close() 				error
}