// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package pb

import (
	"bytes"
	"reflect"

	"github.com/gogo/protobuf/proto"

	"common/storx"
)

// Equal compares two Protobuf messages via serialization.
func Equal(msg1, msg2 proto.Message) bool {
	// reflect.DeepEqual and proto.Equal don't seem work in all cases
	// todo:  see how slow this is compared to custom equality checks
	if msg1 == nil {
		return msg2 == nil
	}
	if reflect.TypeOf(msg1) != reflect.TypeOf(msg2) {
		return false
	}
	msg1Bytes, err := Marshal(msg1)
	if err != nil {
		return false
	}
	msg2Bytes, err := Marshal(msg2)
	if err != nil {
		return false
	}
	return bytes.Equal(msg1Bytes, msg2Bytes)
}

// NodesToIDs extracts Node-s into a list of ids.
func NodesToIDs(nodes []*Node) storx.NodeIDList {
	ids := make(storx.NodeIDList, len(nodes))
	for i, node := range nodes {
		if node != nil {
			ids[i] = node.Id
		}
	}
	return ids
}

// CopyNodeAddress returns a deep copy of a NodeAddress.
func CopyNodeAddress(src *NodeAddress) (dst *NodeAddress) {
	if src == nil {
		return nil
	}
	return proto.Clone(src).(*NodeAddress)
}

// CopyNode returns a deep copy of a node
// It would be better to use `proto.Clone` but it is curently incompatible
// with gogo's customtype extension.
// (see https://github.com/gogo/protobuf/issues/147).
func CopyNode(src *Node) (dst *Node) {
	node := Node{Id: storx.NodeID{}}
	copy(node.Id[:], src.Id[:])
	node.Address = CopyNodeAddress(src.Address)
	return &node
}

// AddressEqual compares two node addresses.
func AddressEqual(a1, a2 *NodeAddress) bool {
	if (a1 == nil) != (a2 == nil) {
		return false
	}
	if a1 == nil {
		return true
	}
	if a1.Address != a2.Address || a1.DebounceLimit != a2.DebounceLimit {
		return false
	}
	if (a1.NoiseInfo == nil) != (a2.NoiseInfo == nil) {
		return false
	}
	if a1.NoiseInfo == nil {
		return true
	}
	return a1.NoiseInfo.Proto == a2.NoiseInfo.Proto &&
		bytes.Equal(a1.NoiseInfo.GetPublicKey(), a2.NoiseInfo.GetPublicKey())

}

// NewRedundancySchemeToStorx creates new storx.RedundancyScheme from the given
// protobuf RedundancyScheme.
func NewRedundancySchemeToStorx(scheme *RedundancyScheme) *storx.RedundancyScheme {
	return &storx.RedundancyScheme{
		Algorithm:      storx.RedundancyAlgorithm(scheme.GetType()),
		ShareSize:      scheme.GetErasureShareSize(),
		RequiredShares: int16(scheme.GetMinReq()),
		RepairShares:   int16(scheme.GetRepairThreshold()),
		OptimalShares:  int16(scheme.GetSuccessThreshold()),
		TotalShares:    int16(scheme.GetTotal()),
	}
}

// Convert converts a *NoiseInfo to a storx.NoiseInfo.
func (n *NoiseInfo) Convert() (rv storx.NoiseInfo) {
	// TODO(jt): the existence of these functions is a
	// disastrous amount of unnecessary boilerplate. i get that
	// we didn't want the storx/common/storx package to have
	// to import github.com/gogo/proto, but at this point, having
	// all these runtime translation layers between a bunch of
	// types is the wrong tradeoff. we should figure out how to
	// make storx/common/pb broken up into a bunch of
	// lightweight type definitions, so we can use them and only
	// define them once. this switch statement could go away.
	if n == nil {
		return rv
	}
	rv.PublicKey = string(n.PublicKey)
	switch n.Proto {
	case NoiseProtocol_NOISE_UNSET:
		rv.Proto = storx.NoiseProto_Unset
	case NoiseProtocol_NOISE_IK_25519_CHACHAPOLY_BLAKE2B:
		rv.Proto = storx.NoiseProto_IK_25519_ChaChaPoly_BLAKE2b
	case NoiseProtocol_NOISE_IK_25519_AESGCM_BLAKE2B:
		rv.Proto = storx.NoiseProto_IK_25519_AESGCM_BLAKE2b
	default:
		rv.Proto = storx.NoiseProto_Unset
	}
	return rv
}

// NoiseInfoConvert converts a storx.NoiseInfo to a *NoiseInfo.
func NoiseInfoConvert(info storx.NoiseInfo) (rv *NoiseInfo) {
	if info.PublicKey == "" && info.Proto == storx.NoiseProto_Unset {
		return nil
	}
	rv = &NoiseInfo{}
	if info.PublicKey != "" {
		rv.PublicKey = []byte(info.PublicKey)
	}
	switch info.Proto {
	case storx.NoiseProto_Unset:
		rv.Proto = NoiseProtocol_NOISE_UNSET
	case storx.NoiseProto_IK_25519_ChaChaPoly_BLAKE2b:
		rv.Proto = NoiseProtocol_NOISE_IK_25519_CHACHAPOLY_BLAKE2B
	case storx.NoiseProto_IK_25519_AESGCM_BLAKE2b:
		rv.Proto = NoiseProtocol_NOISE_IK_25519_AESGCM_BLAKE2B
	default:
		rv.Proto = NoiseProtocol_NOISE_UNSET
	}
	return rv

}

// NodeURL converts a *Node to a storx.NodeURL.
func (n *Node) NodeURL() storx.NodeURL {
	return storx.NodeURL{
		ID:            n.Id,
		Address:       n.Address.Address,
		NoiseInfo:     n.Address.NoiseInfo.Convert(),
		DebounceLimit: int(n.Address.DebounceLimit),
	}
}

// NodeFromNodeURL converts a storx.NodeURL to a *Node.
func NodeFromNodeURL(u storx.NodeURL) *Node {
	return &Node{
		Id: u.ID,
		Address: &NodeAddress{
			Address:       u.Address,
			NoiseInfo:     NoiseInfoConvert(u.NoiseInfo),
			DebounceLimit: int32(u.DebounceLimit),
		},
	}
}
