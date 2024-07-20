package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts {
		ListenAddr: ":6969",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder: DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)
	assert.Equal(t, tr.ListenAddr, ":6969")

	// Basic Server
	assert.Nil(t, tr.ListenAndAccept())
}