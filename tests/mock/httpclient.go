package mock

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

type HTTPClient struct {
	ShouldError   bool
	ResponseBytes []byte
	GotRequest    *http.Request
}

func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if h.ShouldError {
		return nil, errors.New("test error")
	}

	h.GotRequest = req

	body := io.NopCloser(bytes.NewReader(h.ResponseBytes))

	return &http.Response{
		Body:       body,
		StatusCode: 200,
	}, nil
}

type HTTPRequestBuilder struct {
	ShouldError bool
}

func (h *HTTPRequestBuilder) Build(
	_ context.Context,
	_ string,
	_ string,
	_ interface{},
) (*http.Request, error) {
	if h.ShouldError {
		return nil, errors.New("test error")
	}

	return &http.Request{}, nil
}
