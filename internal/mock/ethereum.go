package mock

import (
	"context"

	"github.com/ilkamo/ethparser-go/types"
)

type EthereumClient struct {
	MostRecentBlock uint64
	BlockByNumber   types.Block
	WithError       error
}

func (e EthereumClient) GetMostRecentBlock(_ context.Context) (uint64, error) {
	if e.WithError != nil {
		return 0, e.WithError
	}

	return e.MostRecentBlock, nil
}

func (e EthereumClient) GetBlockByNumber(_ context.Context, _ uint64) (types.Block, error) {
	if e.WithError != nil {
		return types.Block{}, e.WithError
	}

	return e.BlockByNumber, nil
}
