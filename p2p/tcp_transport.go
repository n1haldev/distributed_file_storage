package p2p

import (
	"fmt"
	"net"
	"log"
	"errors"
	// "reflect"
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
	OnPeer		  func(Peer) error
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

func (t *TCPTransport) Close() error {
	return t.listener.Close()
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
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
		}
		fmt.Printf("New incoming connection from %+v\n", conn)

		go t.handleConnection(conn, false)
	}
}

// Dial implements the Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConnection(conn, true)

	return nil
}

func (t *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	var err error
	
	defer func() {
		fmt.Printf("Closing connection: %s", err)
		conn.Close()
	}()
	
	peer := NewTCPPeer(conn, outbound)
		

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	// countDecodeErrors := 0
	// Read Loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc); 

		if err != nil {
			// fmt.Printf("TCP Read error decoding message: %s\n", err)
			// continue
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpcchan <- rpc

		fmt.Printf("message: %+v\n", rpc)
	}
}
