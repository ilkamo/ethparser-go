package jsonrpc

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/ilkamo/ethparser-go/tests/mock"

	"github.com/stretchr/testify/require"
)

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

func TestClient_Call(t *testing.T) {
	ctx := context.TODO()

	t.Run("should call rpc method without errors", func(t *testing.T) {
		expectedCall := `{"jsonrpc":"2.0","id":1,"method":"test_method", "params":["param1","param2"]}`

		mockHTTPClient := &mock.HTTPClient{
			ResponseBytes: []byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`),
		}

		c, err := NewClient(endpoint, WithHTTPClient(mockHTTPClient))
		require.NoError(t, err)

		resp, err := c.Call(ctx, "test_method", []interface{}{"param1", "param2"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, []byte(`"0x1"`), []byte(resp))

		require.Equal(t, http.MethodPost, mockHTTPClient.GotRequest.Method)
		require.Equal(t, endpoint, mockHTTPClient.GotRequest.URL.String())
		require.Equal(t, "application/json", mockHTTPClient.GotRequest.Header.Get("Content-Type"))

		requestBytes, err := io.ReadAll(mockHTTPClient.GotRequest.Body)
		require.NoError(t, err)

		require.JSONEq(t, expectedCall, string(requestBytes))
	})

	t.Run("should call rpc method without error returned by rpc service", func(t *testing.T) {
		expectedCall := `{"jsonrpc":"2.0","id":1,"method":"test_method"}`

		mockHTTPClient := &mock.HTTPClient{
			ResponseBytes: []byte(`{"jsonrpc":"2.0","error": {"code": -32602, "message": "Invalid params"}}`),
		}

		c, err := NewClient(endpoint, WithHTTPClient(mockHTTPClient))
		require.NoError(t, err)

		_, err = c.Call(ctx, "test_method", nil)
		require.ErrorContains(t, err, "rpc error: Invalid params")

		requestBytes, err := io.ReadAll(mockHTTPClient.GotRequest.Body)
		require.NoError(t, err)

		require.JSONEq(t, expectedCall, string(requestBytes))
	})

	t.Run("should return error because of invalid call", func(t *testing.T) {
		mockHTTPClient := &mock.HTTPClient{
			ShouldError: true,
		}

		c, err := NewClient(endpoint, WithHTTPClient(mockHTTPClient))
		require.NoError(t, err)

		_, err = c.Call(ctx, "test", nil)
		require.ErrorContains(t, err, "could not send request: test error")
	})

	t.Run("should return error because of invalid response format", func(t *testing.T) {
		mockHTTPClient := &mock.HTTPClient{
			ResponseBytes: []byte(`invalid json`),
		}

		c, err := NewClient(endpoint, WithHTTPClient(mockHTTPClient))
		require.NoError(t, err)

		_, err = c.Call(ctx, "test", nil)
		require.ErrorContains(t, err, "could not decode response")
	})

	t.Run("should error because of invalid request", func(t *testing.T) {
		mockHTTPRequestBuilder := &mock.HTTPRequestBuilder{
			ShouldError: true,
		}

		c, err := NewClient(endpoint, WithHTTPRequestBuilder(mockHTTPRequestBuilder))
		require.NoError(t, err)

		_, err = c.Call(ctx, "test", nil)
		require.ErrorContains(t, err, "could not create request: test error")
	})
}
