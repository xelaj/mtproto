// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package srp

//! WARNING: if you want to understand this algorithm, go to https://core.telegram.org/api/srp, and try to
//  open this code on right side, and algorith description on left side. Then, try to search via Cmd+F func
//  descriptions in algo descriptions. Then bless your god and drink few whiskey. As far as this way i can help
//  you to understand this secure-like shit created by telegram developers.

import (
	"crypto/sha256"
	"crypto/sha512"
	"math/big"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
	"golang.org/x/crypto/pbkdf2"
)

const (
	randombyteLen = 256 // 2048 bit
)

// GetInputCheckPassword считает нужные для 2FA хеши, описан в доке телеграма:
// https://core.telegram.org/api/srp#checking-the-password-with-srp
func GetInputCheckPassword(password string, srpB []byte, mp *ModPow) (*SrpAnswer, error) {
	return getInputCheckPassword(password, srpB, mp, dry.RandomBytes(randombyteLen))
}

func getInputCheckPassword(
	password string,
	srpB []byte,
	mp *ModPow,
	random []byte,
) (
	*SrpAnswer, error,
) {
	if password == "" {
		return nil, nil
	}

	err := validateCurrentAlgo(srpB, mp)
	if err != nil {
		return nil, errors.Wrap(err, "validating CurrentAlgo")
	}

	p := bytesToBig(mp.P)
	g := big.NewInt(int64(mp.G))
	gBytes := pad256(g.Bytes())

	// random 2048-bit number a
	a := bytesToBig(random)

	// g_a = pow(g, a) mod p
	ga := pad256(bigExp(g, a, p).Bytes())

	// g_b = srp_B
	gb := pad256(srpB)

	// u = H(g_a | g_b)
	u := bytesToBig(calcSHA256(ga, gb))

	// x = PH2(password, salt1, salt2)
	x := bytesToBig(passwordHash2([]byte(password), mp.Salt1, mp.Salt2))

	// v = pow(g, x) mod p
	v := bigExp(g, x, p)

	// k = (k * v) mod p
	k := bytesToBig(calcSHA256(mp.P, gBytes))

	// k_v = (k * v) % p
	kv := k.Mul(k, v).Mod(k, p)

	// t = (g_b - k_v) % p
	t := bytesToBig(srpB)
	if t.Sub(t, kv).Cmp(big.NewInt(0)) == -1 {
		t.Add(t, p)
	}

	// s_a = pow(t, a + u * x) mod p
	sa := pad256(bigExp(t, u.Mul(u, x).Add(u, a), p).Bytes())

	// k_a = H(s_a)
	ka := calcSHA256(sa)

	// M1 := H(H(p) xor H(g) | H2(salt1) | H2(salt2) | g_a | g_b | k_a)
	M1 := calcSHA256(
		dry.BytesXor(calcSHA256(mp.P), calcSHA256(gBytes)),
		calcSHA256(mp.Salt1),
		calcSHA256(mp.Salt2),
		ga,
		gb,
		ka,
	)

	return &SrpAnswer{
		GA: ga,
		M1: M1,
	}, nil
}

// this is simpler struct, copied from PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow
type ModPow struct {
	Salt1 []byte
	Salt2 []byte
	G     int32
	P     []byte
}

// copy of InputCheckPasswordSRPObj
type SrpAnswer struct {
	GA []byte
	M1 []byte
}

// Validating mod pow from server side. just works, don't touch.
func validateCurrentAlgo(srpB []byte, mp *ModPow) error {
	if dhHandshakeCheckConfigIsError(mp.G, mp.P) {
		return errors.New("receive invalid config g")
	}

	p := bytesToBig(mp.P)
	gb := bytesToBig(srpB)

	//?                        awwww so cute ref (^_^), try to guess ↓↓↓
	if big.NewInt(0).Cmp(gb) != -1 || gb.Cmp(p) != -1 || len(srpB) < 248 || len(srpB) > 256 {
		return errors.New("receive invalid value of B")
	}

	return nil
}

// SH(data, salt) := H(salt | data | salt)
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

// joining arrays into single one and calculating hash
// H(a | b | c)
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
	//prime := new(big.Int).SetBytes(primeStr)
	//_ = prime

	// Функция описана здесь, и что-то проверяет.
	// Реализовывать пока что лень.
	// TODO: запилить
	//       или болт положить, ¯\_(ツ)_/¯
	// https://github.com/tdlib/td/blob/f9009cbc01e9c4c77d31120a61feb9c639c6aeda/td/mtproto/DhHandshake.cpp

	return false
}
