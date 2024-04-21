package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ilkamo/ethparser-go/types"
)

const (
	defaultTimeout = time.Second * 30
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPRequestBuilder interface {
	Build(
		ctx context.Context,
		endpoint string,
		rpcMethod string,
		params interface{},
	) (*http.Request, error)
}

type Client struct {
	endpoint           string
	httpClient         HTTPClient
	httpRequestBuilder HTTPRequestBuilder
	log                types.Logger
	// TODO: it would be nice to have some tracing collector here for better observability in production.
}

func NewClient(
	rpcEndpoint string,
	opts ...Option,
) (Client, error) {
	if rpcEndpoint == "" {
		return Client{}, fmt.Errorf("rpc endpoint is required")
	}

	c := &Client{endpoint: rpcEndpoint}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		// use default http client when not provided
		c.httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	if c.log == nil {
		// use default logger when not provided
		c.log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	if c.httpRequestBuilder == nil {
		// use default request builder when not provided
		c.httpRequestBuilder = simpleRequestBuilder{}
	}

	return *c, nil
}

// Call sends an RPC request to the server and returns the result.
func (c Client) Call(
	ctx context.Context,
	method string,
	params interface{},
) (json.RawMessage, error) {
	req, err := c.httpRequestBuilder.Build(ctx, c.endpoint, method, params)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("could not close response body", "error", err)
		}
	}()

	rpcResult, err := c.decodeResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return rpcResult, nil
}

func (c Client) decodeResponse(
	resp *http.Response,
) (json.RawMessage, error) {
	var rpcResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	if rpcResponse.Error != nil {
		c.log.Error("rpc error", "error", rpcResponse.Error)
		return nil, fmt.Errorf("rpc error: %s", rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}
