package p2p

import (
	// "bytes"
	"fmt"
	"net"
	"sync"

	// "github.com/bytedance/sonic/decoder"
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

type TCPTransportOpts struct {
	ListenAddr 		string
	HandshakeFunc	HandshakeFunc
	Decoder			Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listenAddress string
	listener 	net.Listener
	shakeHands 	HandshakeFunc
	decoder		Decoder

	transportLocks 	sync.RWMutex
	peers 		map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
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
		fmt.Printf("New incoming connection from %+v\n", conn)

		go t.handleConnection(conn)
	}
}

type Temp struct {}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Error shaking hands with peer %+v: %s\n", peer, err)
		conn.Close()
		return
	}

	// countDecodeErrors := 0
	// Read Loop
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Println("TCP error decoding message: %s\n", err)
			continue
		}

		fmt.Printf("message: %+v\n", msg)
	}
}

