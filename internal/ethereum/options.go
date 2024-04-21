package ethereum

// WithRPCClient sets the RPC client for the Ethereum client.
func WithRPCClient(rpcClient RPCClient) Option {
	return func(c *Client) {
		c.rpcClient = rpcClient
	}
}
