package store

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	opts := NewStoreOpts(CASPathTransformer)
	store := NewStore(*opts)

	data := []byte("some file data")
	key := "testingPath"
	reader := bytes.NewReader(data)
	assert.Nil(t, store.writeStream(key, reader))

	readed, err := store.readStream(key)
	assert.Nil(t, err)

	// study
	readedDate, err := io.ReadAll(readed)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(data), string(readedDate))
	assert.Nil(t, store.deleteStream(key))
}

func TestCASTransformer(t *testing.T) {
	key := "some key"
	pathKey := CASPathTransformer(key)
	assert.Equal(t, PathKey{
		PathName: "ab0d8/e0ce5/8e6fa/9d1b2/30d25/f2ea0/b44a5/1ebd4",
		FileName:  "ab0d8e0ce58e6fa9d1b230d25f2ea0b44a51ebd4",
	}, pathKey)
}
