package srp

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"
	"golang.org/x/crypto/pbkdf2"
	"math/big"
)

// GetInputCheckPassword считает нужные для 2FA хеши, работает аналогично функции из tdlib:
// https://github.com/tdlib/td/blob/f9009cbc01e9c4c77d31120a61feb9c639c6aeda/td/telegram/PasswordManager.cpp#L72
//
// В random256Bytes нужно передать рандомные 256 байт, можно использовать Random256Bytes().
func GetInputCheckPassword(passwordStr string, accountPassword *telegram.AccountPassword, random256Bytes []byte) (telegram.InputCheckPasswordSRP, error) {
	// У CurrentAlgo должен быть этот самый тип, с длинным названием алгоритма
	// https://github.com/tdlib/td/blob/f9009cbc01e9c4c77d31120a61feb9c639c6aeda/td/telegram/AuthManager.cpp#L537
	current, ok := accountPassword.CurrentAlgo.(*telegram.PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow)
	if !ok {
		return nil, errors.New("invalid CurrentAlgo type")
	}

	password := []byte(passwordStr)
	clientSalt := current.Salt1
	serverSalt := current.Salt2
	g := current.G
	p := current.P
	B := accountPassword.SrpB
	id := accountPassword.SrpId

	if len(password) == 0 {
		return &telegram.InputCheckPasswordEmpty{}, nil
	}

	if dhHandshakeCheckConfigIsError(g, p) {
		return &telegram.InputCheckPasswordEmpty{}, errors.New("receive invalid config g")
	}

	p_bn := new(big.Int).SetBytes(p)
	B_bn := new(big.Int).SetBytes(B)
	zero := big.NewInt(0)

	if zero.Cmp(B_bn) != -1 || B_bn.Cmp(p_bn) != -1 || len(B) < 248 || len(B) > 256 {
		return &telegram.InputCheckPasswordEmpty{}, errors.New("receive invalid value of B")
	}

	g_bn := big.NewInt(int64(g))
	g_padded := make([]byte, 256)
	binary.BigEndian.PutUint32(g_padded[256-4:], uint32(g))

	x := passwordHash2(password, clientSalt, serverSalt)
	x_bn := new(big.Int).SetBytes(x)

	a_bn := new(big.Int).SetBytes(random256Bytes)

	A_bn := new(big.Int).Exp(g_bn, a_bn, p_bn)
	A := pad256(A_bn.Bytes())

	B_pad := make([]byte, 256-len(B))

	u256 := sha256.New()
	u256.Write(A)
	u256.Write(B_pad)
	u256.Write(B)
	u := u256.Sum(nil)
	u_bn := new(big.Int).SetBytes(u)

	k256 := sha256.New()
	k256.Write(p)
	k256.Write(g_padded)
	k := k256.Sum(nil)
	k_bn := new(big.Int).SetBytes(k)

	v_bn := new(big.Int).Exp(g_bn, x_bn, p_bn)

	kv_bn := new(big.Int).Mul(k_bn, v_bn)
	kv_bn.Mod(kv_bn, p_bn)

	t_bn := new(big.Int).Sub(B_bn, kv_bn)
	if t_bn.Cmp(zero) == -1 {
		t_bn.Add(t_bn, p_bn)
	}

	exp_bn := new(big.Int).Mul(u_bn, x_bn)
	exp_bn.Add(exp_bn, a_bn)

	S_bn := new(big.Int).Exp(t_bn, exp_bn, p_bn)
	S := pad256(S_bn.Bytes())
	K := sha256.Sum256(S)

	h1 := sha256.Sum256(p)
	h2 := sha256.Sum256(g_padded)
	for i := range h1 {
		h1[i] ^= h2[i]
	}

	clientSalt256 := sha256.Sum256(clientSalt)
	serverSalt256 := sha256.Sum256(serverSalt)

	M256 := sha256.New()
	M256.Write(h1[:])
	M256.Write(clientSalt256[:])
	M256.Write(serverSalt256[:])
	M256.Write(A)
	M256.Write(pad256(B))
	M256.Write(K[:])

	M := M256.Sum(nil)

	return &telegram.InputCheckPasswordSRPObj{
		SrpId: id,
		A:     A,
		M1:    M[:],
	}, nil
}

func Random256Bytes() []byte {
	return dry.RandomBytes(256)
}

func saltingHashing(data, salt []byte) []byte {
	h := sha256.New()
	h.Write(salt)
	h.Write(data)
	h.Write(salt)

	return h.Sum(nil)
}

func passwordHash1(password, salt1, salt2 []byte) []byte {
	return saltingHashing(saltingHashing(password, salt1), salt2)
}

func passwordHash2(password, salt1, salt2 []byte) []byte {
	return saltingHashing(pbkdf2sha512(passwordHash1(password, salt1, salt2), salt1, 100000), salt2)
}

func pbkdf2sha512(hash1 []byte, salt1 []byte, i int) []byte {
	return pbkdf2.Key(hash1, salt1, i, 64, sha512.New)
}

func pad256(b []byte) []byte {
	if len(b) >= 256 {
		return b[len(b)-256:]
	}

	tmp := make([]byte, 256)
	copy(tmp[256-len(b):], b)

	return tmp
}

func dhHandshakeCheckConfigIsError(gInt int32, primeStr []byte) bool {
	prime := new(big.Int).SetBytes(primeStr)
	_ = prime

	// Функция описана здесь, и что-то проверяет.
	// Реализовывать пока что лень.
	// TODO: запилить
	// https://github.com/tdlib/td/blob/f9009cbc01e9c4c77d31120a61feb9c639c6aeda/td/mtproto/DhHandshake.cpp#L18

	return false
}
