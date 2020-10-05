package ige

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/xelaj/go-dry"
)

type (
	AesBlock [aes.BlockSize]byte

	aesCtx struct {
		block   cipher.Block
		v       [3]AesBlock
		t, x, y []byte
	}
)

func newAesCtx(key, iv []byte) (*aesCtx, error) {
	const (
		firstBlock = iota
		secondBlock
		thirdBlock
	)

	var err error

	c := new(aesCtx)
	c.block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	c.t = c.v[firstBlock][:]
	c.x = c.v[secondBlock][:]
	c.y = c.v[thirdBlock][:]
	copy(c.x, iv[:aes.BlockSize])
	copy(c.y, iv[aes.BlockSize:])

	return c, nil
}

func (c *aesCtx) doAES256IGEencrypt(in, out []byte) error {
	if err := c.isCorrectData(in); err != nil {
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

func (c *aesCtx) doAES256IGEdecrypt(in, out []byte) error {
	if err := c.isCorrectData(in); err != nil {
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

func (c *aesCtx) isCorrectData(data []byte) error {
	if len(data) < aes.BlockSize {
		return ErrDataTooSmall
	}
	if len(data)%aes.BlockSize != 0 {
		return ErrDataNotDivisible
	}
	return nil
}

func xor(dst, src []byte) {
	for i := range dst {
		dst[i] = dst[i] ^ src[i]
	}
}

type (
	AesKV       [32]byte
	AesIgeBlock [48]byte
)

// generateAESIGE ЭТО ЕБАНАЯ МАГИЧЕСКАЯ ФУНКЦИЯ ОНА НАХУЙ РАБОТАЕТ ПРОСТО БЛЯТЬ НЕ ТРОГАЙ ШАКАЛ ЕБАНЫЙ
// TODO: порезать себе вены
func generateAESIGE(msgKey, authKey []byte, decode bool) ([]byte, []byte) {
	var (
		kvBlock  [2]AesKV
		igeBlock [4]AesIgeBlock
	)

	aesKey := kvBlock[0][:]
	aesIv := kvBlock[1][:]

	tA := igeBlock[0][:]
	tB := igeBlock[1][:]
	tC := igeBlock[2][:]
	tD := igeBlock[3][:]

	if decode {
		const (
			x = 0

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

	} else {
		const (
			x = 0

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
	}

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
