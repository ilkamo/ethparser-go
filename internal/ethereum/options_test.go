package ethereum

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/internal/mock"
)

func TestWithRPCClient(t *testing.T) {
	mockedRPCClient := &mock.RPCClient{}

	t.Run("should set rpc client", func(t *testing.T) {
		c, err := NewClient("http://localhost:1212")
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotEqual(t, mockedRPCClient, c.rpcClient)

		c, err = NewClient("http://localhost:1212", WithRPCClient(mockedRPCClient))
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, mockedRPCClient, c.rpcClient)
	})
}
