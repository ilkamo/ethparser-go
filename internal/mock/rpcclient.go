package mock

import (
	"context"
	"encoding/json"
	"errors"
)

type RPCClient struct {
	ShouldError bool
	Response    json.RawMessage
}

func (r RPCClient) Call(_ context.Context, _ string, _ interface{}) (json.RawMessage, error) {
	if r.ShouldError {
		return nil, errors.New("test error")
	}

	return r.Response, nil
}
