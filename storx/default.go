// Copyright (C) 2022 Storx Labs, Inc.
// See LICENSE for copying information.

package storx

import "encoding/base32"

// base32Encoding is base32 without padding.
var base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)
