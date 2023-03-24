// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package context2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"common/context2"
	"common/testcontext"
)

func TestWithoutCancellation(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	parent, cancel := context.WithCancel(ctx)
	cancel()

	without := context2.WithoutCancellation(parent)
	require.Equal(t, error(nil), without.Err())
	require.Equal(t, (<-chan struct{})(nil), without.Done())
}
