package portto

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBlock(*Block) error
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

func (s SQLStore) CreateBlock(b *Block) error {
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
