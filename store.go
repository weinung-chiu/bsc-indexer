package portto

import (
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

type Repository interface {
	StoreBlock(*Block) error
	GetBlock(number uint64) (*Block, error)
	GetLatestNumber() (uint64, error)
	GetNewBlocks(limit int) ([]*Block, error)
}

type SQLStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

func (s SQLStore) GetNewBlocks(limit int) ([]*Block, error) {
	var blocks []*Block

	result := s.db.Limit(limit).Order("number desc").Find(&blocks)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		log.Fatal("failed to get new blocks, ", result.Error)
	}

	return blocks, nil
}

func (s SQLStore) StoreBlock(b *Block) error {
	result := s.db.Create(b)
	if result.Error != nil {
		return fmt.Errorf("failed to create MySQL record")
	}

	return nil
}

func (s SQLStore) GetBlock(number uint64) (*Block, error) {
	b := &Block{}

	result := s.db.First(b, number)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get block, %v", result.Error)
	}

	return b, nil
}

func (s SQLStore) GetLatestNumber() (uint64, error) {
	b := &Block{}
	result := s.db.Last(b)

	if result.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}

	if result.Error != nil {
		return 0, result.Error
	}

	return b.Number, nil
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
