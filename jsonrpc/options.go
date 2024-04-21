package jsonrpc

import "github.com/ilkamo/ethparser-go/types"

type Option func(c *Client)

func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithLogger(logger types.Logger) Option {
	return func(c *Client) {
		c.log = logger
	}
}

func WithHTTPRequestBuilder(httpRequestBuilder HTTPRequestBuilder) Option {
	return func(c *Client) {
		c.httpRequestBuilder = httpRequestBuilder
	}
}
