package fileserver

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"

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
		log.Println("Closing the File Server")
		fs.Transport.Close()
	}()
	log.Println("[INFO] File Server Started Consuming")

	for {
		select {
		case msg := <-fs.Transport.Consume():
			var m p2p.Message
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&m); err != nil {
				log.Println("[Error]", err)
			}
			log.Printf("payload received from %s: %+s\n", msg.From, m.Payload)
		case <-fs.tearDownChan:
			return
		}
	}
}

func (fs *FileServer) handleMessage(msg *p2p.Message) error {
	return nil
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
	buf := bytes.Buffer{}
	tee := io.TeeReader(r, &buf)

	if err := fs.store.Write(key, tee); err != nil {
		return err
	}

	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}

	payload := p2p.Payload{
		Key:  key,
		Data: buf.Bytes(),
	}

	log.Printf("Broadcasting data: %s", buf.Bytes())

	return fs.Broadcast(payload)
}

func (fs *FileServer) Broadcast(payload p2p.Payload) error {
	return fs.Transport.Broadcast(payload)
}

func (fs *FileServer) Dial() {
	fmt.Println("============================")
	for _, addr := range fs.RemoteNodes {
		go func(addr string) {
			log.Printf("Dialing %s\n", addr)
			if err := fs.Transport.Dial(addr); err != nil {
				log.Printf("[ERROR] While Dial in FileServer, err: %s\n", err)
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
