// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"crypto/aes"
	"errors"
	"fmt"

	"github.com/gotd/ige"
)

// default key id is empty key id
var defaultKeyID [8]byte //nolint:gochecknoglobals // allowed

// Encrypt encrypts data using AES-IGE to given buffer.
func (c *cipher) Encrypt(kid [8]byte, data []byte) ([]byte, error) {
	if kid == defaultKeyID {
		kid = c.key.ID()
	} else if kid != c.key.ID() {
		return nil, errors.New("unknown key id")
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("input data is not full AES block")
	}

	messageKey := MessageKey(&c.key, data, c.side)
	key, iv := Keys(&c.key, messageKey, c.side)
	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	encrypted := make([]byte, len(data))
	ige.EncryptBlocks(aesBlock, iv[:], encrypted, data)

	res := make([]byte, 8+16+len(encrypted))
	copy(res[:8], kid[:])
	copy(res[8:8+16], messageKey[:])
	copy(res[8+16:], encrypted)

	return res, nil
}
