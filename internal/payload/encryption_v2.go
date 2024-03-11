// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import "crypto/sha256"

// MessageKey computes message key for provided auth_key and padded payload.
func MessageKey(authKey *Key, plaintextPadded []byte, mode Side) Int128 {
	// `msg_key_large = SHA256 (substr (auth_key, 88+x, 32) + plaintext + random_padding);`
	msgKeyLarge := msgKeyLarge(authKey, plaintextPadded, mode)
	// `msg_key = substr (msg_key_large, 8, 16);`
	return messageKey(msgKeyLarge)
}

// Keys returns (aes_key, aes_iv) pair for AES-IGE.
//
// See https://core.telegram.org/mtproto/description#defining-aes-key-and-initialization-vector
//
// Example:
//
//	key, iv := crypto.Keys(authKey, messageKey, crypto.Client)
//	cipher, err := aes.NewCipher(key[:])
//	if err != nil {
//		return nil, err
//	}
//	encryptor := ige.NewIGEEncrypter(cipher, iv[:])
func Keys(authKey *Key, msgKey Int128, mode Side) (key, iv Int256) {
	x := mode.X()

	// `sha256_a = SHA256 (msg_key + substr (auth_key, x, 36));`
	a := sha256a(authKey, msgKey, x)
	// `sha256_b = SHA256 (substr (auth_key, 40+x, 36) + msg_key);`
	b := sha256b(authKey, msgKey, x)

	return aesKey(a, b), aesIV(a, b)
}

// msg_key = substr (msg_key_large, 8, 16).
func messageKey(messageKeyLarge Int256) (v Int128) {
	b := messageKeyLarge[8 : 16+8]
	copy(v[:len(b)], b)
	return v
}

// aesKey writes aes_key value into v.
//
// aes_key = substr (sha256_a, 0, 8) + substr (sha256_b, 8, 16) + substr (sha256_a, 24, 8);
func aesKey(sha256a, sha256b Int256) (key Int256) {
	copy(key[:8], sha256a[:8])
	copy(key[8:], sha256b[8:16+8])
	copy(key[24:], sha256a[24:24+8])

	return key
}

// aesIV writes aes_iv value into v.
//
// aes_iv = substr (sha256_b, 0, 8) + substr (sha256_a, 8, 16) + substr (sha256_b, 24, 8);
func aesIV(sha256a, sha256b Int256) (iv Int256) {
	// Same as aes_key, but with swapped params.
	return aesKey(sha256b, sha256a)
}

// sha256a returns sha256_a value.
//
// sha256_a = SHA256 (msg_key + substr (auth_key, x, 36));
func sha256a(authKey *Key, msgKey Int128, x int) (sum Int256) {
	h := sha256.New()
	_, _ = h.Write(msgKey[:])
	_, _ = h.Write(authKey[x : x+36])
	copy(sum[:], h.Sum(nil))

	return sum
}

// sha256b returns sha256_b value.
//
// sha256_b = SHA256 (substr (auth_key, 40+x, 36) + msg_key);
func sha256b(authKey *Key, msgKey Int128, x int) (sum Int256) {
	h := sha256.New()
	_, _ = h.Write(authKey[40+x : 40+x+36])
	_, _ = h.Write(msgKey[:])
	copy(sum[:], h.Sum(nil))

	return sum
}

// msgKeyLarge returns msg_key_large value.
func msgKeyLarge(authKey *Key, plaintextPadded []byte, mode Side) (sum Int256) {
	h := sha256.New()

	x := mode.X()
	_, _ = h.Write(authKey[88+x : 32+88+x])
	_, _ = h.Write(plaintextPadded)

	copy(sum[:], h.Sum(nil))

	return sum
}
