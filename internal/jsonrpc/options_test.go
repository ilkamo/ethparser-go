package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/internal/mock"
)

const endpoint = "http://localhost:8545"

func TestWithHTTPClient(t *testing.T) {
	t.Run("with nil http client - should use default", func(t *testing.T) {
		c, err := NewClient(endpoint, WithHTTPClient(nil))
		require.NoError(t, err)
		require.NotNil(t, c.httpClient)
	})

	t.Run("with defined http client", func(t *testing.T) {
		mockedClient := mock.HTTPClient{}

		c, err := NewClient(endpoint, WithHTTPClient(&mockedClient))
		require.NoError(t, err)
		require.Equal(t, &mockedClient, c.httpClient)
	})
}

func TestWithLogger(t *testing.T) {
	t.Run("with nil logger - should use default", func(t *testing.T) {
		c, err := NewClient(endpoint, WithLogger(nil))
		require.NoError(t, err)
		require.NotNil(t, c.log)
	})

	t.Run("with defined logger", func(t *testing.T) {
		mockedLogger := mock.Logger{}

		c, err := NewClient(endpoint, WithLogger(&mockedLogger))
		require.NoError(t, err)
		require.Equal(t, &mockedLogger, c.log)
	})
}

func TestWithHTTPRequestBuilder(t *testing.T) {
	t.Run("with nil http request builder - should use default", func(t *testing.T) {
		c, err := NewClient(endpoint, WithHTTPRequestBuilder(nil))
		require.NoError(t, err)
		require.NotNil(t, c.httpRequestBuilder)
	})

	t.Run("with defined http request builder", func(t *testing.T) {
		mockedBuilder := mock.HTTPRequestBuilder{}

		c, err := NewClient(endpoint, WithHTTPRequestBuilder(&mockedBuilder))
		require.NoError(t, err)
		require.Equal(t, &mockedBuilder, c.httpRequestBuilder)
	})
}
