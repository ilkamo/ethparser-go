package e2e

import (
	"log/slog"
	"os"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/ilkamo/ethparser-go/internal/storage"
	"github.com/ilkamo/ethparser-go/parser"
)

type ParserTestSuite struct {
	suite.Suite
	parser                 *parser.Parser
	transactionsRepository parser.TransactionsRepository
	observerRepository     parser.ObserverRepository
}

func (s *ParserTestSuite) SetupSuite() {
	noNewBlockPause := time.Second * 4
	lastParsedBlock := uint64(19698124)

	s.transactionsRepository = storage.NewTransactionRepositoryWithLatestBlock(lastParsedBlock)
	s.observerRepository = storage.NewObserverRepository()

	p, err := parser.NewParser(
		"https://cloudflare-eth.com",
		slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		parser.WithNoNewBlocksPause(noNewBlockPause),
		parser.WithTransactionsRepo(s.transactionsRepository),
		parser.WithObserverRepo(s.observerRepository),
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.parser = p
}
