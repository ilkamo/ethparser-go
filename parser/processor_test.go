package parser

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/internal/mock"
	"github.com/ilkamo/ethparser-go/internal/storage"
	"github.com/ilkamo/ethparser-go/types"
)

func TestParser_processBlock(t *testing.T) {
	log := &mock.Logger{}
	ethMock := mock.EthereumClient{
		MostRecentBlock: 2,
		BlockByNumber:   types.Block{},
	}

	p, err := NewParser(
		endpoint,
		log,
		WithNoNewBlocksPause(time.Millisecond),
		WithAddressesRepo(mock.AddressesRepository{
			WantError: errors.New("addresses error"),
		}),
		WithEthereumClient(ethMock),
	)
	require.NoError(t, err)
	require.NotNil(t, p)

	tx := types.Transaction{
		Hash:  "0x005295d8C90Fe127932C6fE78daE6D5a4B975098",
		From:  "0x995295d8C90Fe127932C6fE78daE6D5a4B975098",
		To:    "0x225295d8C90Fe127932C6fE78daE6D5a4B975098",
		Value: *big.NewInt(123),
	}

	err = p.processBlock(context.Background(), types.Block{Transactions: []types.Transaction{tx}})
	require.Error(t, err)
}

// TestParser_processBlocks tests the parallel processing of blocks.
func TestParser_processBlocks(t *testing.T) {
	log := &mock.Logger{}
	mostRecentBlockOnChain := uint64(14)
	ethMock := mock.EthereumClient{
		MostRecentBlock: mostRecentBlockOnChain,
		BlockByNumber: types.Block{
			Number: mostRecentBlockOnChain,
			Hash:   "0xasd295d8C90Fe127932C6fE78daE6D5a4B975gs1",
			Transactions: []types.Transaction{
				{
					Hash:  "0x005295d8C90Fe127932C6fE78daE6D5a4B975098",
					From:  "0x995295d8C90Fe127932C6fE78daE6D5a4B975098",
					To:    "0x225295d8C90Fe127932C6fE78daE6D5a4B975098",
					Value: *big.NewInt(123),
				},
				{
					Hash:  "0x005295d8C90Fe127932C6fE78daE6D5a4B975099",
					From:  "0x995295d8C90Fe127932C6fE78daE6D5a4B975099",
					To:    "0x225295d8C90Fe127932C6fE78daE6D5a4B975099",
					Value: *big.NewInt(123),
				},
				{
					Hash:  "0x005295d8C90Fe127932C6fE78daE6D5a4B975100",
					From:  "0x995295d8C90Fe127932C6fE78daE6D5a4B975100",
					To:    "0x225295d8C90Fe127932C6fE78daE6D5a4B975100",
					Value: *big.NewInt(123),
				},
			},
		},
	}

	maxBlocksToProcessInParallelCount := 10

	t.Run("parser should start and process until the latest block", func(t *testing.T) {
		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(time.Millisecond),
			WithTransactionsRepo(storage.NewTransactionRepositoryWithLatestBlock(0)),
			WithEthereumClient(ethMock),
			WithMaxBlocksToProcessInParallel(maxBlocksToProcessInParallelCount),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		ctx := context.TODO()
		err = p.processBlocks(ctx)
		require.NoError(t, err)

		// After the first iteration, the last processed block should be equal to maxBlocksToProcessInParallelCount.
		require.Equal(t, maxBlocksToProcessInParallelCount, p.GetCurrentBlock())

		err = p.processBlocks(ctx)
		require.NoError(t, err)

		// After the second iteration, the last processed block should be equal to mostRecentBlockOnChain as there are no new blocks.
		require.Equal(t, int(mostRecentBlockOnChain), p.GetCurrentBlock())
	})

	t.Run("parser should process blocks and return an error", func(t *testing.T) {
		p, err := NewParser(
			endpoint,
			log,
			WithTransactionsRepo(storage.NewTransactionRepositoryWithLatestBlock(0)),
			WithEthereumClient(ethMock),
			WithMaxBlocksToProcessInParallel(maxBlocksToProcessInParallelCount),
			WithAddressesRepo(mock.AddressesRepository{
				WantError: errors.New("addresses error"),
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		ctx := context.TODO()
		err = p.processBlocks(ctx)
		require.Error(t, err)
	})
}
