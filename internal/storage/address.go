package storage

import (
	"context"
	"strings"
	"sync"
)

// AddressesRepository is a repository for addresses.
// This is an in memory implementation however in production it should be backed by a
// fast cache storage like Redis or similar.
type AddressesRepository struct {
	observedAddresses map[string]struct{}
	sync.RWMutex
}

// NewAddressesRepository creates a new AddressesRepository.
func NewAddressesRepository() *AddressesRepository {
	return &AddressesRepository{
		observedAddresses: make(map[string]struct{}),
	}
}

func (o *AddressesRepository) ObserveAddress(_ context.Context, address string) error {
	o.Lock()
	defer o.Unlock()

	o.observedAddresses[strings.ToLower(address)] = struct{}{}

	return nil
}

func (o *AddressesRepository) IsAddressObserved(_ context.Context, address string) (bool, error) {
	o.RLock()
	defer o.RUnlock()

	_, ok := o.observedAddresses[strings.ToLower(address)]

	return ok, nil
}
