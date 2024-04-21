package e2e

import (
	"log/slog"
	"os"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/ilkamo/ethparser-go/internal/ethereum"
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
	client, err := ethereum.NewClient("https://cloudflare-eth.com")
	s.Require().NoError(err)

	noNewBlockPause := time.Second * 4
	lastParsedBlock := uint64(19698124)

	s.transactionsRepository = storage.NewTransactionRepositoryWithLatestBlock(lastParsedBlock)
	s.observerRepository = storage.NewObserverRepository()

	s.parser = parser.NewParser(
		client,
		slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		noNewBlockPause,
		s.transactionsRepository,
		s.observerRepository,
	)
}
