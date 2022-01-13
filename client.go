package portto

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

func (c *Client) GetTransactionByHash(ctx context.Context, hash string) (*Transaction, error) {
	r := &RawTransaction{}
	err := c.RPCClient.CallContext(ctx, r, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, fmt.Errorf("rpc client get transaction error, %v", err)
	}

	nonce, err := hexToUint64(r.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Nonce %q to int, %v", r.Nonce, err)
	}
	value, err := hexToUint64(r.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Value %q to int, %v", r.Value, err)
	}

	tx := &Transaction{
		Hash:  hash,
		From:  r.From,
		To:    r.To,
		Nonce: nonce,
		Data:  r.Data,
		Value: value,
	}

	return tx, err
}

func (c *Client) GetTransactionReceipt(ctx context.Context, hash string) (*types.Receipt, error) {
	return c.EthClient.TransactionReceipt(ctx, common.HexToHash(hash))
}

type RawTransaction struct {
	Hash  string `json:"tx_hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Nonce string `json:"nonce"`
	Data  string `json:"input"`
	Value string `json:"value"`
}

func hexToUint64(hex string) (uint64, error) {
	cleaned := strings.Replace(hex, "0x", "", -1)
	result, err := strconv.ParseUint(cleaned, 16, 64)

	return uint64(result), err
}
