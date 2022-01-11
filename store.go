package portto

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type Repository interface {
	StoreBlock(block *types.Block) error
	GetBlock(number uint64) *types.Block
	GetLatestNumber() uint64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		blocks: make(map[uint64]*types.Block, 0),
	}
}

type InMemoryStore struct {
	blocks map[uint64]*types.Block
	sync.RWMutex

	latestNumber uint64
}

func (s *InMemoryStore) GetLatestNumber() uint64 {
	return s.latestNumber
}

func (s *InMemoryStore) GetBlock(number uint64) *types.Block {
	block, ok := s.blocks[number]
	if !ok {
		return nil
	}

	return block
}

func (s *InMemoryStore) StoreBlock(block *types.Block) error {
	s.Lock()
	defer s.Unlock()
	s.blocks[block.NumberU64()] = block
	if block.NumberU64() > s.latestNumber {
		s.latestNumber = block.NumberU64()
	}
	return nil
}

func (s *InMemoryStore) ShowBlocks() {
	fmt.Printf("%d block(s) in memory store...\n", len(s.blocks))
	for num, block := range s.blocks {
		fmt.Printf("Block: %d timestamp: %d Hash: %s\n", num, block.Time(), block.Hash().String())
	}
}
