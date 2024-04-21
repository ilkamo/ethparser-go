package parser

import (
	"context"
	"fmt"

	"github.com/ilkamo/ethparser-go/types"
)

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
