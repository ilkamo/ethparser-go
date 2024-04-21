package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/ilkamo/ethparser-go/types"
	"github.com/stretchr/testify/require"
)

func TestTransactionsRepository(t *testing.T) {
	addresses := randomAddresses()
	ctx := context.TODO()

	t.Run("repo should be empty", func(t *testing.T) {
		repo := NewTransactionRepository()

		transactions, err := repo.GetTransactions(ctx, addresses[0])
		require.ErrorIs(t, err, types.ErrAddressNotFound)
		require.Empty(t, transactions)
	})

	t.Run("save transaction", func(t *testing.T) {
		repo := NewTransactionRepository()

		tx0 := types.Transaction{Hash: "0x1", From: addresses[0], To: addresses[1]}

		err := repo.SaveTransactions(ctx, []types.Transaction{tx0})
		require.NoError(t, err)

		// outbound transaction
		transactions, err := repo.GetTransactions(ctx, addresses[0])
		require.NoError(t, err)
		require.Len(t, transactions, 1)
		require.Equal(t, tx0, transactions[0])

		// inbound transaction
		transactions, err = repo.GetTransactions(ctx, addresses[1])
		require.NoError(t, err)
		require.Len(t, transactions, 1)
		require.Equal(t, tx0, transactions[0])
	})

	t.Run("save two transactions with the same from address", func(t *testing.T) {
		repo := NewTransactionRepository()

		tx0 := types.Transaction{Hash: "0x1", From: addresses[0], To: addresses[1]}
		tx1 := types.Transaction{Hash: "0x2", From: addresses[0], To: addresses[2]}

		err := repo.SaveTransactions(ctx, []types.Transaction{tx0, tx1})
		require.NoError(t, err)

		transactions, err := repo.GetTransactions(ctx, addresses[0])
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx0)
		require.Contains(t, transactions, tx1)

		transactions, err = repo.GetTransactions(ctx, addresses[1])
		require.NoError(t, err)
		require.Len(t, transactions, 1)
		require.Contains(t, transactions, tx0)
	})

	t.Run("save two transactions with the same to address", func(t *testing.T) {
		repo := NewTransactionRepository()

		tx0 := types.Transaction{Hash: "0x1", From: addresses[0], To: addresses[2]}
		tx1 := types.Transaction{Hash: "0x2", From: addresses[1], To: addresses[2]}

		err := repo.SaveTransactions(ctx, []types.Transaction{tx0, tx1})
		require.NoError(t, err)

		transactions, err := repo.GetTransactions(ctx, addresses[2])
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx0)
		require.Contains(t, transactions, tx1)

		transactions, err = repo.GetTransactions(ctx, addresses[0])
		require.NoError(t, err)
		require.Len(t, transactions, 1)
		require.Contains(t, transactions, tx0)

		transactions, err = repo.GetTransactions(ctx, addresses[1])
		require.NoError(t, err)
		require.Len(t, transactions, 1)
		require.Contains(t, transactions, tx1)
	})

	t.Run("save multiple transactions", func(t *testing.T) {
		repo := NewTransactionRepository()

		tx0 := types.Transaction{Hash: "0x1", From: addresses[0], To: addresses[1]}
		tx1 := types.Transaction{Hash: "0x2", From: addresses[1], To: addresses[2]}
		tx2 := types.Transaction{Hash: "0x3", From: addresses[2], To: addresses[0]}

		err := repo.SaveTransactions(ctx, []types.Transaction{tx0, tx1, tx2})
		require.NoError(t, err)

		transactions, err := repo.GetTransactions(ctx, addresses[0])
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx0)
		require.Contains(t, transactions, tx2)

		transactions, err = repo.GetTransactions(ctx, addresses[1])
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx0)
		require.Contains(t, transactions, tx1)

		transactions, err = repo.GetTransactions(ctx, addresses[2])
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx1)
		require.Contains(t, transactions, tx2)
	})

	t.Run("get last processed block from empty repository", func(t *testing.T) {
		repo := NewTransactionRepository()

		blockNumber, err := repo.GetLastProcessedBlock(ctx)
		require.NoError(t, err)
		require.Zero(t, blockNumber)
	})

	t.Run("get last processed block", func(t *testing.T) {
		repo := NewTransactionRepository()

		err := repo.SaveLastProcessedBlock(ctx, 100)
		require.NoError(t, err)

		blockNumber, err := repo.GetLastProcessedBlock(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(100), blockNumber)
	})

	t.Run("no case sensitive addresses", func(t *testing.T) {
		repo := NewTransactionRepository()

		tx0 := types.Transaction{Hash: "0x1", From: addresses[0], To: addresses[1]}
		tx1 := types.Transaction{Hash: "0x2", From: addresses[0], To: addresses[2]}

		err := repo.SaveTransactions(ctx, []types.Transaction{tx0, tx1})
		require.NoError(t, err)

		transactions, err := repo.GetTransactions(ctx, addresses[0])
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx0)
		require.Contains(t, transactions, tx1)

		transactions, err = repo.GetTransactions(ctx, strings.ToUpper(addresses[0]))
		require.NoError(t, err)
		require.Len(t, transactions, 2)
		require.Contains(t, transactions, tx0)
		require.Contains(t, transactions, tx1)
	})
}
