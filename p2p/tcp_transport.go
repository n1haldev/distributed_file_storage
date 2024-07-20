package p2p

import (
	"fmt"
	"net"
	"log"
)

// TCPPeer represents node over a TCP connection
type TCPPeer struct {
	// conn is the connection to the peer
	conn net.Conn

	// dial -> true and accept -> false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}	

// Close implements the Peer Interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransport struct {
	TCPTransportOpts
	listener      	net.Listener
	rpcchan	   		chan RPC
	
	// transportLocks sync.RWMutex
	// peers          map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: 	opts,
		rpcchan: 			make(chan RPC),
	}
}

// Implements the Transport interface, which will return read-onlychannel of RPCs
func (t *TCPTransport) Consume() <- chan RPC {
	return t.rpcchan
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.acceptor()

	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

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

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("Error shaking hands with peer %+v: %s\n", peer, err)
		return
	}

	// countDecodeErrors := 0
	// Read Loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error decoding message: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", rpc)
	}
}
