// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package tlsopts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/identity"
	"common/identity/testidentity"
	"common/peertls/tlsopts"
	"common/storx"
)

func TestVerifyIdentity_success(t *testing.T) {
	for i := 0; i < 50; i++ {
		ident, err := testidentity.PregeneratedIdentity(i, storx.LatestIDVersion())
		require.NoError(t, err)

		err = tlsopts.VerifyIdentity(ident.ID)(nil, identity.ToChains(ident.Chain()))
		assert.NoError(t, err)
	}
}

func TestVerifyIdentity_success_signed(t *testing.T) {
	for i := 0; i < 50; i++ {
		ident, err := testidentity.PregeneratedSignedIdentity(i, storx.LatestIDVersion())
		require.NoError(t, err)

		err = tlsopts.VerifyIdentity(ident.ID)(nil, identity.ToChains(ident.Chain()))
		assert.NoError(t, err)
	}
}

func TestVerifyIdentity_error(t *testing.T) {
	ident, err := testidentity.PregeneratedIdentity(0, storx.LatestIDVersion())
	require.NoError(t, err)

	identTheftVictim, err := testidentity.PregeneratedIdentity(1, storx.LatestIDVersion())
	require.NoError(t, err)

	cases := []struct {
		test   string
		nodeID storx.NodeID
	}{
		{"empty node ID", storx.NodeID{}},
		{"garbage node ID", storx.NodeID{0, 1, 2, 3}},
		{"wrong node ID", identTheftVictim.ID},
	}

	for _, cc := range cases {
		testCase := cc
		t.Run(testCase.test, func(t *testing.T) {
			err := tlsopts.VerifyIdentity(testCase.nodeID)(nil, identity.ToChains(ident.Chain()))
			assert.Error(t, err)
		})
	}
}
