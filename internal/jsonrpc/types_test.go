package jsonrpc

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_simpleRequestBuilder_Build(t *testing.T) {
	ctx := context.TODO()

	t.Run("builder without endpoint", func(t *testing.T) {
		builder := simpleRequestBuilder{}
		_, err := builder.Build(ctx, "", "", nil)
		require.ErrorContains(t, err, "endpoint is required")
	})

	t.Run("builder without rpc method", func(t *testing.T) {
		builder := simpleRequestBuilder{}
		_, err := builder.Build(ctx, "https://test.com", "", nil)
		require.ErrorContains(t, err, "method is required")
	})

	t.Run("builder with endpoint and rpc method - with params", func(t *testing.T) {
		expectedBody := `{"jsonrpc":"2.0","id":1,"method":"test_method", "params":["param1","param2"]}`

		builder := simpleRequestBuilder{}
		req, err := builder.Build(ctx, "https://test.com", "test_method", []interface{}{"param1", "param2"})
		require.NoError(t, err)

		require.Equal(t, "https://test.com", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		requestBytes, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		require.JSONEq(t, expectedBody, string(requestBytes))
	})

	t.Run("builder with invalid context", func(t *testing.T) {
		builder := simpleRequestBuilder{}

		//nolint:staticcheck
		_, err := builder.Build(nil, "https://test.com", "test_method", nil)
		require.ErrorContains(t, err, "could not create request: net/http: nil Context")
	})
}
