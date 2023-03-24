// Copyright (C) 2023 Storx Labs, Inc.
// See LICENSE for copying information.

package storx

import (
	"fmt"
	"net/url"

	"common/base58"
)

// NoiseProto represents different possible Noise handshake and cipher suite
// selections.
type NoiseProto int

const (
	// NoiseProto_Unset is an unset protocol.
	NoiseProto_Unset = 0
	// NoiseProto_IK_25519_ChaChaPoly_BLAKE2b is a Noise protocol.
	NoiseProto_IK_25519_ChaChaPoly_BLAKE2b = 1
	// NoiseProto_IK_25519_AESGCM_BLAKE2b is a Noise protocol.
	NoiseProto_IK_25519_AESGCM_BLAKE2b = 2
)

// NoiseInfo represents the information needed to dial a remote Noise peer.
type NoiseInfo struct {
	Proto     NoiseProto
	PublicKey string // byte representation
}

// WriteTo assists in serializing a NoiseInfo to a NodeURL.
func (info *NoiseInfo) WriteTo(values url.Values) {
	if info.Proto != NoiseProto_Unset {
		values.Set("noise_proto", fmt.Sprint(int(info.Proto)))
	}
	if info.PublicKey != "" {
		values.Set("noise_pub", base58.CheckEncode([]byte(info.PublicKey), 0))
	}
}
