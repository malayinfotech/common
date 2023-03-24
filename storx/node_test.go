// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/base58"
	"common/storx"
	"common/testrand"
)

func TestNodeID_Difficulty(t *testing.T) {
	invalidID := storx.NodeID{}
	difficulty, err := invalidID.Difficulty()
	assert.Error(t, err)
	assert.Equal(t, uint16(0), difficulty)

	for _, testcase := range []struct {
		id         string
		difficulty uint16
	}{
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dee00", 9},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dec00", 10},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113de800", 11},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113d7000", 12},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113de000", 13},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dc000", 14},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113d8000", 15},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa11390000", 16},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113e0000", 17},
	} {

		decoded, err := hex.DecodeString(testcase.id)
		if !assert.NoError(t, err) {
			t.Fatal()
		}

		var nodeID storx.NodeID
		n := copy(nodeID[:], decoded)
		if !assert.Equal(t, n, len(nodeID)) {
			t.Fatal()
		}

		difficulty, err := nodeID.Difficulty()
		if !assert.NoError(t, err) {
			t.Fatal()
		}

		assert.Equal(t, testcase.difficulty, difficulty)
	}
}

// TestNodeScan tests (*NodeID).Scan().
func TestNodeScan(t *testing.T) {
	tmpID := &storx.NodeID{}
	require.Error(t, tmpID.Scan(32))
	require.Error(t, tmpID.Scan(false))
	require.Error(t, tmpID.Scan([]byte{}))
	require.NoError(t, tmpID.Scan(tmpID.Bytes()))
}

// TestNodeValue tests NodeID.Value().
func TestNodeValue(t *testing.T) {
	tmpID := storx.NodeID{}
	v, err := tmpID.Value()
	require.NoError(t, err)
	require.IsType(t, v, []byte{})
	require.Len(t, v, storx.NodeIDSize)
}

func TestNodeID_Version(t *testing.T) {
	for _, testcase := range []struct {
		id         string
		difficulty uint16
		version    storx.IDVersionNumber
	}{
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113de500", 8, storx.V0},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dee00", 9, storx.V0},
		{"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dec00", 10, storx.V0},
	} {
		decoded, err := hex.DecodeString(testcase.id)
		require.NoError(t, err)

		var nodeID storx.NodeID
		n := copy(nodeID[:], decoded)
		require.Equal(t, n, len(nodeID))

		difficulty, err := nodeID.Difficulty()
		require.NoError(t, err)

		assert.Equal(t, testcase.difficulty, difficulty)
		assert.Equal(t, testcase.version, nodeID.Version().Number)
	}
}

func TestNodeID_String_Version(t *testing.T) {
	for _, testcase := range []struct {
		hexID    string
		base58ID string
		version  storx.IDVersionNumber
	}{
		{
			"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113de500",
			"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dBZN6Lg2T",
			storx.V0,
		},
		{
			"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dee00",
			"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dG3JN2sdZ",
			storx.V0,
		},
		{
			"fda09d6bed970d7a38fe7389cd2b1b9620cf0ea1fcda2404d353c3fa113dec00",
			"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7",
			storx.V0,
		},
	} {
		decoded, err := hex.DecodeString(testcase.hexID)
		require.NoError(t, err)

		var nodeID storx.NodeID
		n := copy(nodeID[:], decoded)
		require.Equal(t, n, len(nodeID))

		base58Str := nodeID.String()
		binID, version, err := base58.CheckDecode(base58Str)
		require.NoError(t, err)

		idVersion, err := storx.GetIDVersion(storx.IDVersionNumber(version))
		require.NoError(t, err)

		assert.Equal(t, testcase.version, idVersion.Number)
		assert.Equal(t, nodeID[:storx.NodeIDSize-1], binID[:storx.NodeIDSize-1])
	}
}

