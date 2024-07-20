package main

import (
	"log"
	"fmt"
	"github.com/n1haldev/distributed_file_storage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("Custom logic with peer outside of TCPTransport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts {
		ListenAddr: ":3000",		
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg:= <- tr.Consume()
			fmt.Printf("Received message: %+v\n", string(msg.Payload))
		}
	}();

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("Error listening and accepting connections: %s", err)
	}

	select {}
}