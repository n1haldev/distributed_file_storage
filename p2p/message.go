package p2p

// Message can be any arbitrary data that is being sent over a transport between 2 nodes
type P2PMessage struct {
	Payload []byte
}