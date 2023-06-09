// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package signing

import (
	"context"

	"github.com/zeebo/errs"

	"common/pb"
	"common/rpc/rpctracing"
	"common/storx"
)

// Error is the default error class for signing package.
var Error = errs.Class("signing")

// Signer is able to sign data and verify own signature belongs.
type Signer interface {
	ID() storx.NodeID
	HashAndSign(ctx context.Context, data []byte) ([]byte, error)
	HashAndVerifySignature(ctx context.Context, data, signature []byte) error
	SignHMACSHA256(ctx context.Context, data []byte) ([]byte, error)
	VerifyHMACSHA256(ctx context.Context, data, signature []byte) error
}

// SignOrderLimit signs the order limit using the specified signer.
// Signer is a satellite.
func SignOrderLimit(ctx context.Context, satellite Signer, unsigned *pb.OrderLimit) (_ *pb.OrderLimit, err error) {
	defer mon.Task()(&ctx)(&err)
	bytes, err := EncodeOrderLimit(ctx, unsigned)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	signed := *unsigned
	signed.SatelliteSignature, err = satellite.HashAndSign(ctx, bytes)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &signed, nil
}

// SignUplinkOrder signs the order using the specified signer.
// Signer is an uplink.
func SignUplinkOrder(ctx context.Context, privateKey storx.PiecePrivateKey, unsigned *pb.Order) (_ *pb.Order, err error) {
	ctx = rpctracing.WithoutDistributedTracing(ctx)
	defer mon.Task()(&ctx)(&err)

	bytes, err := EncodeOrder(ctx, unsigned)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	signed := *unsigned
	signed.UplinkSignature, err = privateKey.Sign(bytes)
	if err != nil {
		return nil, Error.Wrap(err)
	}
	return &signed, nil
}

// SignPieceHash signs the piece hash using the specified signer.
// Signer is either uplink or storage node.
func SignPieceHash(ctx context.Context, signer Signer, unsigned *pb.PieceHash) (_ *pb.PieceHash, err error) {
	defer mon.Task()(&ctx)(&err)
	bytes, err := EncodePieceHash(ctx, unsigned)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	signed := *unsigned
	signed.Signature, err = signer.HashAndSign(ctx, bytes)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &signed, nil
}

// SignUplinkPieceHash signs the piece hash using the specified signer.
// Signer is either uplink or storage node.
func SignUplinkPieceHash(ctx context.Context, privateKey storx.PiecePrivateKey, unsigned *pb.PieceHash) (_ *pb.PieceHash, err error) {
	defer mon.Task()(&ctx)(&err)
	bytes, err := EncodePieceHash(ctx, unsigned)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	signed := *unsigned
	signed.Signature, err = privateKey.Sign(bytes)
	if err != nil {
		return nil, Error.Wrap(err)
	}
	return &signed, nil
}

// SignExitCompleted signs the ExitCompleted using the specified signer.
// Signer is a satellite.
func SignExitCompleted(ctx context.Context, signer Signer, unsigned *pb.ExitCompleted) (_ *pb.ExitCompleted, err error) {
	defer mon.Task()(&ctx)(&err)
	bytes, err := EncodeExitCompleted(ctx, unsigned)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	signed := *unsigned
	signed.ExitCompleteSignature, err = signer.HashAndSign(ctx, bytes)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &signed, nil
}

// SignExitFailed signs the ExitFailed using the specified signer.
// Signer is a satellite.
func SignExitFailed(ctx context.Context, signer Signer, unsigned *pb.ExitFailed) (_ *pb.ExitFailed, err error) {
	defer mon.Task()(&ctx)(&err)
	bytes, err := EncodeExitFailed(ctx, unsigned)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	signed := *unsigned
	signed.ExitFailureSignature, err = signer.HashAndSign(ctx, bytes)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &signed, nil
}
