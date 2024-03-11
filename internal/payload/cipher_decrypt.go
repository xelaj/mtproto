// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"crypto/aes"
	"fmt"

	"github.com/gotd/ige"
	"github.com/pkg/errors"
)

// Decrypt decrypts data from encrypted message using AES-IGE.
// func (c Cipher) Decrypt(k Key, msgKey Int128, encrypted []byte) (*Envelope, error) {
func (c *cipher) Decrypt(data []byte) ([]byte, error) {
	const minLenEncrypted = 8 + 16
	if len(data) < minLenEncrypted {
		return nil, fmt.Errorf("message is too small: got %v, want %v", len(data), minLenEncrypted)
	}

	gotKeyID := [8]byte(data[:8])
	msgKey := [16]byte(data[8 : 8+16])
	encrypted := data[8+16:]

	if kid := c.key.ID(); gotKeyID != kid {
		return nil, fmt.Errorf("wrong key id: got %v, want %v", gotKeyID, kid)
	}

	plaintext, err := c.decryptMessage(msgKey, encrypted)
	if err != nil {
		return nil, err
	}

	// Checking SHA256 hash value of msg_key
	if msgKey != MessageKey(&c.key, plaintext, c.side.Invert()) {
		return nil, errors.New("msg_key is invalid")
	}

	return plaintext, nil
}

// decryptMessage decrypts data from encrypted message using AES-IGE.
func (c *cipher) decryptMessage(msgKey Int128, encrypted []byte) ([]byte, error) {
	if len(encrypted)%16 != 0 {
		return nil, errors.New("invalid encrypted data padding")
	}

	key, iv := Keys(&c.key, msgKey, c.side.Invert())
	cipher, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(encrypted))
	ige.DecryptBlocks(cipher, iv[:], plaintext, encrypted)

	return plaintext, nil
}
