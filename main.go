package main

import (
	"log"

	"github.com/n1haldev/distributed_file_storage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts {
		ListenAddr: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts {
		StorageRoot: listenAddr + "_network",
		PathTransformFunc: CasPathTransformFunc,
		Transport: tcpTransport,
		BootstrapNodes: nodes,
	}

	fs := NewFileServer(fileServerOpts)

	return fs
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()
}