func TestNodeID_Compare(t *testing.T) {
	var a storx.NodeID
	require.Equal(t, 0, a.Compare(a)) //nolint:gocritic

	for k := 0; k < len(storx.NodeID{}); k++ {
		var a, b storx.NodeID
		a[k], b[k] = 1, 2
		require.Equal(t, 0, a.Compare(a)) //nolint:gocritic
		require.Equal(t, 0, b.Compare(b)) //nolint:gocritic
		require.Equal(t, -1, a.Compare(b))
		require.Equal(t, 1, b.Compare(a))
	}

	for k := 0; k < 100; k++ {
		var x, y storx.NodeID
		a, b := testrand.Int63n(math.MaxInt64), testrand.Int63n(math.MaxInt64)
		binary.BigEndian.PutUint64(x[:], uint64(a))
		binary.BigEndian.PutUint64(y[:], uint64(b))
		require.Equal(t, comp(a, b), x.Compare(y))
	}
}

func comp(a, b int64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func TestNodeID_MarshalJSON(t *testing.T) {
	nodeID, _ := storx.NodeIDFromString("12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7")
	buf, err := json.Marshal(nodeID)
	require.NoError(t, err)
	assert.Equal(t, string(buf), `"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7"`)
}

func TestNodeID_UnmarshalJSON(t *testing.T) {
	var nodeID storx.NodeID
	err := json.Unmarshal([]byte(`"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7"`), &nodeID)
	require.NoError(t, err)
	assert.Equal(t, nodeID.String(), "12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7")

	assert.Error(t, json.Unmarshal([]byte(`""12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7""`), &nodeID))
	assert.Error(t, json.Unmarshal([]byte(`{}`), &nodeID))
	assert.Error(t, json.Unmarshal([]byte(`"fajsd8ccj23m k,ladlf"`), &nodeID))
}

func TestNewVersionedID(t *testing.T) {
	nodeID := testrand.NodeID()

	assert.Equal(t, storx.V0, nodeID.Version().Number)

	for versionNumber, version := range storx.IDVersions {
		versionedNodeID := storx.NewVersionedID(nodeID, version)
		assert.Equal(t, versionNumber, versionedNodeID.Version().Number)
		assert.Equal(t, versionNumber, storx.IDVersionNumber(versionedNodeID[storx.NodeIDSize-1]))
	}
}

func TestNodeIDList_Contains(t *testing.T) {
	nodes := storx.NodeIDList{testrand.NodeID(), testrand.NodeID(), testrand.NodeID()}

	for _, testcase := range []struct {
		list     storx.NodeIDList
		id       storx.NodeID
		contains bool
	}{
		{storx.NodeIDList{}, storx.NodeID{}, false},
		{storx.NodeIDList{}, nodes[0], false},
		{nodes, storx.NodeID{}, false},
		{nodes, nodes[0], true},
		{nodes, nodes[1], true},
		{nodes, nodes[2], true},
		{nodes[0:1], nodes[0], true},
		{nodes[0:1], nodes[2], false},
	} {
		assert.Equal(t, testcase.contains, testcase.list.Contains(testcase.id))
	}
}

func TestUniqueNodeIDs(t *testing.T) {
	var IDs storx.NodeIDList
	id := testrand.NodeID()
	IDs = append(IDs, testrand.NodeID(), id, testrand.NodeID(), id, testrand.NodeID(), id, testrand.NodeID(), id, id)

	list := IDs.Unique()
	assert.Equal(t, len(list), 5)
}

func BenchmarkNodeID_Less(b *testing.B) {
	a := testrand.NodeID()
	b.Run("Same", func(b *testing.B) {
		total := 0
		x, y := a, a
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})

	b.Run("First", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[0]++
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})

	b.Run("Middle", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)/2]++
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})

	b.Run("Last", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)-1]++
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})
}

func btoi(v bool) int {
	if v {
		return 1
	}
	return 0
}

func BenchmarkNodeID_Compare(b *testing.B) {
	a := testrand.NodeID()
	b.Run("Same", func(b *testing.B) {
		total := 0
		x, y := a, a
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})

	b.Run("First", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[0]++
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})

	b.Run("Middle", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)/2]++
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})

	b.Run("Last", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)-1]++
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})
}
