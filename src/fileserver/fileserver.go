package fileserver

import (
	"github.com/sid2199/diststore/src/p2p"
	"github.com/sid2199/diststore/src/store"
)

type FileServerOpts struct {
	ListenAddr string
	Root string
	PathTransformer store.PathTransformer
	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *store.Store
}

func NewFileServerOpts(listenAddr string, root string, pathTransformer store.PathTransformer) *FileServerOpts {
	return &FileServerOpts{
		ListenAddr: listenAddr,
		Root: root,
		PathTransformer: pathTransformer,
	}
}


func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store: store.NewStore(*store.NewStoreOpts(opts.Root, opts.PathTransformer)),
	}
}


func (fs *FileServer) Start() {

}