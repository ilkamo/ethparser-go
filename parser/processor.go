package parser

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ilkamo/ethparser-go/types"
)

// processBlocks processes the blocks in batches. It gets the number of blocks to process, then processes them in parallel.
// It waits for processed blocks and updates the `last processed block indicator` only if all the blocks were processed successfully.
// If an error occurs during processing, the batch will be retried again in the next iteration. This assumes that
// parser repositories are idempotent. It is an all or nothing approach that works well if
// the rpc client is reliable and the number of blocks to process is well tuned. This could be improved with a more sophisticated
// approach that would allow for partial processing of the batch by tracking the processed blocks.
func (p *Parser) processBlocks(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultBlocksProcessTimeout)
	defer cancel()

	blocksToProcessCount, lastBlockNumberOfTheSequence, err := p.getNumberOfBlocksToProcess(ctx)
	if err != nil {
		return fmt.Errorf("could not get number of blocks to process: %w", err)
	}

	if blocksToProcessCount == 0 {
		p.logger.Info("no new blocks, sleeping to avoid spamming the node")
		time.Sleep(p.noNewBlocksPause)
		return nil
	}

	wg := sync.WaitGroup{}
	for i := 0; i < blocksToProcessCount; i++ {
		wg.Add(1)

		go func(blockNumber int) {
			defer wg.Done()

			block, err := p.ethClient.GetBlockByNumber(ctx, uint64(blockNumber))
			if err != nil {
				p.logger.Error("could not get block by number", "block", blockNumber, "error", err)
				p.setProcessingError(err)
				return
			}

			if err := p.processBlock(ctx, block); err != nil {
				p.logger.Error("could not process block", "block", block.Number, "error", err)
				p.setProcessingError(err)
			}
		}(p.GetCurrentBlock() + i + 1)
	}
	wg.Wait()

	if len(p.getProcessingErrors()) > 0 {
		return fmt.Errorf("errors occurred during block processing: %v", p.getProcessingErrors())
	}

	// Clear the processing errors for the next iteration.
	p.clearProcessingErrors()

	// Save the last processed block of the sequence.
	if err = p.transactionsRepo.SaveLastProcessedBlock(ctx, lastBlockNumberOfTheSequence); err != nil {
		return fmt.Errorf("could not save last processed block of the sequence: %w", err)
	}

	// Move the sequence forward.
	p.setLastProcessedBlock(lastBlockNumberOfTheSequence)

	return nil
}

// processBlock processes the block by filtering out observed transactions and saving them to the repository.
func (p *Parser) processBlock(ctx context.Context, block types.Block) error {
	p.logger.Info("processing block", "block", block.Number, "transactions", len(block.Transactions))

	observedTx, err := p.processAndFilterObservedTransactions(ctx, block.Transactions)
	if err != nil {
		return fmt.Errorf("could not filter observed transactions: %w", err)
	}

	p.logger.Info("observed transactions", "transactions", len(observedTx))

	if err := p.transactionsRepo.SaveTransactions(ctx, observedTx); err != nil {
		return fmt.Errorf("could not save transactions: %w", err)
	}

	return nil
}

// processAndFilterObservedTransactions filters out transactions that involve observed addresses.
func (p *Parser) processAndFilterObservedTransactions(
	ctx context.Context,
	transactions []types.Transaction,
) ([]types.Transaction, error) {
	var filtered []types.Transaction

	for _, tx := range transactions {
		okFrom, err := p.addressesRepository.IsAddressObserved(ctx, tx.From)
		if err != nil {
			return nil, fmt.Errorf("could not check if address `from` is observed: %w", err)
		}

		okTo, err := p.addressesRepository.IsAddressObserved(ctx, tx.To)
		if err != nil {
			return nil, fmt.Errorf("could not check if address `to` is observed: %w", err)
		}

		if okFrom || okTo {
			filtered = append(filtered, tx)
		}
	}

	return filtered, nil
}

func (p *Parser) setProcessingError(err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.processingErrs = append(p.processingErrs, err)
}

func (p *Parser) getProcessingErrors() []error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.processingErrs
}

func (p *Parser) clearProcessingErrors() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.processingErrs = nil
}
