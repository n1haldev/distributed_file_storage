package p2p

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestTCPTransport(t *testing.T) {
// 	listenAddr := ":6969"
// 	tr := NewTCPTransport(listenAddr, GOBDecoder{}, NOPHandshakeFunc)
// 	assert.Equal(t, listenAddr, tr.listenAddress)

// 	// Basic Server
// 	assert.Nil(t, tr.ListenAndAccept())

// 	select{}
// }