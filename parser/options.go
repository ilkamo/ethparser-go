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

func WithEthereumClient(client EthereumClient) Option {
	return func(p *Parser) {
		p.ethClient = client
	}
}

func WithTransactionsRepo(repo TransactionsRepository) Option {
	return func(p *Parser) {
		p.transactionsRepo = repo
	}
}

func WithAddressesRepo(repo AddressesRepository) Option {
	return func(p *Parser) {
		p.addressesRepository = repo
	}
}

func WithNoNewBlocksPause(duration time.Duration) Option {
	return func(p *Parser) {
		p.noNewBlocksPause = duration
	}
}
