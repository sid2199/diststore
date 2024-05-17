package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

type StoreOpts struct {
	Root            string
	PathTransformer PathTransformer
}

var DefaultStoreRoot string = "SidStore"

type Store struct {
	StoreOpts
}
type PathKey struct {
	PathName string
	FileName string
}

func (path *PathKey) FullFileName() string {
	return path.PathName + "/" + path.FileName
}

type PathTransformer func(string) PathKey

var DefaultPathTransformer PathTransformer = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

func CASPathTransformer(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:]) // Note: to convert an array to a slice use arr[:]

	blockSize := 5
	sliceLen := len(hashString) / blockSize
	path := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		path[i] = hashString[from:to]
	}

	return PathKey{
		PathName: strings.Join(path, "/"),
		FileName: hashString,
	}
}

func NewStoreOpts(root string, pathTransformer PathTransformer) *StoreOpts {
	if pathTransformer == nil {
		pathTransformer = DefaultPathTransformer
	}
	if root == "" {
		root = DefaultStoreRoot
	}
	return &StoreOpts{
		Root:            root,
		PathTransformer: pathTransformer,
	}
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Read(key string) (io.Reader, error) {
	return s.readStream(key)
}

func (s *Store) readStream(key string) (io.Reader, error) {
	pathKey := s.PathTransformer(key)
	fullFileNameWithRoot := s.Root + "/" + pathKey.FullFileName()

	f, err := os.Open(fullFileNameWithRoot)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	log.Printf("Read %d bytes from disk: %s\n", n, fullFileNameWithRoot)
	return buf, nil
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformer(key)
	fullFileNameWithRoot := s.Root + "/" + pathKey.FullFileName()
	pathNameWithRoot := s.Root + "/" + pathKey.PathName

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(fullFileNameWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Written %d bytes to disk: %s\n", n, fullFileNameWithRoot)
	return nil
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformer(key)
	fullFileNameWithRoot := s.Root + "/" + pathKey.FullFileName()

	_, err := os.Stat(fullFileNameWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	return s.deleteStream(key)
}

func (s *Store) deleteStream(key string) error {
	pathKey := s.PathTransformer(key)
	pathNameWithRoot := s.Root + "/" + strings.Split(pathKey.PathName, "/")[0]

	defer func() {
		log.Printf("deteted from disk: %s\n", pathNameWithRoot)
	}()
	return os.RemoveAll(pathNameWithRoot)
}
