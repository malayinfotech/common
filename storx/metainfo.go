// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx

import "github.com/zeebo/errs"

var (
	// ErrNoPath is an error class for using empty path.
	ErrNoPath = errs.Class("no path specified")

	// ErrObjectNotFound is an error class for non-existing object.
	ErrObjectNotFound = errs.Class("object not found")
)

// ListDirection specifies listing direction.
type ListDirection int8

const (
	// Before lists backwards from cursor, without cursor [NOT SUPPORTED].
	Before = ListDirection(-2)
	// Backward lists backwards from cursor, including cursor [NOT SUPPORTED].
	Backward = ListDirection(-1)
	// Forward lists forwards from cursor, including cursor.
	Forward = ListDirection(1)
	// After lists forwards from cursor, without cursor.
	After = ListDirection(2)
)

// ListOptions lists objects.
type ListOptions struct {
	Prefix    Path
	Cursor    Path // Cursor is relative to Prefix, full path is Prefix + Cursor
	Delimiter rune
	Recursive bool
	Direction ListDirection
	Limit     int
	Status    int32
}

// BucketListOptions lists objects.
type BucketListOptions struct {
	Cursor    string
	Direction ListDirection
	Limit     int
}

// BucketList is a list of buckets.
type BucketList struct {
	More  bool
	Items []Bucket
}

// NextPage returns options for listing the next page.
func (opts BucketListOptions) NextPage(list BucketList) BucketListOptions {
	if !list.More || len(list.Items) == 0 {
		return BucketListOptions{}
	}

	return BucketListOptions{
		Cursor:    list.Items[len(list.Items)-1].Name,
		Direction: After,
		Limit:     opts.Limit,
	}
}
