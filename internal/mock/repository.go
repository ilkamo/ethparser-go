package mock

import (
	"context"

	"github.com/ilkamo/ethparser-go/types"
)

type TransactionsRepository struct {
	GetError  error
	SaveError error
}

func (t TransactionsRepository) GetTransactions(_ context.Context, _ string) ([]types.Transaction, error) {
	if t.GetError != nil {
		return nil, t.GetError
	}

	return nil, nil
}

func (t TransactionsRepository) SaveTransactions(_ context.Context, _ []types.Transaction) error {
	if t.SaveError != nil {
		return t.SaveError
	}

	return nil
}

func (t TransactionsRepository) GetLastProcessedBlock(_ context.Context) (uint64, error) {
	if t.GetError != nil {
		return 0, t.GetError
	}

	return 0, nil
}

func (t TransactionsRepository) SaveLastProcessedBlock(_ context.Context, _ uint64) error {
	if t.SaveError != nil {
		return t.SaveError
	}

	return nil
}

type AddressesRepository struct {
	WantError error
}

func (o AddressesRepository) ObserveAddress(_ context.Context, _ string) error {
	if o.WantError != nil {
		return o.WantError
	}

	return nil
}

func (o AddressesRepository) IsAddressObserved(_ context.Context, _ string) (bool, error) {
	if o.WantError != nil {
		return false, o.WantError
	}

	return false, nil
}
