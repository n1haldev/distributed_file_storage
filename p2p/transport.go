package p2p

// Peer to represent a node
type Peer interface {

}

// Transport to handle communications between nodes. I am using TCP for my project.
type Transport interface {
	ListenAndAccept() error
}