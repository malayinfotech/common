// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/paths"
	"common/storx"
)

func consumeIter(iter paths.Iterator) string {
	var parts []string
	for !iter.Done() {
		parts = append(parts, iter.Next())
	}
	return fmt.Sprintf("%q", parts)
}

func printLookup(revealed map[string]string, remaining paths.Iterator, base *Base) {
	if base == nil {
		fmt.Printf("<%q, %s, nil>\n", revealed, consumeIter(remaining))
	} else {
		fmt.Printf("<%q, %s, <%q, %q, %q, %v>>\n",
			revealed, consumeIter(remaining), base.Unencrypted, base.Encrypted, base.Key[:2], base.Default)
	}
}

func toKey(val string) (out storx.Key) {
	copy(out[:], val)
	return out
}

func abortIfError(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func ExampleStore() {
	s := NewStore()
	ep := paths.NewEncrypted
	up := paths.NewUnencrypted

	// Add a fairly complicated tree to the store.
	abortIfError(s.AddWithCipher("b1", up("u1/u2/u3"), ep("e1/e2/e3"), toKey("k3"), storx.EncAESGCM))
	abortIfError(s.AddWithCipher("b1", up("u1/u2/u3/u4"), ep("e1/e2/e3/e4"), toKey("k4"), storx.EncAESGCM))
	abortIfError(s.AddWithCipher("b1", up("u1/u5"), ep("e1/e5"), toKey("k5"), storx.EncAESGCM))
	abortIfError(s.AddWithCipher("b1", up("u6"), ep("e6"), toKey("k6"), storx.EncAESGCM))
	abortIfError(s.AddWithCipher("b1", up("u6/u7/u8"), ep("e6/e7/e8"), toKey("k8"), storx.EncAESGCM))
	abortIfError(s.AddWithCipher("b2", up("u1"), ep("e1'"), toKey("k1"), storx.EncAESGCM))
	abortIfError(s.AddWithCipher("b3", paths.Unencrypted{}, paths.Encrypted{}, toKey("m1"), storx.EncAESGCM))

	// Look up some complicated queries by the unencrypted path.
	printLookup(s.LookupUnencrypted("b1", up("u1")))
	printLookup(s.LookupUnencrypted("b1", up("u1/u2/u3")))
	printLookup(s.LookupUnencrypted("b1", up("u1/u2/u3/u6")))
	printLookup(s.LookupUnencrypted("b1", up("u1/u2/u3/u4")))
	printLookup(s.LookupUnencrypted("b1", up("u6/u7")))
	printLookup(s.LookupUnencrypted("b2", up("u1")))
	printLookup(s.LookupUnencrypted("b3", paths.Unencrypted{}))
	printLookup(s.LookupUnencrypted("b3", up("z1")))

	fmt.Println()

	// Look up some complicated queries by the encrypted path.
	printLookup(s.LookupEncrypted("b1", ep("e1")))
	printLookup(s.LookupEncrypted("b1", ep("e1/e2/e3")))
	printLookup(s.LookupEncrypted("b1", ep("e1/e2/e3/e6")))
	printLookup(s.LookupEncrypted("b1", ep("e1/e2/e3/e4")))
	printLookup(s.LookupEncrypted("b1", ep("e6/e7")))
	printLookup(s.LookupEncrypted("b2", ep("e1'")))
	printLookup(s.LookupEncrypted("b3", paths.Encrypted{}))
	printLookup(s.LookupEncrypted("b3", ep("z1")))

	// output:
	//
	// <map["e2":"u2" "e5":"u5"], [], nil>
	// <map["e4":"u4"], [], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	// <map[], ["u6"], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	// <map[], [], <"u1/u2/u3/u4", "e1/e2/e3/e4", "k4", false>>
	// <map["e8":"u8"], ["u7"], <"u6", "e6", "k6", false>>
	// <map[], [], <"u1", "e1'", "k1", false>>
	// <map[], [], <"", "", "m1", false>>
	// <map[], ["z1"], <"", "", "m1", false>>
	//
	// <map["u2":"e2" "u5":"e5"], [], nil>
	// <map["u4":"e4"], [], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	// <map[], ["e6"], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	// <map[], [], <"u1/u2/u3/u4", "e1/e2/e3/e4", "k4", false>>
	// <map["u8":"e8"], ["e7"], <"u6", "e6", "k6", false>>
	// <map[], [], <"u1", "e1'", "k1", false>>
	// <map[], [], <"", "", "m1", false>>
	// <map[], ["z1"], <"", "", "m1", false>>
}

func ExampleStore_SetDefaultKey() {
	s := NewStore()
	dk := toKey("dk")
	s.SetDefaultKey(&dk)
	ep := paths.NewEncrypted
	up := paths.NewUnencrypted

	abortIfError(s.AddWithCipher("b1", up("u1/u2/u3"), ep("e1/e2/e3"), toKey("k3"), storx.EncAESGCM))

	printLookup(s.LookupUnencrypted("b1", up("u1")))
	printLookup(s.LookupUnencrypted("b1", up("u1/u2")))
	printLookup(s.LookupUnencrypted("b1", up("u1/u2/u3")))
	printLookup(s.LookupUnencrypted("b1", up("u1/u2/u3/u4")))

	fmt.Println()

	printLookup(s.LookupEncrypted("b1", ep("e1")))
	printLookup(s.LookupEncrypted("b1", ep("e1/e2")))
	printLookup(s.LookupEncrypted("b1", ep("e1/e2/e3")))
	printLookup(s.LookupEncrypted("b1", ep("e1/e2/e3/e4")))

	// output:
	//
	// <map[], ["u1"], <"", "", "dk", true>>
	// <map[], ["u1" "u2"], <"", "", "dk", true>>
	// <map[], [], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	// <map[], ["u4"], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	//
	// <map[], ["e1"], <"", "", "dk", true>>
	// <map[], ["e1" "e2"], <"", "", "dk", true>>
	// <map[], [], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
	// <map[], ["e4"], <"u1/u2/u3", "e1/e2/e3", "k3", false>>
}

func TestStoreErrors(t *testing.T) {
	for _, pathCipher := range []storx.CipherSuite{
		storx.EncNull,
		storx.EncAESGCM,
		storx.EncSecretBox,
	} {
		s := NewStore()
		ep := paths.NewEncrypted
		up := paths.NewUnencrypted

		// Too many encrypted parts
		require.Error(t, s.AddWithCipher("b1", up("u1"), ep("e1/e2/e3"), storx.Key{}, pathCipher))

		// Too many unencrypted parts
		require.Error(t, s.AddWithCipher("b1", up("u1/u2/u3"), ep("e1"), storx.Key{}, pathCipher))

		// Mismatches
		require.NoError(t, s.AddWithCipher("b1", up("u1"), ep("e1"), storx.Key{}, pathCipher))
		require.Error(t, s.AddWithCipher("b1", up("u2"), ep("e1"), storx.Key{}, pathCipher))
		require.Error(t, s.AddWithCipher("b1", up("u1"), ep("f1"), storx.Key{}, pathCipher))
	}
}

func TestStoreErrorState(t *testing.T) {
	s := NewStore()
	ep := paths.NewEncrypted
	up := paths.NewUnencrypted

	// Do an empty lookup.
	revealed1, consumed1, base1 := s.LookupUnencrypted("b1", up("u1/u2"))

	// Attempt to do an addition that fails.
	require.Error(t, s.AddWithCipher("b1", up("u1/u2"), ep("e1/e2/e3"), storx.Key{}, storx.EncNull))
	require.Error(t, s.AddWithCipher("b1", up("u1/u2"), ep("e1/e2/e3"), storx.Key{}, storx.EncAESGCM))
	require.Error(t, s.AddWithCipher("b1", up("u1/u2"), ep("e1/e2/e3"), storx.Key{}, storx.EncSecretBox))

	// Ensure that we get the same results as before
	revealed2, consumed2, base2 := s.LookupUnencrypted("b1", up("u1/u2"))

	assert.Equal(t, revealed1, revealed2)
	assert.Equal(t, consumed1, consumed2)
	assert.Equal(t, base1, base2)
}

func TestStoreIterate(t *testing.T) {
	type storeEntry struct {
		bucket     string
		unenc      paths.Unencrypted
		enc        paths.Encrypted
		key        storx.Key
		pathCipher storx.CipherSuite
	}

	for _, pathCipher := range []storx.CipherSuite{
		storx.EncNull,
		storx.EncAESGCM,
		storx.EncSecretBox,
	} {
		for _, bypass := range []bool{false, true} {
			s := NewStore()
			s.EncryptionBypass = bypass

			ep := paths.NewEncrypted
			up := paths.NewUnencrypted

			expected := map[storeEntry]struct{}{
				{"b1", up("u1/u2/u3"), ep("e1/e2/e3"), toKey("k3"), pathCipher}:         {},
				{"b1", up("u1/u2/u3/u4"), ep("e1/e2/e3/e4"), toKey("k4"), pathCipher}:   {},
				{"b1", up("u1/u5"), ep("e1/e5"), toKey("k5"), pathCipher}:               {},
				{"b1", up("u6"), ep("e6"), toKey("k6"), pathCipher}:                     {},
				{"b1", up("u6/u7/u8"), ep("e6/e7/e8"), toKey("k8"), pathCipher}:         {},
				{"b2", up("u1"), ep("e1'"), toKey("k1"), pathCipher}:                    {},
				{"b3", paths.Unencrypted{}, paths.Encrypted{}, toKey("m1"), pathCipher}: {},
			}

			for entry := range expected {
				require.NoError(t, s.AddWithCipher(entry.bucket, entry.unenc, entry.enc, entry.key, entry.pathCipher))
			}

			got := make(map[storeEntry]struct{})
			require.NoError(t, s.IterateWithCipher(func(bucket string, unenc paths.Unencrypted, enc paths.Encrypted, key storx.Key, pathCipher storx.CipherSuite) error {
				got[storeEntry{bucket, unenc, enc, key, pathCipher}] = struct{}{}
				return nil
			}))
			require.Equal(t, expected, got)
		}
	}
}

func TestStoreEncryptionBypass(t *testing.T) {
	s := NewStore()
	s.SetDefaultKey(new(storx.Key))
	s.SetDefaultPathCipher(storx.EncAESGCM)

	{
		_, _, base := s.LookupUnencrypted("bucket", paths.NewUnencrypted(""))
		require.Equal(t, base.PathCipher, storx.EncAESGCM)
	}

	s.EncryptionBypass = true

	{
		_, _, base := s.LookupUnencrypted("bucket", paths.NewUnencrypted(""))
		require.Equal(t, base.PathCipher, storx.EncNullBase64URL)
	}
}
