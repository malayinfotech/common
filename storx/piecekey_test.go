// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"common/memory"
	"common/storx"
	"common/testrand"
)

func TestPublicPrivatePieceKey(t *testing.T) {
	expectedPublicKey, expectedPrivateKey, err := storx.NewPieceKey()
	require.NoError(t, err)

	publicKey, err := storx.PiecePublicKeyFromBytes(expectedPublicKey.Bytes())
	require.NoError(t, err)
	require.Equal(t, expectedPublicKey, publicKey)

	privateKey, err := storx.PiecePrivateKeyFromBytes(expectedPrivateKey.Bytes())
	require.NoError(t, err)
	require.Equal(t, expectedPrivateKey, privateKey)

	{
		data := testrand.Bytes(10 * memory.KiB)
		signature, err := privateKey.Sign(data)
		require.NoError(t, err)

		err = publicKey.Verify(data, signature)
		require.NoError(t, err)

		err = publicKey.Verify(data, testrand.BytesInt(32))
		require.Error(t, err)

		err = publicKey.Verify(testrand.Bytes(10*memory.KiB), signature)
		require.Error(t, err)
	}

	{
		// to small
		_, err = storx.PiecePublicKeyFromBytes([]byte{1})
		require.Error(t, err)

		// to small
		_, err = storx.PiecePrivateKeyFromBytes([]byte{1})
		require.Error(t, err)

		// to large
		_, err = storx.PiecePublicKeyFromBytes(testrand.Bytes(33))
		require.Error(t, err)

		// to large
		_, err = storx.PiecePrivateKeyFromBytes(testrand.Bytes(65))
		require.Error(t, err)

		// public key from private
		_, err = storx.PiecePublicKeyFromBytes(expectedPrivateKey.Bytes())
		require.Error(t, err)

		// private key from public
		_, err = storx.PiecePrivateKeyFromBytes(expectedPublicKey.Bytes())
		require.Error(t, err)
	}
}
