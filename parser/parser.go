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
	defaultBlockProcessTimeout = 5 * time.Second
	defaultNoNewBlocksPause    = 10 * time.Second // eth new block appears every ~12 seconds
)

type Parser struct {
	blockProcessTimeout time.Duration
	ethClient           EthereumClient
	lastProcessedBlock  uint64
	logger              types.Logger
	noNewBlocksPause    time.Duration
	transactionsRepo    TransactionsRepository
	addressesRepository AddressesRepository
	running             bool
	singleWorkerChannel chan struct{}
	sync.RWMutex
}

func NewParser(
	rpcEndpoint string,
	logger types.Logger,
	opts ...Option,
) (*Parser, error) {
	p := &Parser{
		blockProcessTimeout: defaultBlockProcessTimeout,
		logger:              logger,
		noNewBlocksPause:    defaultNoNewBlocksPause,
		transactionsRepo:    storage.NewTransactionRepository(),
		addressesRepository: storage.NewAddressesRepository(),
		singleWorkerChannel: make(chan struct{}, 1),
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

	p.singleWorkerChannel <- struct{}{}

	return p, nil
}

func (p *Parser) GetCurrentBlock() int {
	p.RLock()
	defer p.RUnlock()

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
		case <-p.singleWorkerChannel:
			ctx, cancel := context.WithTimeout(ctx, defaultBlockProcessTimeout)

			p.logger.Info("fetching and parsing block")

			if err := p.fetchAndParseBlock(ctx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					p.logger.Error("could not fetch and parse because the context expired")
				} else {
					p.logger.Error("could not fetch and parse block", "error", err)
				}
			}
			if err == nil {
				p.logger.Info("block fetched and parsed", "block", p.lastProcessedBlock)
			}

			p.singleWorkerChannel <- struct{}{}
			cancel()
		}
	}
}

func (p *Parser) isRunning() bool {
	p.RLock()
	defer p.RUnlock()

	return p.running
}

func (p *Parser) setIsRunning(running bool) {
	p.Lock()
	defer p.Unlock()

	p.running = running
}

func (p *Parser) fetchAndParseBlock(ctx context.Context) error {
	lastBlockNumber, err := p.ethClient.GetMostRecentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("could not get most recent block: %w", err)
	}

	// Check if there are new blocks to process. If not, sleep for a while to avoid spamming the node.
	if !p.shouldProcessBlock(lastBlockNumber) {
		p.logger.Info("no new blocks, sleeping to avoid spamming the node")
		time.Sleep(p.noNewBlocksPause)
		return nil
	}

	// Get the next block to process in the sequence.
	block, err := p.ethClient.GetBlockByNumber(ctx, p.lastProcessedBlock+1)
	if err != nil {
		return fmt.Errorf("could not get block by number: %w", err)
	}

	// Process the block.
	if err = p.processBlock(ctx, block); err != nil {
		return fmt.Errorf("could not process block: %w", err)
	}

	// Save the last processed block.
	if err = p.transactionsRepo.SaveLastProcessedBlock(ctx, block.Number); err != nil {
		return fmt.Errorf("could not save last processed block: %w", err)
	}

	// Move to sequence to the next block.
	p.setLastProcessedBlock(block.Number)

	return nil
}

// shouldProcessBlock checks if there are new blocks to process.
func (p *Parser) shouldProcessBlock(lastBlockNumber uint64) bool {
	return p.lastProcessedBlock < lastBlockNumber
}

// setLastProcessedBlock sets the last processed block number.
func (p *Parser) setLastProcessedBlock(blockNumber uint64) {
	p.Lock()
	defer p.Unlock()

	p.lastProcessedBlock = blockNumber
}
