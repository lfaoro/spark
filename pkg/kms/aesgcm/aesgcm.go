// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"

	"github.com/lfaoro/spark/pkg/kms"
)

type AESGCM struct {
	key string
}

func New(key string) (kms.Service, error) {
	if key == "" {
		key = "oUnmr6cm9G7yeRBMsZrqk84w8M22jcxM"
	}
	if len(key) < 32 {
		return nil, errors.Errorf("aesgcm: key size must be 32 characters, currently: %v", len(key))
	}
	return AESGCM{key: key}, nil
}

func (a AESGCM) Encrypt(plaintext []byte) ([]byte, error) {
	_cipher, err := aes.NewCipher([]byte(a.key))
	if err != nil {
		return nil, errors.Wrap(err, "aesgcm: failed to create a new cipher")
	}

	gcm, err := cipher.NewGCM(_cipher)
	if err != nil {
		return nil, errors.Wrap(err, "aesgcm: failed to wrap cipher in GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "aesgcm: failed to read random nonce")
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (a AESGCM) Decrypt(ciphertext []byte) ([]byte, error) {
	_cipher, err := aes.NewCipher([]byte(a.key))
	if err != nil {
		return nil, errors.Wrap(err, "aesgcm: unable to create new cipher")
	}

	gcm, err := cipher.NewGCM(_cipher)
	if err != nil {
		return nil, errors.Wrap(err, "aesgcm: unable to wrap cipher in GCM")
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.Wrap(err, "aesgcm: unable to read random nonce")
	}

	nonce, plaintext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, plaintext, nil)
}
