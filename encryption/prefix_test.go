// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"

	"common/paths"
	"common/storx"
)

func TestPrefixInfo(t *testing.T) {
	store := NewStore()
	store.SetDefaultKey(new(storx.Key))
	store.SetDefaultPathCipher(storx.EncAESGCM)

	ep := paths.NewEncrypted
	up := paths.NewUnencrypted

	type PathPair struct {
		Unenc string
		Enc   string
	}

	type TestCase struct {
		Path   PathPair
		Parent PathPair
	}

	encrypt := func(path string) string {
		t.Helper()

		enc, err := EncryptPathWithStoreCipher("bucket", up(path), store)
		require.NoError(t, err)
		return enc.Raw()
	}

	// define the set of encrypted<->unencrypted pairs of strings the test will use.
	// note that this is not a bijection: multiple strings decrypt into the empty string.
	// specifically, both "" and encryptPathComponent("", storx.EncAESGCM, key).
	// that is why there is both pairNone and pairEmpty. pairNone is in the parent slot
	// when there is no parent component, and pairEmpty is used in the parent slot when
	// the parent is the empty string.

	var (
		pairNone            = PathPair{"", ""}
		pairEmpty           = PathPair{"", encrypt("/")[:13]}
		pairEmptyEmpty      = PathPair{"/", encrypt("/")}
		pairEmptyEmptyEmpty = PathPair{"//", encrypt("//")}
		pairFoo             = PathPair{"foo", encrypt("foo")}
		pairFooBar          = PathPair{"foo/bar", encrypt("foo/bar")}
		pairFooEmpty        = PathPair{"foo/", encrypt("foo/")}
		pairFooEmptyBar     = PathPair{"foo//bar", encrypt("foo//bar")}
	)

	// define the tests to run. the variables are named such that pairFooBarBaz represents
	// the path "foo/bar/baz".

	tests := []TestCase{
		{Path: pairNone, Parent: pairNone},
		{Path: pairEmptyEmpty, Parent: pairEmpty},
		{Path: pairEmptyEmptyEmpty, Parent: pairEmptyEmpty},
		{Path: pairFoo, Parent: pairNone},
		{Path: pairFooBar, Parent: pairFoo},
		{Path: pairFooEmpty, Parent: pairFoo},
		{Path: pairFooEmptyBar, Parent: pairFooEmpty},
	}

	for _, test := range tests {
		pi, err := GetPrefixInfo("bucket", up(test.Path.Unenc), store)
		require.NoError(t, err)
		require.NotNil(t, pi)

		require.Equal(t, pi.Bucket, "bucket")
		require.Equal(t, pi.Cipher, storx.EncAESGCM)
		require.Equal(t, pi.PathUnenc, up(test.Path.Unenc))
		require.Equal(t, pi.PathEnc, ep(test.Path.Enc))
		require.Equal(t, pi.ParentUnenc, up(test.Parent.Unenc))
		require.Equal(t, pi.ParentEnc, ep(test.Parent.Enc))
	}
}
