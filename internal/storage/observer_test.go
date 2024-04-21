package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObserverRepository(t *testing.T) {
	addresses := randomAddresses()
	ctx := context.TODO()

	t.Run("repo should be empty", func(t *testing.T) {
		repo := NewObserverRepository()

		isObserved, err := repo.IsAddressObserved(ctx, addresses[0])
		require.NoError(t, err)
		require.False(t, isObserved, "should not be observed")
	})

	t.Run("observe address", func(t *testing.T) {
		repo := NewObserverRepository()

		err := repo.ObserveAddress(ctx, addresses[0])
		require.NoError(t, err)

		isObserved, err := repo.IsAddressObserved(ctx, addresses[0])
		require.NoError(t, err)
		require.True(t, isObserved, "should be observed")

		isObserved, err = repo.IsAddressObserved(ctx, addresses[1])
		require.NoError(t, err)
		require.False(t, isObserved, "should not be observed")
	})

	t.Run("observe more than one address", func(t *testing.T) {
		repo := NewObserverRepository()

		err := repo.ObserveAddress(ctx, addresses[0])
		require.NoError(t, err)

		err = repo.ObserveAddress(ctx, addresses[1])
		require.NoError(t, err)

		isObserved, err := repo.IsAddressObserved(ctx, addresses[0])
		require.NoError(t, err)
		require.True(t, isObserved, "should be observed")

		isObserved, err = repo.IsAddressObserved(ctx, addresses[1])
		require.NoError(t, err)
		require.True(t, isObserved, "should be observed")

		isObserved, err = repo.IsAddressObserved(ctx, addresses[2])
		require.NoError(t, err)
		require.False(t, isObserved, "should not be observed")
	})

	t.Run("repo should be idempotent", func(t *testing.T) {
		repo := NewObserverRepository()

		err := repo.ObserveAddress(ctx, addresses[1])
		require.NoError(t, err)

		err = repo.ObserveAddress(ctx, addresses[1])
		require.NoError(t, err)

		isObserved, err := repo.IsAddressObserved(ctx, addresses[1])
		require.NoError(t, err)
		require.True(t, isObserved)
	})

	t.Run("no difference for uppercase and lowercase addresses", func(t *testing.T) {
		repo := NewObserverRepository()

		err := repo.ObserveAddress(ctx, addresses[1])
		require.NoError(t, err)

		isObserved, err := repo.IsAddressObserved(ctx, strings.ToUpper(addresses[1]))
		require.NoError(t, err)
		require.True(t, isObserved)
	})
}

func randomAddresses() []string {
	return []string{
		"0x056Fc2ceC04BF827d2A3a6e0A9588a05d6f87B57",
		"0x63feFeeD9eF48706B402A6b94Bf9F63747B3D5Da",
		"0x4d52a27740DD522F7f02E269bDE3AdB189da84aC",
	}
}
