package main

import (
	"time"

	"portto"
)

func main() {
	rpcEndpoint := "https://bsc-dataseed1.ninicoin.io/"
	repo := &portto.InMemoryStore{}

	i := portto.NewIndexer(rpcEndpoint, repo)

	go i.Run()
	i.GetBlock(14266189)
	i.GetBlock(14266190)

	time.Sleep(5 * time.Second)
}
