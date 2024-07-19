package main

import (
	"log"
	"github.com/n1haldev/distributed_file_storage/p2p"
	)

func main() {
	tr := p2p.NewTCPTransport(":6969")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("Error listening and accepting connections: %s", err)
	}

	select {}
}