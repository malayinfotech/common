// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx

import (
	"time"

	"github.com/zeebo/errs"

	"common/uuid"
)

var (
	// ErrBucket is an error class for general bucket errors.
	ErrBucket = errs.Class("bucket")

	// ErrNoBucket is an error class for using empty bucket name.
	ErrNoBucket = errs.Class("no bucket specified")

	// ErrBucketNotFound is an error class for non-existing bucket.
	ErrBucketNotFound = errs.Class("bucket not found")

	// ErrBucketNotEmpty is an error class for using non-empty bucket in operation that requires empty bucket.
	ErrBucketNotEmpty = errs.Class("bucket must be empty")
)

// Bucket contains information about a specific bucket.
type Bucket struct {
	ID                          uuid.UUID
	Name                        string
	ProjectID                   uuid.UUID
	PartnerID                   uuid.UUID
	UserAgent                   []byte
	Created                     time.Time
	PathCipher                  CipherSuite
	DefaultSegmentsSize         int64
	DefaultRedundancyScheme     RedundancyScheme
	DefaultEncryptionParameters EncryptionParameters
	Placement                   PlacementConstraint
}
