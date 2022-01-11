package portto

import (
	"fmt"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type Repository interface {
	StoreBlock(block *types.Block) error
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		blocks: make([]*types.Block, 0),
	}
}

type InMemoryStore struct {
	blocks []*types.Block
	sync.RWMutex
}

func (s *InMemoryStore) StoreBlock(block *types.Block) error {
	s.Lock()
	defer s.Unlock()
	s.blocks = append(s.blocks, block)

	return nil
}

func (s *InMemoryStore) ShowBlocks() {
	sort.Slice(s.blocks, func(i, j int) bool {
		return s.blocks[i].Number().Uint64() > s.blocks[j].Number().Uint64()
	})

	fmt.Printf("%d block(s) in memory store...\n", len(s.blocks))
	for _, block := range s.blocks {
		fmt.Printf("Block: %d timestamp: %d Hash: %s\n", block.Number(), block.Time(), block.Hash().String())
	}
}
