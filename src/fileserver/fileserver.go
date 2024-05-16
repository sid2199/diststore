package fileserver

import (
	"fmt"
	"io"

	"github.com/sid2199/diststore/src/p2p"
	"github.com/sid2199/diststore/src/store"
)

type FileServerOpts struct {
	Root            string
	PathTransformer store.PathTransformer
	Transport       p2p.Transport
	RemoteNodes     []string
}

type FileServer struct {
	FileServerOpts
	store        *store.Store
	tearDownChan chan int
}

func NewFileServerOpts(root string, pathTransformer store.PathTransformer, transport p2p.Transport, remoteNodes []string) *FileServerOpts {
	return &FileServerOpts{
		Root:            root,
		PathTransformer: pathTransformer,
		Transport:       transport,
		RemoteNodes:     remoteNodes,
	}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store:          store.NewStore(*store.NewStoreOpts(opts.Root, opts.PathTransformer)),
		tearDownChan:   make(chan int, 1),
	}
}

func (fs *FileServer) Close() {
	close(fs.tearDownChan)
}

func (fs *FileServer) Consume() {
	defer func() {
		fmt.Println("Closing the File Server")
		fs.Transport.Close()
	}()
	fmt.Println("Started Consuming")

	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.tearDownChan:
			return
		}
	}
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.Dial()
	fs.Consume()

	return nil
}

func (fs *FileServer) Store(key string, r io.Reader) error {
	return fs.store.Write(key, r)
}

func (fs *FileServer) Dial() {
	for _, addr := range fs.RemoteNodes {
		go func(addr string) {
			fmt.Printf("Dialing %s\n", addr)
			if err := fs.Transport.Dial(addr); err != nil {
				fmt.Printf("[ERROR] While Dial in FileServer, err: %s\n", err)
			}
		}(addr)
	}
}

func MakeServer(listenAddr string, nodes ...string) *FileServer {
	fsOpts := NewFileServerOpts(listenAddr+"_nw", store.CASPathTransformer,
		p2p.NewTCPTransport(*p2p.NewTCPTransportOpts(
			listenAddr, p2p.NOPHandshake, p2p.NewDefaultDecoder(), nil),
		), nodes)
	return NewFileServer(*fsOpts)
}