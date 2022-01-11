package portto

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Indexer struct {
	endpoint string

	repo Repository

	jobs   chan uint64
	errors chan error
}

const MaxWorker = 3

func NewIndexer(endpoint string, repo Repository) *Indexer {
	return &Indexer{
		endpoint: endpoint,
		repo:     repo,
		jobs:     make(chan uint64),
		errors:   make(chan error),
	}
}

func (idx *Indexer) Run() {
	for i := 0; i < MaxWorker; i++ {
		go idx.Worker(i, idx.endpoint)
	}

	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			log.Println("tick...")
		case err := <-idx.errors:
			log.Printf("error received : %v", err)
		}
	}
}

func (idx *Indexer) GetBlock(number uint64) {
	idx.jobs <- number
}

func (idx *Indexer) Worker(id int, endpoint string) {
	log.Printf("Worker %d generated\n", id)
	client, err := NewClient(endpoint)
	if err != nil {
		idx.errors <- fmt.Errorf("failed to create Client, %v", err)
		return
	}

	for number := range idx.jobs {
		block, err := client.GetBlockByNumber(context.TODO(), number)
		if err != nil {
			idx.errors <- fmt.Errorf("failed to get block, %v", err)
		}

		log.Printf("Worker %d got block %d\n", id, block.Number().Uint64())
		log.Printf("Block Hash : %s", block.Hash().String())
	}
}
