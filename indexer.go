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

// Interval represent add new block to fetch queue (in second)
const Interval = 10

// MinimalBlockNumber when crawl
const MinimalBlockNumber = 14274321

func NewIndexer(endpoint string, repo Repository) (*Indexer, error) {
	c, err := NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		endpoint:  endpoint,
		ethClient: c,
		repo:      repo,
		jobs:      make(chan uint64, 10),
		errors:    make(chan error),
	}, nil
}

func (idx *Indexer) Run() {
	for i := 0; i < MaxWorker; i++ {
		go idx.Worker(i, idx.endpoint)
	}

	for {
		select {
		case <-time.Tick(Interval * time.Second):
			if len(idx.jobs) > 0 {
				log.Printf("got %d job(s) to do, new job skiped", len(idx.jobs))
			} else {
				go idx.addNewBlockToJobQueue()
			}
		case err := <-idx.errors:
			log.Printf("error received : %v", err)
		}
	}
}

func (idx *Indexer) addNewBlockToJobQueue() {
	latestInChain, err := idx.ethClient.GetBlockNumber(context.TODO())
	if err != nil {
		idx.errors <- fmt.Errorf("failed to get latest number on chain, %v", err)
		return
	}

	latestInDB := idx.repo.GetLatestNumber()

	var from uint64
	if latestInDB > MinimalBlockNumber {
		from = latestInDB + 1
	} else {
		from = MinimalBlockNumber + 1
	}

	log.Printf("adding new jobs to queue : from %d to %d\n", from, latestInChain)
	for i := from; i <= latestInChain; i++ {
		idx.jobs <- i
	}
}

func (idx *Indexer) GetBlock(number uint64) (*types.Block, error) {
	block := idx.repo.GetBlock(number)
	if block == nil {
		return nil, fmt.Errorf("block not found")
	}

	return block, nil
}

func (idx *Indexer) Worker(id int, endpoint string) {
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
	}
}
