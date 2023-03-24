// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"common/internal/hmacsha512"
	"common/storx"
)

const (
	// AESGCMNonceSize is the size of an AES-GCM nonce.
	AESGCMNonceSize = 12
	// unit32Size is the number of bytes in the uint32 type.
	uint32Size = 4
)

// AESGCMNonce represents the nonce used by the AES-GCM protocol.
type AESGCMNonce [AESGCMNonceSize]byte

// ToAESGCMNonce returns the nonce as a AES-GCM nonce.
func ToAESGCMNonce(nonce *storx.Nonce) *AESGCMNonce {
	aes := new(AESGCMNonce)
	copy((*aes)[:], nonce[:AESGCMNonceSize])
	return aes
}

// Increment increments the nonce with the given amount.
func Increment(nonce *storx.Nonce, amount int64) (truncated bool, err error) {
	return incrementBytes(nonce[:], amount)
}

// Encrypt encrypts data with the given cipher, key and nonce.
func Encrypt(data []byte, cipher storx.CipherSuite, key *storx.Key, nonce *storx.Nonce) (cipherData []byte, err error) {
	// Don't encrypt empty slice
	if len(data) == 0 {
		return []byte{}, nil
	}

	switch cipher {
	case storx.EncNull:
		return data, nil
	case storx.EncAESGCM:
		return EncryptAESGCM(data, key, ToAESGCMNonce(nonce))
	case storx.EncSecretBox:
		return EncryptSecretBox(data, key, nonce)
	case storx.EncNullBase64URL:
		return nil, ErrInvalidConfig.New("base64 encoding not supported for this operation")
	default:
		return nil, ErrInvalidConfig.New("encryption type %d is not supported", cipher)
	}
}

// Decrypt decrypts cipherData with the given cipher, key and nonce.
func Decrypt(cipherData []byte, cipher storx.CipherSuite, key *storx.Key, nonce *storx.Nonce) (data []byte, err error) {
	// Don't decrypt empty slice
	if len(cipherData) == 0 {
		return []byte{}, nil
	}

	switch cipher {
	case storx.EncNull:
		return cipherData, nil
	case storx.EncAESGCM:
		return DecryptAESGCM(cipherData, key, ToAESGCMNonce(nonce))
	case storx.EncSecretBox:
		return DecryptSecretBox(cipherData, key, nonce)
	case storx.EncNullBase64URL:
		return nil, ErrInvalidConfig.New("base64 encoding not supported for this operation")
	default:
		return nil, ErrInvalidConfig.New("encryption type %d is not supported", cipher)
	}
}

// NewEncrypter creates a Transformer using the given cipher, key and nonce to encrypt data passing through it.
func NewEncrypter(cipher storx.CipherSuite, key *storx.Key, startingNonce *storx.Nonce, encryptedBlockSize int) (Transformer, error) {
	switch cipher {
	case storx.EncNull:
		return &NoopTransformer{}, nil
	case storx.EncAESGCM:
		return NewAESGCMEncrypter(key, ToAESGCMNonce(startingNonce), encryptedBlockSize)
	case storx.EncSecretBox:
		return NewSecretboxEncrypter(key, startingNonce, encryptedBlockSize)
	case storx.EncNullBase64URL:
		return nil, ErrInvalidConfig.New("base64 encoding not supported for this operation")
	default:
		return nil, ErrInvalidConfig.New("encryption type %d is not supported", cipher)
	}
}

// NewDecrypter creates a Transformer using the given cipher, key and nonce to decrypt data passing through it.
func NewDecrypter(cipher storx.CipherSuite, key *storx.Key, startingNonce *storx.Nonce, encryptedBlockSize int) (Transformer, error) {
	switch cipher {
	case storx.EncNull:
		return &NoopTransformer{}, nil
	case storx.EncAESGCM:
		return NewAESGCMDecrypter(key, ToAESGCMNonce(startingNonce), encryptedBlockSize)
	case storx.EncSecretBox:
		return NewSecretboxDecrypter(key, startingNonce, encryptedBlockSize)
	case storx.EncNullBase64URL:
		return nil, ErrInvalidConfig.New("base64 encoding not supported for this operation")
	default:
		return nil, ErrInvalidConfig.New("encryption type %d is not supported", cipher)
	}
}

// EncryptKey encrypts keyToEncrypt with the given cipher, key and nonce.
func EncryptKey(keyToEncrypt *storx.Key, cipher storx.CipherSuite, key *storx.Key, nonce *storx.Nonce) (storx.EncryptedPrivateKey, error) {
	return Encrypt(keyToEncrypt[:], cipher, key, nonce)
}

// DecryptKey decrypts keyToDecrypt with the given cipher, key and nonce.
func DecryptKey(keyToDecrypt storx.EncryptedPrivateKey, cipher storx.CipherSuite, key *storx.Key, nonce *storx.Nonce) (*storx.Key, error) {
	plainData, err := Decrypt(keyToDecrypt, cipher, key, nonce)
	if err != nil {
		return nil, err
	}

	var decryptedKey storx.Key
	copy(decryptedKey[:], plainData)

	return &decryptedKey, nil
}

// DeriveKey derives new key from the given key and message using HMAC-SHA512.
func DeriveKey(key *storx.Key, message string) (*storx.Key, error) {
	mac := hmacsha512.New(key[:])
	mac.Write([]byte(message))
	derived := new(storx.Key)
	sum := mac.SumAndReset()
	copy(derived[:], sum[:])
	return derived, nil
}

// CalcEncryptedSize calculates what would be the size of the cipher data after
// encrypting data with dataSize using a Transformer with the given encryption
// parameters.
func CalcEncryptedSize(dataSize int64, parameters storx.EncryptionParameters) (int64, error) {
	transformer, err := NewEncrypter(parameters.CipherSuite, new(storx.Key), new(storx.Nonce), int(parameters.BlockSize))
	if err != nil {
		return 0, err
	}
	return CalcTransformerEncryptedSize(dataSize, transformer), nil
}

// CalcTransformerEncryptedSize calculates what would be the size of the
// cipher data after encrypting data with dataSize using the given Transformer.
func CalcTransformerEncryptedSize(dataSize int64, transformer Transformer) int64 {
	inBlockSize := int64(transformer.InBlockSize())
	blocks := (dataSize + uint32Size + inBlockSize - 1) / inBlockSize
	encryptedSize := blocks * int64(transformer.OutBlockSize())
	return encryptedSize
}
