package portto

import (
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

type Repository interface {
	StoreBlock(block *types.Block) error
	GetBlock(number uint64) *types.Block
	GetLatestNumber() uint64
}

type SQLStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

func (s SQLStore) StoreBlock(b *types.Block) error {
	model := &block{
		Number:     b.NumberU64(),
		Hash:       b.Hash().String(),
		Time:       b.Time(),
		ParentHash: b.ParentHash().String(),
	}
	result := s.db.Create(model)
	if result.Error != nil {
		return fmt.Errorf("failed to create MySQL record")
	}

	return nil
}

func (s SQLStore) GetBlock(number uint64) *types.Block {
	b := &block{}

	result := s.db.First(b, number)

	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}

	if result.Error != nil {
		// todo: should return error here
		log.Fatal("failed to get block, ", result.Error)
	}

	h := &types.Header{
		ParentHash: common.HexToHash(b.ParentHash),
		Number:     big.NewInt(int64(b.Number)),
		Time:       b.Time,
	}
	blockWithHeader := types.NewBlockWithHeader(h)

	// todo: should return custom struct instead of *types.Block
	return blockWithHeader
}

func (s SQLStore) GetLatestNumber() uint64 {
	b := &block{}
	result := s.db.Last(b)

	if result.Error == gorm.ErrRecordNotFound {
		return 0
	}

	if result.Error != nil {
		// todo: should return error here
		log.Fatal("failed to get latest number, ", result.Error)
	}

	return b.Number
}

type block struct {
	Number     uint64
	Hash       string
	Time       uint64
	ParentHash string
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
