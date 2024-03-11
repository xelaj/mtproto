// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testAuthKey = Key{ //nolint:gochecknoglobals // for testing
	93, 46, 125, 101, 244, 158, 194, 139, 208, 41, 168, 135, 97, 234, 39, 184, 164, 199,
	159, 18, 34, 101, 37, 68, 62, 125, 124, 89, 110, 243, 48, 53, 48, 219, 33, 7, 232, 154,
	169, 151, 199, 160, 22, 74, 182, 148, 24, 122, 222, 255, 21, 107, 214, 239, 113, 24,
	161, 150, 35, 71, 117, 60, 14, 126, 137, 160, 53, 75, 142, 195, 100, 249, 153, 126,
	113, 188, 105, 35, 251, 134, 232, 228, 52, 145, 224, 16, 96, 106, 108, 232, 69, 226,
	250, 1, 148, 9, 119, 239, 10, 163, 42, 223, 90, 151, 219, 246, 212, 40, 236, 4, 52,
	215, 23, 162, 211, 173, 25, 98, 44, 192, 88, 135, 100, 33, 19, 199, 150, 95, 251, 134,
	42, 62, 60, 203, 10, 185, 90, 221, 218, 87, 248, 146, 69, 219, 215, 107, 73, 35, 72,
	248, 233, 75, 213, 167, 192, 224, 184, 72, 8, 82, 60, 253, 30, 168, 11, 50, 254, 154,
	209, 152, 188, 46, 16, 63, 206, 183, 213, 36, 146, 236, 192, 39, 58, 40, 103, 75, 201,
	35, 238, 229, 146, 101, 171, 23, 160, 2, 223, 31, 74, 162, 197, 155, 129, 154, 94, 94,
	29, 16, 94, 193, 23, 51, 111, 92, 118, 198, 177, 135, 3, 125, 75, 66, 112, 206, 233,
	204, 33, 7, 29, 151, 233, 188, 162, 32, 198, 215, 176, 27, 153, 140, 242, 229, 205,
	185, 165, 14, 205, 161, 133, 42, 54, 230, 53, 105, 12, 142,
}

func checkSame(t *testing.T, a, b Cipher) {
	asserts := require.New(t)

	sessionID, err := rand.Int(rand.Reader, big.NewInt(2345512351))
	asserts.NoError(err)

	msg := []byte("data")
	buf, err := a.Encrypt([8]byte{}, BuildEnvelope(0, uint64(sessionID.Int64()), 0, 0, msg, rand.Reader))
	asserts.NoError(err)

	decrypt, err := b.Decrypt(buf)
	asserts.NoError(err)

	e, err := DeserializeEnvelope(decrypt)
	asserts.NoError(err)

	asserts.Equal(sessionID.Int64(), int64(e.SessionID))
	asserts.Equal(msg, e.Msg)
}

func TestCipher(t *testing.T) {
	client := NewClientCipher(rand.Reader, testAuthKey)
	server := NewServerCipher(rand.Reader, testAuthKey)

	checkSame(t, client, server)
	checkSame(t, server, client)
}

func TestEncrypt(t *testing.T) {
	var authKey Key
	for i := 0; i < 256; i++ {
		authKey[i] = byte(i)
	}

	c := NewClientCipher(nopReader{}, authKey)

	// Testing vector from grammers.
	got, err := c.Encrypt(
		[8]byte{},
		[]byte("Hello, world! This data should remain secure!"),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		50, 209, 88, 110, 164, 87, 223, 200, 168, 23, 41, 212, 109, 181, 64, 25, 162, 191, 215,
		247, 68, 249, 185, 108, 79, 113, 108, 253, 196, 71, 125, 178, 162, 193, 95, 109, 219,
		133, 35, 95, 185, 85, 47, 29, 132, 7, 198, 170, 234, 0, 204, 132, 76, 90, 27, 246, 172,
		68, 183, 155, 94, 220, 42, 35, 134, 139, 61, 96, 115, 165, 144, 153, 44, 15, 41, 117,
		36, 61, 86, 62, 161, 128, 210, 24, 238, 117, 124, 154,
	}
	require.Equal(t, want, got)
}

