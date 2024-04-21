package parser

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/internal/mock"
	"github.com/ilkamo/ethparser-go/internal/storage"
	"github.com/ilkamo/ethparser-go/types"
)

const endpoint = "https://test:80"

func TestParser(t *testing.T) {
	mostRecentBlockOnChain := uint64(19698125)
	noNewBlockPauseDuration := time.Millisecond * 100

	address0 := "0x115295d8C90Fe127932C6fE78daE6D5a4B975098"
	expectedTx := types.Transaction{
		Hash:  "0x005295d8C90Fe127932C6fE78daE6D5a4B975098",
		From:  address0,
		To:    "0x225295d8C90Fe127932C6fE78daE6D5a4B975098",
		Value: *big.NewInt(123),
	}
	expectedBlock := types.Block{
		Number:     mostRecentBlockOnChain,
		Hash:       "0xasd295d8C90Fe127932C6fE78daE6D5a4B975gs1",
		ParentHash: "0xasd295d8C90Fe127932C6fE78daE6D5a4B975gs2",
		Timestamp:  time.Now(),
		Transactions: []types.Transaction{
			expectedTx,
		},
	}

	ethClient := mock.EthereumClient{
		MostRecentBlock: mostRecentBlockOnChain,
		BlockByNumber:   expectedBlock,
	}

	t.Run("parser should start and process until the latest block", func(t *testing.T) {
		lastParsedBlock := uint64(19698124)
		log := &mock.Logger{}

		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(noNewBlockPauseDuration),
			WithTransactionsRepo(storage.NewTransactionRepositoryWithLatestBlock(lastParsedBlock)),
			WithEthereumClient(ethClient),
		)
		require.NoError(t, err)
		require.NotNil(t, p)
		require.Equal(t, 0, p.GetCurrentBlock())
		require.Empty(t, p.GetTransactions(address0))

		ctx, cancel := context.WithCancel(context.TODO())

		p.Subscribe(address0)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			err := p.Run(ctx)
			require.NoError(t, err)
			wg.Done()
		}()

		require.Eventually(t, func() bool {
			return p.isRunning()
		}, time.Second*2, time.Millisecond*100)

		// Should observe the subscribed transaction.
		require.Len(t, p.GetTransactions(address0), 1)
		require.Contains(t, p.GetTransactions(address0), expectedTx)
		require.Empty(t, log.GotErrors())

		require.Eventually(t, func() bool {
			return p.GetCurrentBlock() == int(mostRecentBlockOnChain)
		}, time.Second*2, time.Millisecond*100)

		require.Contains(t, log.GotInfos(), "no new blocks, sleeping to avoid spamming the node")

		cancel()
		wg.Wait()
	})

	t.Run("parser should error because of transactions repo", func(t *testing.T) {
		log := &mock.Logger{}

		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(noNewBlockPauseDuration),
			WithTransactionsRepo(mock.TransactionsRepository{GetError: errors.New("transactions error")}),
			WithEthereumClient(ethClient),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		require.Empty(t, p.GetTransactions(address0))
		require.NotEmpty(t, log.GotErrors())
		require.Equal(t, "could not get transactions", log.GotErrors()[0])
	})

	t.Run("parser could not start - should error because of context timeout", func(t *testing.T) {
		log := &mock.Logger{}

		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(noNewBlockPauseDuration),
			WithTransactionsRepo(mock.TransactionsRepository{GetError: context.DeadlineExceeded}),
			WithEthereumClient(ethClient),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		err = p.Run(context.TODO())
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("parser should log error during processing", func(t *testing.T) {
		log := &mock.Logger{}
		ethMock := mock.EthereumClient{
			MostRecentBlock: 2,
			BlockByNumber:   types.Block{},
		}

		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(noNewBlockPauseDuration),
			WithTransactionsRepo(mock.TransactionsRepository{SaveError: errors.New("save error")}),
			WithEthereumClient(ethMock),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		ctx, cancel := context.WithCancel(context.TODO())

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			err := p.Run(ctx)
			require.NoError(t, err)
			wg.Done()
		}()

		require.Eventually(t, func() bool {
			return p.isRunning() && len(log.GotErrors()) > 10
		}, time.Second*2, time.Millisecond*100)
		require.Contains(t, log.GotErrors(), "could not process blocks")

		cancel()
		wg.Wait()
	})

	t.Run("parser should log error because of addresses repo", func(t *testing.T) {
		log := &mock.Logger{}

		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(noNewBlockPauseDuration),
			WithAddressesRepo(mock.AddressesRepository{WantError: errors.New("addresses error")}),
			WithEthereumClient(ethClient),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		require.Empty(t, p.Subscribe(address0))
		require.NotEmpty(t, log.GotErrors())
		require.Equal(t, "could not observe address", log.GotErrors()[0])
	})

	t.Run("parser should error with concurrent run", func(t *testing.T) {
		log := &mock.Logger{}
		ethMock := mock.EthereumClient{
			MostRecentBlock: 2,
			BlockByNumber:   types.Block{},
		}

		p, err := NewParser(
			endpoint,
			log,
			WithNoNewBlocksPause(noNewBlockPauseDuration),
			WithEthereumClient(ethMock),
		)
		require.NoError(t, err)
		require.NotNil(t, p)

		ctx, cancel := context.WithCancel(context.TODO())

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			err := p.Run(ctx)
			require.NoError(t, err)
			wg.Done()
		}()

		require.Eventually(t, func() bool {
			return p.isRunning()
		}, time.Second*2, time.Millisecond*100)

		err = p.Run(context.TODO())
		require.ErrorAs(t, err, &types.ErrAlreadyRunning)

		cancel()
		wg.Wait()
	})
}
