package portto

import "github.com/ethereum/go-ethereum/core/types"

type Repository interface {
	StoreBlock(block *types.Block) error
}

type InMemoryStore struct{}

func (InMemoryStore) StoreBlock(block *types.Block) error {
	//TODO implement me
	panic("implement me")
}
