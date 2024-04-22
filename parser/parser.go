package parser

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ilkamo/ethparser-go/internal/ethereum"
	"github.com/ilkamo/ethparser-go/internal/storage"
	"github.com/ilkamo/ethparser-go/types"
)

const (
	defaultBlocksProcessTimeout       = 5 * time.Second
	defaultNoNewBlocksPause           = 10 * time.Second // eth new block appears every ~12 seconds
	defaultMaxNumberOfBlocksToProcess = 10
)

type Parser struct {
	blocksProcessTimeout                 time.Duration
	ethClient                            EthereumClient
	lastProcessedBlock                   uint64
	logger                               types.Logger
	noNewBlocksPause                     time.Duration
	transactionsRepo                     TransactionsRepository
	addressesRepository                  AddressesRepository
	running                              bool
	batchesWorker                        chan struct{}
	maxNumberOfBlocksToProcessInParallel int
	processingErrs                       []error
	mutex                                sync.RWMutex
}

func NewParser(
	rpcEndpoint string,
	logger types.Logger,
	opts ...Option,
) (*Parser, error) {
	p := &Parser{
		blocksProcessTimeout:                 defaultBlocksProcessTimeout,
		logger:                               logger,
		noNewBlocksPause:                     defaultNoNewBlocksPause,
		transactionsRepo:                     storage.NewTransactionRepository(),
		addressesRepository:                  storage.NewAddressesRepository(),
		batchesWorker:                        make(chan struct{}, 1),
		maxNumberOfBlocksToProcessInParallel: defaultMaxNumberOfBlocksToProcess,
		processingErrs:                       make([]error, 0),
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.ethClient == nil {
		ethClient, err := ethereum.NewClient(rpcEndpoint)
		if err != nil {
			return nil, fmt.Errorf("could not create Ethereum client: %w", err)
		}

		p.ethClient = ethClient
	}

	p.batchesWorker <- struct{}{}

	return p, nil
}

func (p *Parser) GetCurrentBlock() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// 256-bit unsigned integer is the maximum value for an Ethereum block number.
	// It isn't safe to cast it to int directly, but I will do it here to implement the `Parser` interface.
	return int(p.lastProcessedBlock)
}

// Subscribe adds an address to the list of addresses to watch for transactions.
// It is not clear to me what the returned bool means. I assumed it returns true if the address was successfully
// added, false if it not. I would add an error to the return value to provide more information about the failure.
// Additionally, I would add a context to the method signature.
func (p *Parser) Subscribe(address string) bool {
	err := p.addressesRepository.ObserveAddress(context.Background(), address)
	if err != nil {
		p.logger.Error("could not observe address", err)
		return false
	}

	p.logger.Info("started observing address", "address", address)

	return true
}

// GetTransactions returns a list of transactions for an address.
// I cannot change the signature of the method as it is defined in the `Parser` interface.
// However, IMO it would be better to return an error if something goes wrong.
// In addition, I would add a context to the method signature to handle timeouts and cancellations.
func (p *Parser) GetTransactions(address string) []types.Transaction {
	transactions, err := p.transactionsRepo.GetTransactions(context.Background(), address)
	if err != nil {
		if errors.Is(err, types.ErrAddressNotFound) {
			return nil
		}

		p.logger.Error("could not get transactions", err)
		return nil
	}

	return transactions
}

// getNumberOfBlocksToProcess calculates the number of blocks that the parser should process in the next iteration.
func (p *Parser) getNumberOfBlocksToProcess(ctx context.Context) (int, uint64, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	lastBlockNumber, err := p.ethClient.GetMostRecentBlockNumber(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("could not get most recent block: %w", err)
	}

	blocksToProcessCount := int(lastBlockNumber - uint64(p.GetCurrentBlock()))

	if blocksToProcessCount > p.maxNumberOfBlocksToProcessInParallel {
		blocksToProcessCount = p.maxNumberOfBlocksToProcessInParallel
	}

	lastBlockOfTheSequence := p.GetCurrentBlock() + blocksToProcessCount

	p.logger.Info("calculated blocks to process",
		"blocks", blocksToProcessCount, "lastBlockOfTheSequence", lastBlockOfTheSequence)

	return blocksToProcessCount, uint64(lastBlockOfTheSequence), nil
}

// Run starts the parser and listens for new blocks.
// This method is not specified in the task `Parser` interface, but I added it to start the
// parser explicitly (not in the constructor).
// I also added a context to handle timeouts and cancellations.
// When called, it starts processing blocks in a loop until the context is canceled.
// The starting block is the last processed block from the repository so that the parser
// can continue from where it left off after a restart.
func (p *Parser) Run(ctx context.Context) error {
	if p.isRunning() {
		return errors.New("parser is already running")
	}

	p.setIsRunning(true)
	defer func() {
		p.setIsRunning(false)
	}()

	latestProcessed, err := p.transactionsRepo.GetLastProcessedBlock(ctx)
	if err != nil {
		return err
	}

	p.setLastProcessedBlock(latestProcessed)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("stopping parser")
			return nil
		case <-p.batchesWorker:
			if err := p.processBlocks(ctx); err != nil {
				p.logger.Error("could not process blocks", "error", err)
			}

			p.batchesWorker <- struct{}{}
		}
	}
}

func (p *Parser) isRunning() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.running
}

func (p *Parser) setIsRunning(running bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.running = running
}

// setLastProcessedBlock sets the last processed block number.
func (p *Parser) setLastProcessedBlock(blockNumber uint64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.lastProcessedBlock = blockNumber
}
