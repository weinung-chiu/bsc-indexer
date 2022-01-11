package portto

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type Indexer struct {
	endpoint  string
	ethClient *Client

	repo Repository

	jobs   chan uint64
	errors chan error
}

const MaxWorker = 3

func NewIndexer(endpoint string, repo Repository) (*Indexer, error) {
	c, err := NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		endpoint:  endpoint,
		ethClient: c,
		repo:      repo,
		jobs:      make(chan uint64),
		errors:    make(chan error),
	}, nil
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

func (idx *Indexer) GetBlock(number uint64) (*types.Block, error) {
	block := idx.repo.GetBlock(number)
	if block != nil {
		return block, nil
	}

	// todo : these code should merge to worker
	block, err := idx.ethClient.GetBlockByNumber(context.TODO(), number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number, %v", err)
	}

	err = idx.repo.StoreBlock(block)
	if err != nil {
		idx.errors <- fmt.Errorf("failed to store block, %v", err)
	}

	log.Printf("indexer client got block %d and store to repository\n", block.Number().Uint64())
	// end of code should move

	return block, nil
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

		_ = idx.repo.StoreBlock(block)

		log.Printf("Worker %d got block %d and store to repository\n", id, block.Number().Uint64())
	}
}
