// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"os"

	"common/version"
)

func main() {
	versionstr, err := version.FromBuild("common")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("%#v", versionstr)
}
