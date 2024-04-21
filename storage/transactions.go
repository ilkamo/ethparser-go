package storage

import (
	"context"
	"strings"
	"sync"

	"github.com/ilkamo/ethparser-go/types"
)

type TransactionsRepository struct {
	latestBlock uint64
	// A simple in-memory storage for transactions -> map[address]map[txHash]tx
	// I am using a map instead of a slice to avoid duplicates in the storage in case of reprocessing
	// because of a failure.
	transactionsPerAddress map[string]map[string]types.Transaction
	sync.RWMutex
}

func NewTransactionRepository() *TransactionsRepository {
	return newTransactionRepository(0)
}

func NewTransactionRepositoryWithLatestBlock(latestBlock uint64) *TransactionsRepository {
	return newTransactionRepository(latestBlock)
}

func newTransactionRepository(
	latestBlock uint64,
) *TransactionsRepository {
	return &TransactionsRepository{
		latestBlock:            latestBlock,
		transactionsPerAddress: make(map[string]map[string]types.Transaction),
	}
}

func (t *TransactionsRepository) GetTransactions(_ context.Context, address string) ([]types.Transaction, error) {
	t.RLock()
	defer t.RUnlock()

	transactions, ok := t.transactionsPerAddress[strings.ToLower(address)]
	if !ok {
		return nil, types.ErrAddressNotFound
	}

	result := make([]types.Transaction, 0, len(transactions))
	for _, tx := range transactions {
		result = append(result, tx)
	}

	return result, nil
}

func (t *TransactionsRepository) SaveTransactions(_ context.Context, transactions []types.Transaction) error {
	t.Lock()
	defer t.Unlock()

	for _, tx := range transactions {
		txFrom := strings.ToLower(tx.From)
		txTo := strings.ToLower(tx.To)
		txHash := strings.ToLower(tx.Hash)

		transactionsFrom, ok := t.transactionsPerAddress[txFrom]
		if !ok {
			transactionsFrom = make(map[string]types.Transaction)
			t.transactionsPerAddress[txFrom] = transactionsFrom
		}

		transactionsFrom[txHash] = tx

		transactionsTo, ok := t.transactionsPerAddress[txTo]
		if !ok {
			transactionsTo = make(map[string]types.Transaction)
			t.transactionsPerAddress[txTo] = transactionsTo
		}

		transactionsTo[txHash] = tx
	}

	return nil
}

func (t *TransactionsRepository) SaveLastProcessedBlock(_ context.Context, blockNumber uint64) error {
	t.Lock()
	defer t.Unlock()

	t.latestBlock = blockNumber

	return nil
}

func (t *TransactionsRepository) GetLastProcessedBlock(_ context.Context) (uint64, error) {
	t.RLock()
	defer t.RUnlock()

	return t.latestBlock, nil
}
