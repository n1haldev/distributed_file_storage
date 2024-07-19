package p2p

import (
	"net"
	"sync"
	"fmt"
)

// TCPPeer represents node over a TCP connection
type TCPPeer struct {
	// conn is the connection to the peer
	conn net.Conn

	// dial -> true and accept -> false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer {
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener 	net.Listener

	transportLocks 	sync.RWMutex
	peers 		map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	t.listener = ln
	go t.acceptor()
	return nil
}

func (t *TCPTransport) acceptor() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
		}

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	fmt.Printf("New incoming connection from %+v\n", peer)
}

