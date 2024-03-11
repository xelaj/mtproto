package handshake

import (
	"bytes"
	"crypto/aes"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"slices"

	"github.com/go-faster/xor"
	"github.com/gotd/ige"
	"github.com/xelaj/tl"
)

const (
	rsaLen                = 256
	rsaPadDataLimit       = 144
	dataWithPaddingLength = 192
	dataWithHashLength    = dataWithPaddingLength + sha256.Size
	tempKeySize           = 32
)

// RSAPad encrypts given data with RSA, prefixing with a hash.
//
// See https://core.telegram.org/mtproto/auth_key#presenting-proof-of-work-server-authentication.
func RSAPad(data []byte, key *rsa.PublicKey, randomSource io.Reader) ([]byte, error) {
	// 1) data_with_padding := data + random_padding_bytes; — where random_padding_bytes are
	// chosen so that the resulting length of data_with_padding is precisely 192 bytes, and
	// data is the TL-serialized data to be encrypted as before.
	//
	// One has to check that data is not longer than 144 bytes.
	if len(data) > rsaPadDataLimit {
		return nil, fmt.Errorf("data length is bigger that 144 (%d)", len(data))
	}

	dataWithPadding := make([]byte, dataWithPaddingLength)
	// Filling data_with_padding with random bytes firstly. It's easier to do it
	// in first place for testing purposes,in result there wil be no difference.
	if _, err := io.ReadFull(randomSource, dataWithPadding); err != nil {
		return nil, fmt.Errorf("pad data with random: %w", err)
	}

	// then copying data to start of data_with_padding.
	copy(dataWithPadding, data)

	// Make a copy.
	dataPadReversed := make([]byte, dataWithPaddingLength)
	copy(dataPadReversed, dataWithPadding)
	// 2) data_pad_reversed := BYTE_REVERSE(data_with_padding);
	slices.Reverse(dataPadReversed)

	for {
		// 3) A random 32-byte temp_key is generated.
		tempKey := make([]byte, tempKeySize)
		if _, err := io.ReadFull(randomSource, tempKey); err != nil {
			return nil, fmt.Errorf("generate temp_key: %w", err)
		}

		// 4) data_with_hash := data_pad_reversed + SHA256(temp_key + data_with_padding);
		// — after this assignment, data_with_hash is exactly 224 bytes long.
		dataWithHash := make([]byte, 0, dataWithHashLength)
		dataWithHash = append(dataWithHash, dataPadReversed...)
		{
			h := sha256.New()
			_, _ = h.Write(tempKey)
			_, _ = h.Write(dataWithPadding)
			dataWithHash = h.Sum(dataWithHash)
			dataWithHash = dataWithHash[:dataWithHashLength]
		}

		// 5) aes_encrypted := AES256_IGE(data_with_hash, temp_key, 0); — AES256-IGE encryption with zero IV.
		aesEncrypted := make([]byte, len(dataWithHash))
		{
			aesBlock, err := aes.NewCipher(tempKey)
			if err != nil {
				return nil, fmt.Errorf("create cipher: %w", err)
			}
			var zeroIV [32]byte
			ige.EncryptBlocks(aesBlock, zeroIV[:], aesEncrypted, dataWithHash)
		}

		// 6) temp_key_xor := temp_key XOR SHA256(aes_encrypted); — adjusted key, 32 bytes
		tempKeyXor := make([]byte, tempKeySize)
		{
			aesEncryptedHash := sha256.Sum256(aesEncrypted)
			xor.Bytes(tempKeyXor, tempKey, aesEncryptedHash[:])
		}

		// 7) key_aes_encrypted := temp_key_xor + aes_encrypted; — exactly 256 bytes (2048 bits) long.
		keyAESEncrypted := make([]byte, 0, tempKeySize+dataWithHashLength)
		keyAESEncrypted = append(keyAESEncrypted, tempKeyXor...)
		keyAESEncrypted = append(keyAESEncrypted, aesEncrypted...)

		// 8) The value of key_aes_encrypted is compared with the RSA-modulus of server_pubkey
		// as a big-endian 2048-bit (256-byte) unsigned integer. If key_aes_encrypted turns out to be
		// greater than or equal to the RSA modulus, the previous steps starting from the generation
		// of new random temp_key are repeated.
		keyAESEncryptedBig := big.NewInt(0).SetBytes(keyAESEncrypted)
		if keyAESEncryptedBig.Cmp(key.N) >= 0 {
			continue
		}
		// Otherwise the final step is performed:

		// 9) encrypted_data := RSA(key_aes_encrypted, server_pubkey);
		// — 256-byte big-endian integer is elevated to the requisite power from the RSA public key
		// modulo the RSA modulus, and the result is stored as a big-endian integer consisting of
		// exactly 256 bytes (with leading zero bytes if required).
		//
		// Encrypting "key_aes_encrypted" with RSA.
		res := rsaEncrypt(keyAESEncrypted, key)
		return res, nil
	}
}

