// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package fpath_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"common/fpath"
	"common/testcontext"
)

func TestAtomicWriteFile(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	err := fpath.AtomicWriteFile(ctx.File("example.txt"), []byte{1, 2, 3}, 0600)
	require.NoError(t, err)

	data, err := os.ReadFile(ctx.File("example.txt"))
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3}, data)
}
