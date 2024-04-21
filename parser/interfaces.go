package parser

import (
	"context"

	"github.com/ilkamo/ethparser-go/types"
)

type TransactionsRepository interface {
	// GetTransactions returns a list of transactions for an address.
	GetTransactions(ctx context.Context, address string) ([]types.Transaction, error)

	// SaveTransactions saves transactions to the repository.
	SaveTransactions(ctx context.Context, t []types.Transaction) error

	// GetLastProcessedBlock returns the last processed block number.
	GetLastProcessedBlock(ctx context.Context) (uint64, error)

	// SaveLastProcessedBlock saves the last processed block number.
	SaveLastProcessedBlock(ctx context.Context, blockNumber uint64) error
}

type ObserverRepository interface {
	// ObserveAddress adds an address to the list of observed addresses.
	ObserveAddress(ctx context.Context, address string) error

	// IsAddressObserved checks if an address is observed.
	IsAddressObserved(ctx context.Context, address string) (bool, error)
}

type EthereumClient interface {
	// GetMostRecentBlock returns the most recent block number.
	GetMostRecentBlock(ctx context.Context) (uint64, error)

	// GetBlockByNumber returns a block by its number.
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (types.Block, error)
}
