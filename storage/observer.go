package storage

import (
	"context"
	"strings"
	"sync"
)

// ObserverRepository is a repository for observing addresses.
// This is an in memory implementation however in production it should be backed by a
// fast cache storage like Redis or similar.
type ObserverRepository struct {
	observedAddresses map[string]struct{}
	sync.RWMutex
}

// NewObserverRepository creates a new ObserverRepository.
func NewObserverRepository() *ObserverRepository {
	return &ObserverRepository{
		observedAddresses: make(map[string]struct{}),
	}
}

func (o *ObserverRepository) ObserveAddress(_ context.Context, address string) error {
	o.Lock()
	defer o.Unlock()

	o.observedAddresses[strings.ToLower(address)] = struct{}{}

	return nil
}

func (o *ObserverRepository) IsAddressObserved(_ context.Context, address string) (bool, error) {
	o.RLock()
	defer o.RUnlock()

	_, ok := o.observedAddresses[strings.ToLower(address)]

	return ok, nil
}
