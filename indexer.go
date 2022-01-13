package portto

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Indexer struct {
	endpoint  string
	ethClient *Client

	repo Repository

	jobs   chan uint64
	errors chan error

	wg *sync.WaitGroup
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
		//todo: should we use buffered channel here?
		jobs:   make(chan uint64, 10),
		errors: make(chan error),
		wg:     &sync.WaitGroup{},
	}, nil
}

func (idx *Indexer) Run(ctx context.Context) {
	idx.wg.Add(MaxWorker)
	for i := 0; i < MaxWorker; i++ {
		go idx.Worker(ctx, i, idx.endpoint)
	}

	for {
		select {
		case <-time.Tick(Interval * time.Second):
			if len(idx.jobs) > 0 {
				log.Printf("got %d job(s) to do, new job skiped", len(idx.jobs))
			} else {
				go idx.addNewBlockToJobQueue(ctx)
			}
		case err := <-idx.errors:
			log.Printf("error received : %v", err)
		case <-ctx.Done():
			return
		}
	}
}

func (idx *Indexer) addNewBlockToJobQueue(ctx context.Context) {
	latestInChain, err := idx.ethClient.GetBlockNumber(context.TODO())
	if err != nil {
		idx.errors <- fmt.Errorf("failed to get latest number on chain, %v", err)
		return
	}

	latestInDB, err := idx.repo.GetLatestNumber()
	if err != nil {
		idx.errors <- fmt.Errorf("failed to get latest number in DB, %v", err)
		return
	}

	var from uint64
	if latestInDB > MinimalBlockNumber {
		from = latestInDB + 1
	} else {
		from = MinimalBlockNumber + 1
	}

	log.Printf("adding new jobs to queue : from %d to %d\n", from, latestInChain)
	for i := from; i <= latestInChain; i++ {
		select {
		case <-ctx.Done():
			log.Println("stop add new job to queue")
			close(idx.jobs)
			return
		default:
			idx.jobs <- i
		}
	}
}

func (idx *Indexer) GetBlock(number uint64) (*Block, error) {
	block, err := idx.repo.FindBlock(number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block from repo, %v", err)
	}

	return block, nil
}

func (idx *Indexer) Worker(ctx context.Context, id int, endpoint string) {
	defer idx.wg.Done()
	client, err := NewClient(endpoint)
	if err != nil {
		idx.errors <- fmt.Errorf("failed to create Client, %v", err)
		return
	}

	for {
		select {
		case number, ok := <-idx.jobs:
			if !ok {
				log.Printf("jobs channel closed, stop worker %d", id)
				return
			}
			blockRaw, err := client.GetBlockByNumber(context.TODO(), number)
			if err != nil {
				idx.errors <- fmt.Errorf("failed to get block, %v", err)
			}

			hashes := make([]string, len(blockRaw.Transactions()))
			for i, transaction := range blockRaw.Transactions() {
				hashes[i] = transaction.Hash().String()
			}
			blockModel := &Block{
				Number:       blockRaw.NumberU64(),
				Hash:         blockRaw.Hash().String(),
				Time:         blockRaw.Time(),
				ParentHash:   blockRaw.ParentHash().String(),
				Transactions: hashes,
			}

			_ = idx.repo.CreateBlock(blockModel)

			if err != nil {
				log.Printf("[DEV] create transactions error, %v", err)
				return
			}
		case <-ctx.Done():
			log.Printf("receive cancel singal, stop worker %d", id)
			return
		}
	}
}

func (idx *Indexer) StopWait() {
	log.Println("waiting for everything stop...")
	idx.wg.Wait()
}

func (idx Indexer) GetNewBlocks(limit int) ([]*Block, error) {
	blocks, err := idx.repo.GetNewBlocks(limit)
	for i := range blocks {
		blocks[i].Transactions = nil
	}

	return blocks, err
}

func (idx *Indexer) GetTransaction(hash string) (*Transaction, error) {
	tx, err := idx.repo.FindTransaction(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to FindTransaction, %v", err)
	}

	if tx == nil {
		tx, err := idx.ethClient.GetTransactionByHash(context.TODO(), hash)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction from chain, %v", err)
		}

		txReceipt, err := idx.ethClient.GetTransactionReceipt(context.TODO(), hash)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction receipt from chain, %v", err)
		}

		logs := make([]Log, len(txReceipt.Logs))
		for i, l := range txReceipt.Logs {
			logs[i] = Log{
				Index: uint64(l.Index),
				Data:  common.BytesToHash(l.Data).String(),
			}
		}

		t := &Transaction{
			Hash:      hash,
			From:      tx.From,
			To:        tx.To,
			Nonce:     tx.Nonce,
			Data:      tx.Data,
			Value:     tx.Value,
			Logs:      logs,
			BlockHash: txReceipt.BlockHash.String(),
		}

		err = idx.repo.CreateTransaction(t)
		if err != nil {
			return nil, fmt.Errorf("failed to store transaction to repo, %v", err)
		}

		return t, nil
	}

	return tx, nil
}
