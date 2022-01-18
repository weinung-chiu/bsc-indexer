package portto

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetLatestNumber() (uint64, error)
	GetNewBlocks(limit int) ([]*Block, error)

	GetUnconfirmedBlocks() ([]*Block, error)
	ConfirmBlocks([]*Block) error

	CreateBlock(*Block) error
	FindBlock(number uint64) (*Block, error)

	CreateTransaction(*Transaction) error
	FindTransaction(hash string) (*Transaction, error)
}

type SQLStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

func (s SQLStore) GetLatestNumber() (uint64, error) {
	b := &Block{}
	result := s.db.Select("number").Last(b)

	if result.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}

	if result.Error != nil {
		return 0, result.Error
	}

	return b.Number, nil
}

func (s SQLStore) GetNewBlocks(limit int) ([]*Block, error) {
	var blocks []*Block

	result := s.db.Limit(limit).Order("number desc").Find(&blocks)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get new blocks, %v", result.Error)
	}

	return blocks, nil
}

func (s SQLStore) GetUnconfirmedBlocks() ([]*Block, error) {
	var blocks []*Block
	//result := s.db.Where("confirmed = ?", false).Order("number asc").Find(&blocks)
	result := s.db.Where(&Block{Confirmed: false}).Order("number asc").Find(&blocks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get unconfirmed blocks, %v", result.Error)
	}

	return blocks, nil
}

func (s SQLStore) ConfirmBlocks(blocks []*Block) error {
	IDs := make([]uint64, len(blocks))
	for i, block := range blocks {
		IDs[i] = block.Number
	}

	result := s.db.Model(&Block{}).Where("number IN ?", IDs).Select("confirmed").Updates(Block{Confirmed: true})

	return result.Error
}

func (s SQLStore) CreateBlock(b *Block) error {
	result := s.db.Create(b)
	if result.Error != nil {
		return fmt.Errorf("failed to create MySQL record")
	}

	return nil
}

func (s SQLStore) FindBlock(number uint64) (*Block, error) {
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

func (s SQLStore) CreateTransaction(transaction *Transaction) error {
	return s.db.Create(transaction).Error
}

func (s SQLStore) FindTransaction(hash string) (*Transaction, error) {
	var tx = &Transaction{}
	result := s.db.Where("hash = ?", hash).First(tx)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transaction, %v", result.Error)
	}

	return tx, nil
}
