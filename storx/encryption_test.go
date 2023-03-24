// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/storx"
	"common/testrand"
)

func TestNewKey(t *testing.T) {
	t.Run("nil humanReadableKey", func(t *testing.T) {
		t.Parallel()

		key, err := storx.NewKey(nil)
		require.NoError(t, err)
		require.True(t, key.IsZero(), "key isn't zero value")
	})

	t.Run("empty humanReadableKey", func(t *testing.T) {
		t.Parallel()

		key, err := storx.NewKey([]byte{})
		require.NoError(t, err)
		require.True(t, key.IsZero(), "key isn't zero value")
	})

	t.Run("humanReadableKey is of KeySize length", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.Bytes(storx.KeySize)

		key, err := storx.NewKey(humanReadableKey)
		require.NoError(t, err)
		require.Equal(t, humanReadableKey, key[:])
	})

	t.Run("humanReadableKey is shorter than KeySize", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.BytesInt(testrand.Intn(storx.KeySize))

		key, err := storx.NewKey(humanReadableKey)
		require.NoError(t, err)
		require.Equal(t, humanReadableKey, key[:len(humanReadableKey)])
	})

	t.Run("humanReadableKey is larger than KeySize", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.BytesInt(testrand.Intn(10) + storx.KeySize + 1)

		key, err := storx.NewKey(humanReadableKey)
		require.NoError(t, err)
		assert.Equal(t, humanReadableKey[:storx.KeySize], key[:])
	})

	t.Run("same human readable key produce the same key", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.BytesInt(testrand.Intn(10) + storx.KeySize + 1)

		key1, err := storx.NewKey(humanReadableKey)
		require.NoError(t, err)

		key2, err := storx.NewKey(humanReadableKey)
		require.NoError(t, err)

		assert.Equal(t, key1, key2, "keys are equal")
	})
}

func TestKey_IsZero(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var key *storx.Key
		require.True(t, key.IsZero())

		wrapperFn := func(key *storx.Key) bool {
			return key.IsZero()
		}
		require.True(t, wrapperFn(nil))
	})

	t.Run("zero", func(t *testing.T) {
		key := &storx.Key{}
		require.True(t, key.IsZero())
	})

	t.Run("no nil/zero", func(t *testing.T) {
		key := &storx.Key{'k'}
		require.False(t, key.IsZero())
	})
}

// TestNonce_Scan tests (*Nonce).Scan().
func TestNonce_Scan(t *testing.T) {
	tmp := storx.Nonce{}
	require.Error(t, tmp.Scan(32))
	require.Error(t, tmp.Scan(false))
	require.Error(t, tmp.Scan([]byte{}))

	require.NoError(t, tmp.Scan(nil))
	require.True(t, tmp.IsZero())
	require.NoError(t, tmp.Scan(tmp.Bytes()))
	require.True(t, tmp.IsZero())
}

// TestEncryptedPrivateKey_Scan tests (*EncryptedPrivateKey).Scan().
func TestEncryptedPrivateKey_Scan(t *testing.T) {
	tmp := storx.EncryptedPrivateKey{}
	require.Error(t, tmp.Scan(32))
	require.Error(t, tmp.Scan(false))
	require.NoError(t, tmp.Scan([]byte{}))
	require.NoError(t, tmp.Scan([]byte{1, 2, 3, 4}))

	ref := []byte{1, 2, 3}
	require.NoError(t, tmp.Scan(ref))
	ref[0] = 0xFF
	require.Equal(t, storx.EncryptedPrivateKey{1, 2, 3}, tmp)
}
