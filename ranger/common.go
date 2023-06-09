// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package ranger

import (
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
)

// Error is the errs class of standard Ranger errors.
var Error = errs.Class("ranger")

var mon = monkit.Package()
