package portto

import (
	"context"
	"fmt"
	"log"
	"time"
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

func (idx *Indexer) GetBlock(number uint64) (*Block, error) {
	block, err := idx.repo.GetBlock(number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block from repo, %v", err)
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
		blockRaw, err := client.GetBlockByNumber(context.TODO(), number)
		if err != nil {
			idx.errors <- fmt.Errorf("failed to get block, %v", err)
		}

		blockModel := &Block{
			Number:     blockRaw.NumberU64(),
			Hash:       blockRaw.Hash().String(),
			Time:       blockRaw.Time(),
			ParentHash: blockRaw.ParentHash().String(),
		}

		_ = idx.repo.StoreBlock(blockModel)
	}
}

func (idx Indexer) GetNewBlocks(limit int) ([]*Block, error) {
	return idx.repo.GetNewBlocks(limit)
}
