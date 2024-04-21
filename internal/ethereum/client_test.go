package ethereum

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/internal/mock"
	"github.com/ilkamo/ethparser-go/internal/testdata"
)

const endpoint = "http://localhost:1212"

func TestNewClient(t *testing.T) {
	t.Run("should create new client", func(t *testing.T) {
		c, err := NewClient(endpoint)
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("should return error because of empty endpoint", func(t *testing.T) {
		_, err := NewClient("")
		require.Error(t, err)
	})
}

func TestClient_GetMostRecentBlock(t *testing.T) {
	ctx := context.TODO()

	t.Run("should error because of bad response", func(t *testing.T) {
		mockRPCClient := &mock.RPCClient{ShouldError: false}

		c, err := NewClient(endpoint, WithRPCClient(mockRPCClient))
		require.NoError(t, err)

		_, err = c.GetMostRecentBlock(ctx)
		require.ErrorContains(t, err, "could not unmarshal block number: unexpected end of JSON input")
	})

	t.Run("should return most recent block number", func(t *testing.T) {
		mockRPCClient := &mock.RPCClient{Response: []byte(`"0x1"`)}

		c, err := NewClient(endpoint, WithRPCClient(mockRPCClient))
		require.NoError(t, err)

		blockNumber, err := c.GetMostRecentBlock(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), blockNumber)
	})

	t.Run("should error because of rpc error", func(t *testing.T) {
		mockRPCClient := &mock.RPCClient{ShouldError: true}

		c, err := NewClient(endpoint, WithRPCClient(mockRPCClient))
		require.NoError(t, err)

		_, err = c.GetMostRecentBlock(ctx)
		require.ErrorContains(t, err, "could not call rpc method: test error")
	})
}

func TestClient_GetBlockByNumber(t *testing.T) {
	ctx := context.TODO()

	t.Run("should error because of bad response", func(t *testing.T) {
		mockRPCClient := &mock.RPCClient{ShouldError: false}

		c, err := NewClient(endpoint, WithRPCClient(mockRPCClient))
		require.NoError(t, err)

		_, err = c.GetBlockByNumber(ctx, 1)
		require.ErrorContains(t, err, "could not unmarshal block: unexpected end of JSON input")
	})

	t.Run("should return valid block", func(t *testing.T) {
		mockRPCClient := &mock.RPCClient{Response: testdata.BlockJSON}

		c, err := NewClient(endpoint, WithRPCClient(mockRPCClient))
		require.NoError(t, err)

		expectedBlockNumber := uint64(19697111)
		gotBlock, err := c.GetBlockByNumber(ctx, expectedBlockNumber)
		require.NoError(t, err)
		require.Equal(t, expectedBlock, gotBlock)
	})

	t.Run("should error because of rpc error", func(t *testing.T) {
		mockRPCClient := &mock.RPCClient{ShouldError: true}

		c, err := NewClient(endpoint, WithRPCClient(mockRPCClient))
		require.NoError(t, err)

		_, err = c.GetBlockByNumber(ctx, 1)
		require.ErrorContains(t, err, "could not call rpc method: test error")
	})
}
