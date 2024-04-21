package parser

import (
	"time"

	"github.com/ilkamo/ethparser-go/types"
)

type Option func(p *Parser)

func WithBlockProcessTimeout(timeout time.Duration) Option {
	return func(p *Parser) {
		p.blockProcessTimeout = timeout
	}
}

func WithLogger(logger types.Logger) Option {
	return func(p *Parser) {
		p.logger = logger
	}
}
