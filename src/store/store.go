package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type PathKey struct {
	PathName string
	FileName string
}

func (path *PathKey) FullFileName() string {
	return path.PathName + "/" + path.FileName
}

type PathTransformer func(string) PathKey

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

type StoreOpts struct {
	PathTransformer PathTransformer
}

type Store struct {
	StoreOpts
}

func NewStoreOpts(pathTransformer PathTransformer) *StoreOpts {
	return &StoreOpts{
		PathTransformer: pathTransformer,
	}
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) readStream(key string) (io.Reader, error) {
	pathKey := CASPathTransformer(key)

	f, err := os.Open(pathKey.FullFileName())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Read %d bytes from disk: %s\n", n, pathKey.FullFileName())
	return buf, nil
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformer(key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(pathKey.FullFileName())
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes to disk: %s\n", n, pathKey.FullFileName())
	return nil
}

func (s *Store) deleteStream(key string) error {
	pathKey := CASPathTransformer(key)
	dir := strings.Split(pathKey.PathName, "/")[0]
	defer func() {
		fmt.Printf("deteted from disk: %s", dir)
	}()
	return os.RemoveAll(dir)
}
