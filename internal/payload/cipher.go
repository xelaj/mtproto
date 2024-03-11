// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import "io"

type Cipher interface {
	Encrypt(keyID [8]byte, content []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type cipher struct {
	rand io.Reader

	// side describes who created Cipher object. For example, if you are
	// implementing mtproto client — use ClientSide. If you are a server, which
	// handles requests — ClientServer is a right value
	side Side

	key Key
}

var _ Cipher = (*cipher)(nil)

// NewClientCipher creates new client-side Cipher.
//
//nolint:gocritic // hugeParam, doing it once
func NewClientCipher(rand io.Reader, key Key) Cipher { return newCipher(rand, &key, SideClient) }

// NewServerCipher creates new server-side Cipher.
//
//nolint:gocritic // hugeParam, doing it once
func NewServerCipher(rand io.Reader, key Key) Cipher { return newCipher(rand, &key, SideServer) }

var emptyKey Key //nolint:gochecknoglobals // it's cheaper

func newCipher(rand io.Reader, key *Key, side Side) Cipher {
	if rand == nil {
		rand = nopReader{}
	}
	if *key == emptyKey {
		panic("key is empty!")
	}

	return &cipher{rand: rand, side: side, key: *key}
}

type nopReader struct{}

func (nopReader) Read(p []byte) (n int, err error) { clear(p); return len(p), nil }
