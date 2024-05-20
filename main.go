package main

import (
	"bytes"
	"time"

	"github.com/sid2199/diststore/src/fileserver"
	"github.com/sid2199/diststore/src/logger"
)

var log = logger.Logger

func main() {
	log.Info.Println("[DISTRIBUTED STORAGE]")


	fs1 := fileserver.MakeServer(":3000")
	fs2 := fileserver.MakeServer(":4000", ":3000")

	go func() {
		log.Error.Fatalf("[ERROR] %s\n", fs1.Start())
	}()
	time.Sleep(time.Second * 1)
	go func() {
		log.Error.Fatalf("[ERROR] %s\n", fs2.Start())
	}()
	time.Sleep(time.Second * 1)

	data := bytes.NewReader([]byte("very very big file..."))
	fs2.Store("my private data", data)

	select {}
	return
}
