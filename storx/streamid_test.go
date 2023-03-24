// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"common/storx"
	"common/testrand"
)

func TestStreamID_Encode(t *testing.T) {
	for i := 0; i < 10; i++ {
		expectedSize := testrand.Intn(255)
		streamID := testrand.StreamID(expectedSize)

		fromString, err := storx.StreamIDFromString(streamID.String())
		require.NoError(t, err)
		require.Equal(t, streamID.String(), fromString.String())

		fromBytes, err := storx.StreamIDFromBytes(streamID.Bytes())
		require.NoError(t, err)
		require.Equal(t, streamID.Bytes(), fromBytes.Bytes())

		require.Equal(t, streamID, fromString)
		require.Equal(t, expectedSize, fromString.Size())
		require.Equal(t, streamID, fromBytes)
		require.Equal(t, expectedSize, fromBytes.Size())
	}
}
