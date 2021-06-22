// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package ige

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
)

type Cipher struct {
	block   cipher.Block
	v       [3]AesBlock
	t, x, y []byte
}

// NewCipher
func NewCipher(key, iv []byte) (*Cipher, error) {
	const (
		firstBlock = iota
		secondBlock
		thirdBlock
	)

	var err error

	c := new(Cipher)
	c.block, err = aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "creating new cipher")
	}

	c.t = c.v[firstBlock][:]
	c.x = c.v[secondBlock][:]
	c.y = c.v[thirdBlock][:]
	copy(c.x, iv[:aes.BlockSize])
	copy(c.y, iv[aes.BlockSize:])

	return c, nil
}

func (c *Cipher) doAES256IGEencrypt(in, out []byte) error { //nolint:dupl потому что алгоритм на тоненького
	if err := isCorrectData(in); err != nil {
		return err
	}

	for i := 0; i < len(in); i += aes.BlockSize {
		xor(c.x, in[i:i+aes.BlockSize])
		c.block.Encrypt(c.t, c.x)
		xor(c.t, c.y)
		c.x, c.y = c.t, in[i:i+aes.BlockSize]
		copy(out[i:], c.t)
	}
	return nil
}

func (c *Cipher) doAES256IGEdecrypt(in, out []byte) error { //nolint:dupl потому что алгоритм на тоненького
	if err := isCorrectData(in); err != nil {
		return err
	}

	for i := 0; i < len(in); i += aes.BlockSize {
		xor(c.y, in[i:i+aes.BlockSize])
		c.block.Decrypt(c.t, c.y)
		xor(c.t, c.x)
		c.y, c.x = c.t, in[i:i+aes.BlockSize]
		copy(out[i:], c.t)
	}
	return nil
}

func isCorrectData(data []byte) error {
	if len(data) < aes.BlockSize {
		return ErrDataTooSmall
	}
	if len(data)%aes.BlockSize != 0 {
		return ErrDataNotDivisible
	}
	return nil
}

// --------------------------------------------------------------------------------------------------

// generateAESIGEv2 это переписанная функция generateAESIGE, которая выглядить чуточку более понятно.
func generateAESIGEv2(msgKey, authKey []byte, decode bool) (aesKey, aesIv []byte) { //nolint:deadcode wait for it
	var (
		kvBlock  [2]AesKV
		igeBlock [4]AesIgeBlock
	)

	aesKey = kvBlock[0][:]
	aesIv = kvBlock[1][:]

	tA := igeBlock[0][:]
	tB := igeBlock[1][:]
	tC := igeBlock[2][:]
	tD := igeBlock[3][:]

	var x int
	if decode {
		x = 8
	} else {
		x = 0
	}

	var (
		step       = 32
		tAOffStart = x
		tAOffEnd   = tAOffStart + step

		tBOffP0Start = tAOffEnd
		tBOffP0End   = tBOffP0Start + aes.BlockSize

		tBOffP1Start = tAOffEnd + aes.BlockSize
		tBOffP1End   = tBOffP1Start + aes.BlockSize

		tCOffStart = x + 64
		tCOffEnd   = tCOffStart + step

		tDOffStart = x + 96
		tDOffEnd   = tDOffStart + step
	)

	tA = append(tA, msgKey...)
	tA = append(tA, authKey[tAOffStart:tAOffEnd]...)

	tB = append(tB, authKey[tBOffP0Start:tBOffP0End]...)
	tB = append(tB, msgKey...)
	tB = append(tB, authKey[tBOffP1Start:tBOffP1End]...)

	tC = append(tC, authKey[tCOffStart:tCOffEnd]...)
	tC = append(tC, msgKey...)

	tD = append(tD, msgKey...)
	tD = append(tD, authKey[tDOffStart:tDOffEnd]...)

	sha1PartA := dry.Sha1Byte(tA)
	sha1PartB := dry.Sha1Byte(tB)
	sha1PartC := dry.Sha1Byte(tC)
	sha1PartD := dry.Sha1Byte(tD)

	aesKey = append(aesKey, sha1PartA[0:8]...)
	aesKey = append(aesKey, sha1PartB[8:8+12]...)
	aesKey = append(aesKey, sha1PartC[4:4+12]...)

	aesIv = append(aesIv, sha1PartA[8:8+12]...)
	aesIv = append(aesIv, sha1PartB[0:8]...)
	aesIv = append(aesIv, sha1PartC[16:16+4]...)
	aesIv = append(aesIv, sha1PartD[0:8]...)

	return aesKey, aesIv
}

// generateAESIGE ЭТО ЕБАНАЯ МАГИЧЕСКАЯ ФУНКЦИЯ ОНА НАХУЙ РАБОТАЕТ ПРОСТО БЛЯТЬ НЕ ТРОГАЙ ШАКАЛ ЕБАНЫЙ
//nolint:godox ты че ебанулся // TODO: порезать себе вены
func generateAESIGE(msg_key, auth_key []byte, decode bool) ([]byte, []byte) {
	var x int
	if decode {
		x = 8
	} else {
		x = 0
	}

	if len(auth_key) < 96+x+32 {
		panic(fmt.Sprintf("wrong len of auth key, got %v want at least %v", len(auth_key), 96+x+32))
	}

	aes_key := make([]byte, 0, 32)
	aes_iv := make([]byte, 0, 32)
	t_a := make([]byte, 0, 48)
	t_b := make([]byte, 0, 48)
	t_c := make([]byte, 0, 48)
	t_d := make([]byte, 0, 48)

	t_a = append(t_a, msg_key...)
	t_a = append(t_a, auth_key[x:x+32]...)

	t_b = append(t_b, auth_key[32+x:32+x+16]...)
	t_b = append(t_b, msg_key...)
	t_b = append(t_b, auth_key[48+x:48+x+16]...)

	t_c = append(t_c, auth_key[64+x:64+x+32]...)
	t_c = append(t_c, msg_key...)

	t_d = append(t_d, msg_key...)
	t_d = append(t_d, auth_key[96+x:96+x+32]...)

	sha1_a := dry.Sha1Byte(t_a)
	sha1_b := dry.Sha1Byte(t_b)
	sha1_c := dry.Sha1Byte(t_c)
	sha1_d := dry.Sha1Byte(t_d)

	aes_key = append(aes_key, sha1_a[0:8]...)
	aes_key = append(aes_key, sha1_b[8:8+12]...)
	aes_key = append(aes_key, sha1_c[4:4+12]...)

	aes_iv = append(aes_iv, sha1_a[8:8+12]...)
	aes_iv = append(aes_iv, sha1_b[0:8]...)
	aes_iv = append(aes_iv, sha1_c[16:16+4]...)
	aes_iv = append(aes_iv, sha1_d[0:8]...)

	return aes_key, aes_iv
}
