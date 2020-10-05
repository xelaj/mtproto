package ige

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/xelaj/go-dry"
)

type aesCtx struct {
	block   cipher.Block
	v       [3 * aes.BlockSize]byte
	t, x, y []byte
}

func newAesCtx(key, iv []byte) (*aesCtx, error) {
	const (
		firstOffStart = iota * aes.BlockSize
		secondOffStart
		thirdOffStart
		fourthOffsetStart
	)

	var err error

	c := new(aesCtx)

	c.block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	c.t = c.v[firstOffStart:secondOffStart]
	c.x = c.v[secondOffStart:thirdOffStart]
	c.y = c.v[thirdOffStart:fourthOffsetStart]
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

// generateAESIGE ЭТО ЕБАНАЯ МАГИЧЕСКАЯ ФУНКЦИЯ ОНА НАХУЙ РАБОТАЕТ ПРОСТО БЛЯТЬ НЕ ТРОГАЙ ШАКАЛ ЕБАНЫЙ
// TODO: порезать себе вены
func generateAESIGE(msgKey, authKey []byte, decode bool) ([]byte, []byte) {
	var x int
	if decode {
		x = 8
	} else {
		x = 0
	}
	aesKey := make([]byte, 0, 32)
	aesIv := make([]byte, 0, 32)
	tA := make([]byte, 0, 48)
	tB := make([]byte, 0, 48)
	tC := make([]byte, 0, 48)
	tD := make([]byte, 0, 48)

	tA = append(tA, msgKey...)
	tA = append(tA, authKey[x:x+32]...)

	tB = append(tB, authKey[32+x:32+x+16]...)
	tB = append(tB, msgKey...)
	tB = append(tB, authKey[48+x:48+x+16]...)

	tC = append(tC, authKey[64+x:64+x+32]...)
	tC = append(tC, msgKey...)

	tD = append(tD, msgKey...)
	tD = append(tD, authKey[96+x:96+x+32]...)

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
