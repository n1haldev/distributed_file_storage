package p2p

import "net"

// Message can be any arbitrary data that is being sent over a transport between 2 nodes
type RPC struct {
	From 	net.Addr
	Payload []byte
}