func TestDecrypt(t *testing.T) {
	// Test vector from grammers.
	c := NewClientCipher(nopReader{}, testAuthKey)
	b := []byte{
		122, 113, 131, 194, 193, 14, 79, 77, 249, 69, 250, 154, 154, 189, 53, 231, 195, 132,
		11, 97, 240, 69, 48, 79, 57, 103, 76, 25, 192, 226, 9, 120, 79, 80, 246, 34, 106, 7,
		53, 41, 214, 117, 201, 44, 191, 11, 250, 140, 153, 167, 155, 63, 57, 199, 42, 93, 154,
		2, 109, 67, 26, 183, 64, 124, 160, 78, 204, 85, 24, 125, 108, 69, 241, 120, 113, 82,
		78, 221, 144, 206, 160, 46, 215, 40, 225, 77, 124, 177, 138, 234, 42, 99, 97, 88, 240,
		148, 89, 169, 67, 119, 16, 216, 148, 199, 159, 54, 140, 78, 129, 100, 183, 100, 126,
		169, 134, 18, 174, 254, 148, 44, 93, 146, 18, 26, 203, 141, 176, 45, 204, 206, 182,
		109, 15, 135, 32, 172, 18, 160, 109, 176, 88, 43, 253, 149, 91, 227, 79, 54, 81, 24,
		227, 186, 184, 205, 8, 12, 230, 180, 91, 40, 234, 197, 109, 205, 42, 41, 55, 78,
	}

	plaintext, err := c.Decrypt(b)
	if err != nil {
		t.Fatal(err)
	}
	expectedPlaintext := []byte{
		252, 130, 106, 2, 36, 139, 40, 253, 96, 242, 196, 130, 36, 67, 173, 104, 1, 240, 193,
		194, 145, 139, 48, 94, 2, 0, 0, 0, 88, 0, 0, 0, 220, 248, 241, 115, 2, 0, 0, 0, 1, 168,
		193, 194, 145, 139, 48, 94, 1, 0, 0, 0, 28, 0, 0, 0, 8, 9, 194, 158, 196, 253, 51, 173,
		145, 139, 48, 94, 24, 168, 142, 166, 7, 238, 88, 22, 252, 130, 106, 2, 36, 139, 40,
		253, 1, 204, 193, 194, 145, 139, 48, 94, 2, 0, 0, 0, 20, 0, 0, 0, 197, 115, 119, 52,
		196, 253, 51, 173, 145, 139, 48, 94, 100, 8, 48, 0, 0, 0, 0, 0, 252, 230, 103, 4, 163,
		205, 142, 233, 208, 174, 111, 171, 103, 44, 96, 192, 74, 63, 31, 212, 73, 14, 81, 246,
	}
	if !bytes.Equal(expectedPlaintext, plaintext) {
		t.Error("mismatch")
	}
}

func TestCipher_Decrypt(t *testing.T) {
	var key Key
	if _, err := io.ReadFull(predictedRand([]byte{10}), key[:]); err != nil {
		t.Fatal(err)
	}

	c := NewClientCipher(nopReader{}, key)
	s := NewServerCipher(nopReader{}, key)
	tests := []struct {
		name    string
		data    []byte
		dataLen int
		wantErr assert.ErrorAssertionFunc
	}{
		{"NegativeLength", []byte{1, 2, 3, 4}, -1, assert.Error},
		{"NoPadBy4", []byte{1, 2, 3}, 3, assert.Error},
		{"Good", bytes.Repeat([]byte{1, 2, 3, 4}, 4), 16, assert.NoError},
	}

	for _, test := range tests {
		test.wantErr = noErrAsDefault(test.wantErr)

		t.Run(test.name, func(t *testing.T) {
			a := require.New(t)
			e := BuildEnvelope(0, 0, 0, 0, test.data, nil)
			encrypted, err := s.Encrypt(key.ID(), e)
			if !test.wantErr(t, err) || err != nil {
				return
			}

			decrypted, err := c.Decrypt(encrypted)
			if !test.wantErr(t, err) || err != nil {
				return
			}

			a.Equal(e, decrypted)
		})
	}
}
