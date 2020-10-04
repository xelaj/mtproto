package ige

import (
	"bytes"
	"crypto/aes"
	"errors"
	"math/big"

	"github.com/xelaj/go-dry"
)

func MessageKey(msg []byte) []byte {
	return dry.Sha1(string(msg))[4:20]
}

func Encrypt(msg, key []byte) ([]byte, error) {
	msgKey := MessageKey(msg)
	aesKey, aesIV := generateAESIGE(msgKey, key, false)

	y := make([]byte, len(msg)+((16-(len(msg)%16))&15)) // СУДЯ ПО ВСЕМУ вообще не уверен, но это видимо паддинг для добива блока, чтоб он делился на 256 бит
	copy(y, msg)
	return doAES256IGEencrypt(y, aesKey, aesIV)

}

// checkData это msgkey в понятиях мтпрото, нужно что бы проверить, успешно ли прошла расшифровка
func Decrypt(msg, key, checkData []byte) ([]byte, error) {
	aesKey, aesIV := generateAESIGE(checkData, key, true)
	result, err := doAES256IGEdecrypt(msg, aesKey, aesIV)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// generateAESIGE ЭТО ЕБАНАЯ МАГИЧЕСКАЯ ФУНКЦИЯ ОНА НАХУЙ РАБОТАЕТ ПРОСТО БЛЯТЬ НЕ ТРОГАЙ ШАКАЛ ЕБАНЫЙ
// TODO: порезать себе вены
func generateAESIGE(msg_key, auth_key []byte, decode bool) ([]byte, []byte) {
	var x int
	if decode {
		x = 8
	} else {
		x = 0
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

func doAES256IGEencrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("AES256IGE: data too small to encrypt")
	}
	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("AES256IGE: data not divisible by block size")
	}

	t := make([]byte, aes.BlockSize)
	x := make([]byte, aes.BlockSize)
	y := make([]byte, aes.BlockSize)
	copy(x, iv[:aes.BlockSize])
	copy(y, iv[aes.BlockSize:])
	encrypted := make([]byte, len(data))

	i := 0
	for i < len(data) {
		xor(x, data[i:i+aes.BlockSize])
		block.Encrypt(t, x)
		xor(t, y)
		x, y = t, data[i:i+aes.BlockSize]
		copy(encrypted[i:], t)
		i += aes.BlockSize
	}

	return encrypted, nil
}

func doAES256IGEdecrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("AES256IGE: data too small to decrypt")
	}
	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("AES256IGE: data not divisible by block size")
	}

	t := make([]byte, aes.BlockSize)
	x := make([]byte, aes.BlockSize)
	y := make([]byte, aes.BlockSize)
	copy(x, iv[:aes.BlockSize])
	copy(y, iv[aes.BlockSize:])
	decrypted := make([]byte, len(data))

	i := 0
	for i < len(data) {
		xor(y, data[i:i+aes.BlockSize])
		block.Decrypt(t, y)
		xor(t, x)
		y, x = t, data[i:i+aes.BlockSize]
		copy(decrypted[i:], t)
		i += aes.BlockSize
	}

	return decrypted, nil

}

func xor(dst, src []byte) {
	for i := range dst {
		dst[i] = dst[i] ^ src[i]
	}
}

// DecryptMessageWithTempKeys дешифрует сообщение паролем, которые получены в процессе обмена ключами диффи хеллмана
func DecryptMessageWithTempKeys(msg []byte, nonceSecond, nonceServer *big.Int) []byte {
	key, iv := generateTempKeys(nonceSecond, nonceServer)
	decodedWithHash, err := doAES256IGEdecrypt(msg, key, iv)
	dry.PanicIfErr(err)

	// decodedWithHash := SHA1(answer) + answer + (0-15 рандомных байт); длина должна делиться на 16;
	decodedHash := decodedWithHash[:20]
	decodedMessage := decodedWithHash[20:]

	// режем последние 0-15 байт ориентируюясь по хешу
	for i := len(decodedMessage) - 1; i > len(decodedMessage)-16; i-- {
		if bytes.Equal(decodedHash, dry.Sha1Byte(decodedMessage[:i])) {
			return decodedMessage[:i]
		}
	}

	panic("couldn't trim message: hashes incompatible on more than 16 tries")
}

// EncryptMessageWithTempKeys шифрует сообщение паролем, которые получены в процессе обмена ключами диффи хеллмана
func EncryptMessageWithTempKeys(msg []byte, nonceSecond, nonceServer *big.Int) []byte {
	hash := dry.Sha1Byte(msg)

	// добавляем остаток рандомных байт в сообщение, что бы суммарно оно делилось на 16
	totalLen := len(hash) + len(msg)
	overflowedLen := totalLen % 16
	needToAdd := 16 - overflowedLen
	dry.RandomBytes(needToAdd)

	msg = bytes.Join([][]byte{hash, msg, dry.RandomBytes(needToAdd)}, []byte{})

	key, iv := generateTempKeys(nonceSecond, nonceServer)
	encodedWithHash, err := doAES256IGEencrypt(msg, key, iv)
	dry.PanicIfErr(err)

	return encodedWithHash
}

// https://tlgrm.ru/docs/mtproto/auth_key#server-otvecaet-dvuma-sposobami
// generateTempKeys генерирует временные ключи для шифрования в процессе обемна ключами.
func generateTempKeys(nonceSecond, nonceServer *big.Int) (key, iv []byte) {
	t1 := make([]byte, 48) // nonceSecond + nonceServer
	copy(t1[0:], nonceSecond.Bytes())
	copy(t1[32:], nonceServer.Bytes())
	hash1 := dry.Sha1Byte(t1) // SHA1(nonceSecond + nonceServer)

	t2 := make([]byte, 48) // nonceServer + nonceSecond
	copy(t2[0:], nonceServer.Bytes())
	copy(t2[16:], nonceSecond.Bytes())
	hash2 := dry.Sha1Byte(t2) // SHA1(nonceServer + nonceSecond)

	tmpAESKey := make([]byte, 32)     // SHA1(nonceSecond + nonceServer) + substr (SHA1(nonceServer + nonceSecond), 0, 12);
	copy(tmpAESKey[0:], hash1)        // SHA1(nonceSecond + nonceServer)
	copy(tmpAESKey[20:], hash2[0:12]) // substr (SHA1(nonceServer + nonceSecond), 0, 12)

	t3 := make([]byte, 64) // nonceSecond + nonceSecond
	copy(t3[0:], nonceSecond.Bytes())
	copy(t3[32:], nonceSecond.Bytes())
	hash3 := dry.Sha1Byte(t3) // SHA1(nonceSecond + nonceSecond)

	tmpAESIV := make([]byte, 32)                  // substr (SHA1(server_nonce + new_nonce), 12, 8) + SHA1(new_nonce + new_nonce) + substr (new_nonce, 0, 4);
	copy(tmpAESIV[0:], hash2[12:12+8])            // substr (SHA1(nonceServer + nonceSecond), 12, 8)
	copy(tmpAESIV[8:], hash3)                     // SHA1(nonceSecond + nonceSecond)
	copy(tmpAESIV[28:], nonceSecond.Bytes()[0:4]) // substr (nonceSecond, 0, 4)

	return tmpAESKey, tmpAESIV
}
