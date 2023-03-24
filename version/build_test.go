// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package version_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"common/testcontext"
)

func TestFromBuild(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "testbuild"

	data, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, `"v0.0.0-00010101000000-000000000000"`, string(data))
}