// DecodeRSAPad implements server-side decoder of RSAPad.
func DecodeRSAPad(data []byte, key *rsa.PrivateKey) ([]byte, error) {
	var encryptedData [256]byte
	if !rsaDecrypt(data, key, encryptedData[:]) {
		return nil, errors.New("invalid encrypted_data")
	}

	tempKeyXor := encryptedData[:tempKeySize]
	aesEncrypted := encryptedData[tempKeySize:]

	tempKey := make([]byte, tempKeySize)
	{
		aesEncryptedHash := sha256.Sum256(aesEncrypted)
		xor.Bytes(tempKey, tempKeyXor, aesEncryptedHash[:])
	}

	dataWithHash := make([]byte, len(aesEncrypted))
	{
		aesBlock, err := aes.NewCipher(tempKey)
		if err != nil {
			return nil, fmt.Errorf("create cipher: %w", err)
		}
		var zeroIV [32]byte
		ige.DecryptBlocks(aesBlock, zeroIV[:], dataWithHash, aesEncrypted)
	}

	dataWithPadding := dataWithHash[:dataWithPaddingLength]
	slices.Reverse(dataWithPadding)

	hash := dataWithHash[dataWithPaddingLength:]
	{
		h := sha256.New()
		h.Write(tempKey)
		h.Write(dataWithPadding)

		if !bytes.Equal(hash, h.Sum(nil)) {
			return nil, errors.New("hash mismatch")
		}
	}

	return dataWithPadding, nil
}

func rsaEncrypt(data []byte, key *rsa.PublicKey) []byte {
	z := new(big.Int).SetBytes(data)
	e := big.NewInt(int64(key.E))
	c := new(big.Int).Exp(z, e, key.N)
	res := make([]byte, rsaLen)
	c.FillBytes(res)
	return res
}

func rsaDecrypt(data []byte, key *rsa.PrivateKey, to []byte) bool {
	c := new(big.Int).SetBytes(data)
	m := new(big.Int).Exp(c, key.D, key.N)
	return FillBytes(m, to)
}

// FillBytes is safe version of (*big.Int).FillBytes.
// Returns false if to length is not exact equal to big.Int's.
// Otherwise fills to using b and returns true.
func FillBytes(b *big.Int, to []byte) bool {
	bits := b.BitLen()
	if (bits+7)/8 > len(to) {
		return false
	}
	b.FillBytes(to)
	return true
}

func lookupForKeys(keys []*rsa.PublicKey, fingerprints []uint64) (index int) {
	return slices.IndexFunc(keys, func(pk *rsa.PublicKey) bool {
		return slices.Contains(fingerprints, rsaFingerprint(pk))
	})
}

// https://core.telegram.org/mtproto/auth_key#dh-exchange-initiation
// rsa_public_key n:string e:string = RSAPublicKey
func rsaFingerprint(key *rsa.PublicKey) uint64 {
	keyE := big.NewInt(int64(key.E))

	n, err := tl.Marshal(ptr(key.N.Bytes()))
	check(err)
	e, err := tl.Marshal(ptr(keyE.Bytes()))
	check(err)

	hash := sha1.Sum(append(n, e...))
	return binary.LittleEndian.Uint64(hash[12:20])
}

func ptr[T any](v T) *T { return &v }

func check(err error) {
	if err != nil {
		panic(err)
	}
}
