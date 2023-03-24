// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package encryption_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/encryption"
	"common/memory"
	"common/storx"
	"common/testcontext"
	"common/testrand"
)

const (
	uint32Size = 4
)

func TestCalcEncryptedSize(t *testing.T) {
	_ = testcontext.New(t)

	forAllCiphers(func(cipher storx.CipherSuite) {
		for i, dataSize := range []int64{
			0,
			1,
			1*memory.KiB.Int64() - uint32Size,
			1 * memory.KiB.Int64(),
			32*memory.KiB.Int64() - uint32Size,
			32 * memory.KiB.Int64(),
			32*memory.KiB.Int64() + 100,
		} {
			errTag := fmt.Sprintf("%d-%d. %+v", cipher, i, dataSize)

			parameters := storx.EncryptionParameters{CipherSuite: cipher, BlockSize: 1 * memory.KiB.Int32()}

			calculatedSize, err := encryption.CalcEncryptedSize(dataSize, parameters)
			require.NoError(t, err, errTag)

			encrypter, err := encryption.NewEncrypter(parameters.CipherSuite, new(storx.Key), new(storx.Nonce), int(parameters.BlockSize))
			require.NoError(t, err, errTag)

			randReader := io.NopCloser(io.LimitReader(testrand.Reader(), dataSize))
			reader := encryption.TransformReader(encryption.PadReader(randReader, encrypter.InBlockSize()), encrypter, 0)

			cipherData, err := io.ReadAll(reader)
			assert.NoError(t, err, errTag)
			assert.EqualValues(t, calculatedSize, len(cipherData), errTag)
		}
	})
}

func forAllCiphers(test func(cipher storx.CipherSuite)) {
	for _, cipher := range []storx.CipherSuite{
		storx.EncNull,
		storx.EncAESGCM,
		storx.EncSecretBox,
	} {
		test(cipher)
	}
}
