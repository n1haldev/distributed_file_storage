package main

import (
	"log"
	"github.com/n1haldev/distributed_file_storage/p2p"
	)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Decoder: p2p.GOBDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("Error listening and accepting connections: %s", err)
	}

	select {}
}