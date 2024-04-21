package ethereum

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilkamo/ethparser-go/hexutils"
	"github.com/ilkamo/ethparser-go/jsonrpc"
	"github.com/ilkamo/ethparser-go/types"
)

type RPCClient interface {
	Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error)
}

type Option func(c *Client)

type Client struct {
	rpcClient RPCClient
}

func NewClient(endpoint string, opts ...Option) (Client, error) {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.rpcClient == nil {
		rpcClient, err := jsonrpc.NewClient(endpoint)
		if err != nil {
			return Client{}, fmt.Errorf("could not create rpc client: %w", err)
		}

		c.rpcClient = rpcClient
	}

	return *c, nil
}

// GetMostRecentBlock returns the number of the most recent block.
func (c Client) GetMostRecentBlock(ctx context.Context) (uint64, error) {
	resp, err := c.rpcClient.Call(ctx, "eth_blockNumber", nil)
	if err != nil {
		return 0, fmt.Errorf("could not call rpc method: %w", err)
	}

	var blockNumber string
	if err := json.Unmarshal(resp, &blockNumber); err != nil {
		return 0, fmt.Errorf("could not unmarshal block number: %w", err)
	}

	return hexutils.DecodeEthNumberToUint(blockNumber)
}

// GetBlockByNumber returns a block by its number.
func (c Client) GetBlockByNumber(ctx context.Context, blockNumber uint64) (types.Block, error) {
	resp, err := c.rpcClient.Call(
		ctx,
		"eth_getBlockByNumber",
		[]interface{}{hexutils.EncodeToEthNumber(blockNumber), true},
	)
	if err != nil {
		return types.Block{}, fmt.Errorf("could not call rpc method: %w", err)
	}

	var b block
	if err := json.Unmarshal(resp, &b); err != nil {
		return types.Block{}, fmt.Errorf("could not unmarshal block: %w", err)
	}

	return b.ToBlock()
}
