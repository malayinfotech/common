// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package testpeertls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/identity"
	"common/identity/testidentity"
	"common/peertls"
	"common/pkcrypto"
	"common/storx"
)

func TestNewCertChain(t *testing.T) {
	testidentity.CompleteIdentityVersionsTest(t, func(t *testing.T, version storx.IDVersion, ident *identity.FullIdentity) {
		for length := 2; length < 4; length++ {
			t.Logf("length: %d", length)
			keys, chain, err := NewCertChain(length, version.Number)
			require.NoError(t, err)

			assert.Len(t, chain, length)
			assert.Len(t, keys, length)

			pubKey, err := pkcrypto.PublicKeyFromPrivate(keys[peertls.CAIndex])
			require.NoError(t, err)
			assert.Equal(t, pubKey, chain[peertls.CAIndex].PublicKey)

			pubKey, err = pkcrypto.PublicKeyFromPrivate(keys[peertls.LeafIndex])
			require.NoError(t, err)
			assert.Equal(t, pubKey, chain[peertls.LeafIndex].PublicKey)

			err = peertls.VerifyPeerCertChains(nil, identity.ToChains(chain))
			assert.NoError(t, err)

			assert.True(t, chain[peertls.CAIndex].IsCA)
			assert.False(t, chain[peertls.LeafIndex].IsCA)
		}
	})
}
