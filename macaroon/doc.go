// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

// Package macaroon implements contextual caveats and authorization.
package macaroon

//go:generate protoc  -I../pb -I. --pico_out=paths=source_relative:. types.proto
//go:generate goimports -local storx -w .
