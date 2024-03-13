package handshake

import (
	"bytes"
	"crypto/aes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/gotd/ige"
)

// tempAESKeys returns tmp_aes_key and tmp_aes_iv based on new_nonce and
// server_nonce as defined in "Creating an Authorization Key".
//
// See https://core.telegram.org/mtproto/auth_key#6-server-responds-with
//
// tmp_aes_key := SHA1(new_nonce + server_nonce) + substr (SHA1(server_nonce + new_nonce), 0, 12);
func TempAESKeys(newNonce Int256, serverNonce Int128) (key, iv Int256) {
	// n is newNonce, s is serverNonce
	ns := sha1.Sum(append(newNonce[:], serverNonce[:]...))
	sn := sha1.Sum(append(serverNonce[:], newNonce[:]...))
	nn := sha1.Sum(append(newNonce[:], newNonce[:]...))

	return tempAESKey(ns, sn), tempAESIV(sn, nn, newNonce)
}

// tmp_aes_key := SHA1(new_nonce + server_nonce) + substr (SHA1(server_nonce + new_nonce), 0, 12);
func tempAESKey(ns, sn [sha1.Size]byte) (key Int256) {
	// SHA1(new_nonce + server_nonce)
	copy(key[:20], ns[:])

	// substr(SHA1(server_nonce + new_nonce), 0, 12);
	copy(key[20:32], sn[:12])

	return key
}

// tmp_aes_iv := substr(SHA1(server_nonce + new_nonce), 12, 8) + SHA1(new_nonce + new_nonce) + substr (new_nonce, 0, 4);
func tempAESIV(sn, nn [sha1.Size]byte, newNonce Int256) (iv Int256) {
	// substr(SHA1(server_nonce + new_nonce), 12, 8)
	copy(iv[0:8], sn[12:12+8])

	// SHA1(new_nonce + new_nonce)
	copy(iv[8:8+20], nn[:])

	// substr(new_nonce, 0, 4)
	copy(iv[8+20:8+20+4], newNonce[:4])

	return iv
}

func encryptHandshake(rand io.Reader, key, iv Int256, data []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}

	answerWithHash := DataWithHash(data, rand)
	dst := make([]byte, len(answerWithHash))
	ige.EncryptBlocks(cipher, iv[:], dst, answerWithHash)

	return dst, nil
}

func DataWithHash(data []byte, rand io.Reader) []byte {
	dataWithHashAndPad := make([]byte, len(data)+sha1.Size+pad(len(data)+sha1.Size, 16))
	if _, err := io.ReadFull(rand, dataWithHashAndPad); err != nil {
		panic(err)
	}

	h := sha1.Sum(data)
	copy(dataWithHashAndPad, h[:])
	copy(dataWithHashAndPad[sha1.Size:], data)

	return dataWithHashAndPad
}

func pad(l, n int) int { return (n - l%n) % n }

func decryptHandshake(key, iv Int256, data []byte) ([]byte, error) {
	if data == nil {
		// Most common cause of this error is invalid crypto implementation,
		// i.e. invalid keys are used to decrypt payload which lead to
		// decrypt failure, so data does not match sha1 with any padding.
		return nil, errors.New("guess data from data_with_hash")
	}

	cipher, err := aes.NewCipher(key[:])
	if err != nil {
		panic(fmt.Errorf("create aes cipher: %w", err))
	}

	dataWithHash := make([]byte, len(data))
	if len(dataWithHash)%cipher.BlockSize() != 0 {
		return nil, fmt.Errorf("invalid len of data_with_hash (%d %% 16 != 0)", len(dataWithHash))
	}
	ige.DecryptBlocks(cipher, iv[:], dataWithHash, data)

	return GuessDataWithHash(dataWithHash), nil
}

// guessDataWithHash guesses data from data_with_hash.
func GuessDataWithHash(dataWithHash []byte) []byte {
	// data_with_hash := SHA1(data) + data + (0-15 random bytes);
	// such that length be divisible by 16;
	if len(dataWithHash) <= sha1.Size {
		return nil // Data length too small.
	}

	v := dataWithHash[:sha1.Size]
	for i := 0; i < 16; i++ {
		if len(dataWithHash)-i < sha1.Size {

			return nil // End of slice reached.
		}
		data := dataWithHash[sha1.Size : len(dataWithHash)-i]
		h := sha1.Sum(data)
		if bytes.Equal(h[:], v) {
			return data // Found.
		}
	}
	return nil
}
