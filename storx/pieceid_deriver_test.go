// Copyright (C) 2022 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/identity/testidentity"
	"common/storx"
	"common/testrand"
)

func TestPieceID_PieceDeriver(t *testing.T) {
	pieceIDA := storx.NewPieceID()
	pieceIDB := storx.NewPieceID()
	a := pieceIDA.Deriver()
	b := pieceIDB.Deriver()

	n0 := testidentity.MustPregeneratedIdentity(0, storx.LatestIDVersion()).ID
	n1 := testidentity.MustPregeneratedIdentity(1, storx.LatestIDVersion()).ID

	require.Equal(t, pieceIDA.Derive(n0, 1), a.Derive(n0, 1))
	require.Equal(t, pieceIDB.Derive(n1, 1), b.Derive(n1, 1))

	assert.NotEqual(t, a.Derive(n0, 0), a.Derive(n1, 0), "a(n0, 0) != a(n1, 0)")
	assert.NotEqual(t, b.Derive(n0, 0), b.Derive(n1, 0), "b(n0, 0) != b(n1, 0)")
	assert.NotEqual(t, a.Derive(n0, 0), b.Derive(n0, 0), "a(n0, 0) != b(n0, 0)")
	assert.NotEqual(t, a.Derive(n1, 0), b.Derive(n1, 0), "a(n1, 0) != b(n1, 0)")

	assert.NotEqual(t, a.Derive(n0, 0), a.Derive(n0, 1), "a(n0, 0) != a(n0, 1)")

	// idempotent
	assert.Equal(t, a.Derive(n0, 0), a.Derive(n0, 0), "a(n0, 0)")
	assert.Equal(t, a.Derive(n1, 0), a.Derive(n1, 0), "a(n1, 0)")

	assert.Equal(t, b.Derive(n0, 0), b.Derive(n0, 0), "b(n0, 0)")
	assert.Equal(t, b.Derive(n1, 0), b.Derive(n1, 0), "b(n1, 0)")
}

func TestPieceID_PieceDeriver_Golden(t *testing.T) {
	root, err := storx.PieceIDFromString("ZAMSQD6IMH5V6NXBTGRWQTQ4QLUGF7M7M5GDDQJ5YD64LN5QJJ5A")
	require.NoError(t, err)

	deriver := root.Deriver()
	verify := func(node string, pieceNum int32, expected string) {
		nodeid, err := storx.NodeIDFromString(node)
		require.NoError(t, err)

		derived := deriver.Derive(nodeid, pieceNum)
		require.Equal(t, expected, derived.String())
	}

	verify("1dYvWGgmzmerRxa2Rzv6dqjDogfCZE7dwSuDnfgaSfGT98GjQG", 0, "HPPSS55DNSNWN2YCNLQSFQVX3XIIAEIGNXZ2L5GODLDAXB6DVK4A")
	verify("1z3SQSAQjZxLZQ2sMQMbBtm12P3jSSTfjDnApo5Vu3XS7aGYAf", 1, "WGQN2CTMRQNRAL7YCGBORCWRNQI4PCGYYX2OUXMRP7QWTGEYO3AA")
	verify("12nmSD9xEp7EogWwEdKdQu4MwKwMdfmzFv7Cri5Jmyn1jVz3Yw4", 2, "SQM4U23GNGV4UVRLKEWO3ZXSUR6VSP2G6ESMZLPO6SS2DT6LLG2A")
	verify("1jA7TXHcCdZPf9kPLGa7L4KeRvD5xYmCrrLxHE9S39CbKjkoSH", 3, "MQVWWQWTOZ72B7VWNI5NQNVEN22AFB3U36E7KC2BYBVDPLIWMWGQ")
	verify("1kfXetMUDQM5YpN5Q5tMwSiZuYbQU7BofSjk9LYALrY5h2LaoM", 4, "OZJEUCWRWE7VCADPR6OADNOXWFJFESTMLVKH3ZEA2Y366EAFLIFQ")
	verify("16GRrTkVqo5fGVsWDXKhgtYcwhR7JtsG7PiGEyecavefEi5bzK", 5, "IWU4EVIZ564KNFURH57MW7S3TN64YW3NY3TR5KTW4THFJN6MBLLQ")
	verify("1pYwuNrMLsMAhu2kPuh1iT9uHVa5NGwyuy1gsiQtf21Xw2SvqY", 6, "2NG5NNLD3QONXJ2SR5TZ4VGK3Y5KVMJ7TZLDW2UGTMIMQ74HB5HQ")
	verify("12wKZHgYdqgcaWpWym9f8tzFRNx75DTVQ2kYhdwxUCWNvJRaRvV", 7, "HZ54AGD45N35TAY67LO4UFPFMFXJ2DFU62UQO5AQDWKPZXMTQSFQ")
	verify("1KMs9SyfKKq37xUkURKDMX7PVouqEZpTpefeiEKHiCm83JXYeh", 8, "FVVSEQ6RQ3Q3MGH4LBOM37XUH2TQUYGYVRKAWCQEEQVMHN7DLGLA")
	verify("1256p7W3uttfHoBqeKR2BaNyM7ZsJF3PeZqMbfRNR5fPBZja3kz", 9, "OL4TZZYOHVDXT6V2BMDADPKAPNC7G7PEPAY4CLCQKPICYD6JEMHQ")
}

func BenchmarkDeriver(b *testing.B) {
	pieceID := storx.NewPieceID()
	n0 := testidentity.MustPregeneratedIdentity(0, storx.LatestIDVersion()).ID

	b.Run("Derive", func(b *testing.B) {
		for k := 0; k < b.N; k++ {
			sink = pieceID.Derive(n0, 0)
		}
	})

	b.Run("Deriver", func(b *testing.B) {
		deriver := pieceID.Deriver()
		for k := 0; k < b.N; k++ {
			sink = deriver.Derive(n0, 0)
		}
	})
}

var sink storx.PieceID

func BenchmarkDeriver100(b *testing.B) {
	nodes := make([]storx.NodeID, 100)
	for i := range nodes {
		nodes[i] = testrand.NodeID()
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pieceID := storx.PieceID{byte(i), byte(i >> 8)}
		deriver := pieceID.Deriver()
		for i := range nodes {
			sink = deriver.Derive(nodes[i], int32(i))
		}
	}
}
