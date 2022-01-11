package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	fmt.Println("Press Ctrl+C to interrupt...")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	for {
		select {
		case <-done:
			repo.ShowBlocks()
			fmt.Println("Bye Bye...")
			os.Exit(1)
		}
	}
}
