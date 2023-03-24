// Copyright (C) 2022 Storx Labs, Inc.
// See LICENSE for copying information.

//go:build go1.18 && !go1.21 && !noquic
// +build go1.18,!go1.21,!noquic

package quic

import (
	"common/rpc"
)

const quicConnectorPriority = 20

func init() {
	rpc.RegisterCandidateConnectorType("quic", func() rpc.Connector {
		return NewDefaultConnector(nil)
	}, quicConnectorPriority)
}
