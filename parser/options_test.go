package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/internal/mock"
)

func TestWithBlockProcessTimeout(t *testing.T) {
	t.Run("set block process timeout opt", func(t *testing.T) {
		p, err := NewParser(endpoint, nil, WithBlockProcessTimeout(2000))
		require.NoError(t, err)
		require.Equal(t, time.Duration(2000), p.blockProcessTimeout)
	})
}

func TestWithLogger(t *testing.T) {
	t.Run("set logger opt", func(t *testing.T) {
		log := &mock.Logger{}

		p, err := NewParser(endpoint, nil, WithLogger(log))
		require.NoError(t, err)
		require.NotNil(t, p.logger)
		require.Equal(t, log, p.logger)
	})
}

func TestWithEthereumClient(t *testing.T) {
	t.Run("set ethereum client opt", func(t *testing.T) {
		ethClient := mock.EthereumClient{}

		p, err := NewParser("", nil, WithEthereumClient(ethClient))
		require.NoError(t, err)
		require.NotNil(t, p.ethClient)
		require.Equal(t, ethClient, p.ethClient)
	})
}

func TestWithTransactionsRepo(t *testing.T) {
	t.Run("set transactions repo opt", func(t *testing.T) {
		repo := mock.TransactionsRepository{}

		p, err := NewParser(endpoint, nil, WithTransactionsRepo(repo))
		require.NoError(t, err)
		require.NotNil(t, p.transactionsRepo)
		require.Equal(t, repo, p.transactionsRepo)
	})
}

func TestWithObserverRepo(t *testing.T) {
	t.Run("set observer repo opt", func(t *testing.T) {
		repo := mock.ObserverRepository{}

		p, err := NewParser(endpoint, nil, WithObserverRepo(repo))
		require.NoError(t, err)
		require.NotNil(t, p.observerRepo)
		require.Equal(t, repo, p.observerRepo)
	})
}

func TestWithNoNewBlocksPause(t *testing.T) {
	t.Run("set no new blocks pause opt", func(t *testing.T) {
		p, err := NewParser(endpoint, nil, WithNoNewBlocksPause(2000))
		require.NoError(t, err)
		require.Equal(t, time.Duration(2000), p.noNewBlocksPause)
	})
}
