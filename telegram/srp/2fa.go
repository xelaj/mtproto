package srp

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"
	"golang.org/x/crypto/pbkdf2"
	"math/big"
)

// GetInputCheckPassword считает нужные для 2FA хеши, описан в доке телеграма:
// https://core.telegram.org/api/srp#checking-the-password-with-srp
//
// В random256Bytes нужно передать рандомные 256 байт, можно использовать Random256Bytes().
func GetInputCheckPassword(password string, accountPassword *telegram.AccountPassword, random256Bytes []byte) (telegram.InputCheckPasswordSRP, error) {
	if password == "" {
		return &telegram.InputCheckPasswordEmpty{}, nil
	}

	// У CurrentAlgo должен быть этот самый тип, с длинным названием алгоритма
	// https://github.com/tdlib/td/blob/f9009cbc01e9c4c77d31120a61feb9c639c6aeda/td/telegram/AuthManager.cpp#L537
	current, ok := accountPassword.CurrentAlgo.(*telegram.PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow)
	if !ok {
		return nil, errors.New("invalid CurrentAlgo type")
	}

	if dhHandshakeCheckConfigIsError(current.G, current.P) {
		return &telegram.InputCheckPasswordEmpty{}, errors.New("receive invalid config g")
	}

	p := bytesToBig(current.P)
	gbBig := bytesToBig(accountPassword.SrpB)
	zero := big.NewInt(0)

	if zero.Cmp(gbBig) != -1 || gbBig.Cmp(p) != -1 || len(accountPassword.SrpB) < 248 || len(accountPassword.SrpB) > 256 {
		return &telegram.InputCheckPasswordEmpty{}, errors.New("receive invalid value of B")
	}

	g := big.NewInt(int64(current.G))
	gBytes := pad256(g.Bytes())

	// random 2048-bit number a
	a := bytesToBig(random256Bytes)

	// g_a = pow(g, a) mod p
	ga := pad256(bigExp(g, a, p).Bytes())

	// g_b = srp_B
	gb := pad256(accountPassword.SrpB)

	// u = H(g_a | g_b)
	u := bytesToBig(calcSHA256(ga, gb))

	// x = PH2(password, salt1, salt2)
	x := bytesToBig(passwordHash2([]byte(password), current.Salt1, current.Salt2))

	// v = pow(g, x) mod p
	v := bigExp(g, x, p)

	// k = (k * v) mod p
	k := bytesToBig(calcSHA256(current.P, gBytes))

	// k_v = (k * v) % p
	kv := k.Mul(k, v).Mod(k, p)

	// t = (g_b - k_v) % p
	t := new(big.Int).Sub(gbBig, kv)
	if t.Cmp(zero) == -1 {
		t.Add(t, p)
	}

	// s_a = pow(t, a + u * x) mod p
	sa := pad256(bigExp(t, u.Mul(u, x).Add(u, a), p).Bytes())

	// k_a = H(s_a)
	ka := calcSHA256(sa)

	// M1 := H(H(p) xor H(g) | H2(salt1) | H2(salt2) | g_a | g_b | k_a)
	M1 := calcSHA256(
		xorBytes(calcSHA256(current.P), calcSHA256(gBytes)),
		calcSHA256(current.Salt1),
		calcSHA256(current.Salt2),
		ga,
		gb,
		ka,
	)

	return &telegram.InputCheckPasswordSRPObj{
		SrpId: accountPassword.SrpId,
		A:     ga,
		M1:    M1,
	}, nil
}

func Random256Bytes() []byte {
	return dry.RandomBytes(256)
}

func saltingHashing(data, salt []byte) []byte {
	return calcSHA256(salt, data, salt)
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

func xorBytes(a, b []byte) []byte {
	c := make([]byte, len(a))
	for i := range c {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func calcSHA256(arrays ...[]byte) []byte {
	h := sha256.New()
	for _, arr := range arrays {
		h.Write(arr)
	}
	return h.Sum(nil)
}

func bytesToBig(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func bigExp(x, y, m *big.Int) *big.Int {
	return new(big.Int).Exp(x, y, m)
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
