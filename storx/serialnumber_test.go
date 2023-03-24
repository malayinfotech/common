// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"common/storx"
	"common/testrand"
)

func TestSerialNumber_Encode(t *testing.T) {
	_, err := storx.SerialNumberFromString("likn43kilfzd")
	assert.Error(t, err)

	_, err = storx.SerialNumberFromBytes([]byte{1, 2, 3, 4, 5})
	assert.Error(t, err)

	for i := 0; i < 10; i++ {
		serialNumber := testrand.SerialNumber()

		fromString, err := storx.SerialNumberFromString(serialNumber.String())
		assert.NoError(t, err)
		fromBytes, err := storx.SerialNumberFromBytes(serialNumber.Bytes())
		assert.NoError(t, err)

		assert.Equal(t, serialNumber, fromString)
		assert.Equal(t, serialNumber, fromBytes)
	}
}
