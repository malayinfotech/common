// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package pb

import (
	"encoding/base64"
	"encoding/json"

	"common/encryption"
	"common/storx"
)

type encryptionAccessStoreEntryMarshal struct {
	Bucket               string                `json:"bucket,omitempty"`
	UnencryptedPath      string                `json:"unencrypted_path,omitempty"`
	EncryptedPath        string                `json:"encrypted_path,omitempty"`
	Key                  string                `json:"key,omitempty"`
	PathCipher           CipherSuite           `json:"path_cipher,omitempty"`
	EncryptionParameters *EncryptionParameters `json:"encryption_parameters,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (se *EncryptionAccess_StoreEntry) MarshalJSON() ([]byte, error) {
	key, err := storx.NewKey([]byte{})
	if err != nil {
		return nil, err
	}

	path, err := encryption.DecryptPathRaw(string(se.EncryptedPath), storx.EncNullBase64URL, key)
	if err != nil {
		return nil, err
	}

	return json.Marshal(encryptionAccessStoreEntryMarshal{
		Bucket:               string(se.Bucket),
		UnencryptedPath:      string(se.UnencryptedPath),
		EncryptedPath:        path,
		Key:                  base64.URLEncoding.EncodeToString(se.Key),
		PathCipher:           se.PathCipher,
		EncryptionParameters: se.EncryptionParameters,
	})
}
