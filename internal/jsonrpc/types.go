package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int64       `json:"id,omitempty"`
}

type Response struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
	ID      int64           `json:"id,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func newRequestBody(method string, params interface{}) (io.Reader, error) {
	if method == "" {
		return nil, fmt.Errorf("method is required")
	}

	buff := new(bytes.Buffer)

	if err := json.NewEncoder(buff).Encode(Request{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1, // TODO: generate unique ID
	}); err != nil {
		return nil, err
	}

	return buff, nil
}

type simpleRequestBuilder struct{}

func (s simpleRequestBuilder) Build(
	ctx context.Context,
	endpoint string,
	rpcMethod string,
	params interface{},
) (*http.Request, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	reqBody, err := newRequestBody(rpcMethod, params)
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
