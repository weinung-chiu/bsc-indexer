package main

import (
	"log"
	"time"

	"portto"
)

func main() {
	rpcEndpoint := "https://bsc-dataseed1.ninicoin.io/"
	repo := portto.NewInMemoryStore()

	i, err := portto.NewIndexer(rpcEndpoint, repo)
	if err != nil {
		log.Fatal("failed to make new indexer, ", err)
	}

	go i.Run()
	block, err := i.GetBlock(14266189)
	if err != nil {
		log.Printf("GetBlock Error, %v", err)
	} else {
		log.Printf("got block %d", block.Number().Uint64())
	}

	time.Sleep(5 * time.Second)

	repo.ShowBlocks()
}
