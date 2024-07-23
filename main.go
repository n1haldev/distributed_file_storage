package main

import (
	"log"
	"time"

	"github.com/n1haldev/distributed_file_storage/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		// TODO: onPeer func
		// OnPeer: ,
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts {
		StorageRoot: "nihal_network",
		PathTransformFunc: CasPathTransformFunc,
		Transport: tcpTransport,
		BootstrapNodes: []string{":4000"},
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}