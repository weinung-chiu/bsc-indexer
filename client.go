package portto

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	RPCClient *rpc.Client
	EthClient *ethclient.Client
}

func NewClient(rpcEndpoint string) (*Client, error) {
	rpcClient, err := rpc.Dial(rpcEndpoint)
	if err != nil {
		return nil, err
	}

	ethClient := ethclient.NewClient(rpcClient)

	return &Client{rpcClient, ethClient}, nil

}

func (c *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	return c.EthClient.BlockNumber(ctx)
}

func (c *Client) GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	return c.EthClient.BlockByNumber(ctx, big.NewInt(int64(number)))
}
