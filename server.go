package main

import (
	"log"

	"github.com/n1haldev/distributed_file_storage/p2p"
)

type FileServerOpts struct {
	StorageRoot 		string
	PathTransformFunc 	PathTransformFunc
	Transport 			p2p.Transport
	BootstrapNodes		[]string
}

type FileServer struct {
	FileServerOpts
	store *Store
	quit chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts {
		Root: opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(storeOpts),
		quit: make(chan struct{}),
	}
}

func (fs *FileServer) Stop() {
	close(fs.quit)
}

func (fs *FileServer) loop() {
	defer func() {
		log.Println("Stopping FileServer due to quit signal")
		fs.Transport.Close()
	}()

	for {
		select {
		case msg := <- fs.Transport.Consume():
			log.Println(msg)
		case <-fs.quit:
			return
		}
	}
}

func (fs *FileServer) bootStrapNetwork() error {
	for _, addr := range fs.BootstrapNodes {
		go func(addr string) {
			if err := fs.Transport.Dial(addr); err != nil {
				log.Println(err)
			}
		}(addr)
	}

	return nil
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.bootStrapNetwork() 

	fs.loop()

	return nil
}