package handshake

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
)

func MakeGAB(rand io.Reader, g, gA, dh *big.Int) (gB, b *big.Int, err error) {
	if err := CheckDH(g, dh); err != nil {
		return nil, nil, fmt.Errorf("check dh: %w", err)
	}

	// 6. Random number b is computed:
	b, err = crand.Int(rand, big.NewInt(0).SetBit(big.NewInt(0), 2048, 1))
	if err != nil {
		panic(fmt.Errorf("number b generation: %w", err))
	}
	// g_b = g^b mod dh_prime
	gB = big.NewInt(0).Exp(g, b, dh)

	// Checking key exchange parameters.
	if err := CheckDHParams(dh, g, gA, gB); err != nil {
		return nil, nil, fmt.Errorf("key exchange failed: invalid params: %w", err)
	}

	return gB, b, nil
}

// CheckDH performs DH parameters check described in Telegram docs.
//
//	Client is expected to check whether p is a safe 2048-bit prime (meaning that both p and (p-1)/2 are prime,
//	and that 2^2047 < p < 2^2048), and that g generates a cyclic subgroup of prime order (p-1)/2, i.e.
//	is a quadratic residue mod p. Since g is always equal to 2, 3, 4, 5, 6 or 7, this is easily done using quadratic
//	reciprocity law, yielding a simple condition on p mod 4g — namely, p mod 8 = 7 for g = 2; p mod 3 = 2 for g = 3;
//	no extra condition for g = 4; p mod 5 = 1 or 4 for g = 5; p mod 24 = 19 or 23 for g = 6; and p mod 7 = 3,
//	5 or 6 for g = 7.
//
// See https://core.telegram.org/mtproto/auth_key#presenting-proof-of-work-server-authentication.
//
// See https://core.telegram.org/api/srp#checking-the-password-with-srp.
//
// See https://core.telegram.org/api/end-to-end#sending-a-request.
func CheckDH(g, p *big.Int) error {
	// The client is expected to check whether p is a safe 2048-bit prime
	// (meaning that both p and (p-1)/2 are prime, and that 2^2047 < p < 2^2048).
	// FIXME(tdakkota): we check that 2^2047 <= p < 2^2048
	// 	but docs says to check 2^2047 < p < 2^2048.
	//
	// TDLib check 2^2047 <= too:
	// https://github.com/tdlib/td/blob/d161323858a782bc500d188b9ae916982526c262/td/mtproto/DhHandshake.cpp#L23
	if p.BitLen() != 2048 {
		return errors.New("p should be 2^2047 < p < 2^2048")
	}

	if err := CheckGP(g, p); err != nil {
		return err
	}

	return checkPrime(p)
}

func checkPrime(p *big.Int) error {
	if !Prime(p) {
		return errors.New("p is not prime number")
	}

	sub := big.NewInt(0).Sub(p, big.NewInt(1))
	pr := sub.Quo(sub, big.NewInt(2))
	if !Prime(pr) {
		return errors.New("(p-1)/2 is not prime number")
	}

	return nil
}

// Prime checks that given number is prime.
func Prime(p *big.Int) bool {
	// TODO(tdakkota): maybe it should be smaller?
	// 1 - 1/4^64 is equal to 0.9999999999999999999999999999999999999970612641229442812300781587
	//
	// TDLib uses nchecks = 64
	// See https://github.com/tdlib/td/blob/d161323858a782bc500d188b9ae916982526c262/tdutils/td/utils/BigNum.cpp#L155.
	const probabilityN = 64

	// ProbablyPrime is mutating, so we need a copy
	return p.ProbablyPrime(probabilityN)
}

// CheckGP checks whether g generates a cyclic subgroup of prime order (p-1)/2, i.e. is a quadratic residue mod p.
// Also check that g is 2, 3, 4, 5, 6 or 7.
//
// This function is needed by some Telegram algorithms(Key generation, SRP 2FA).
//
// See https://core.telegram.org/mtproto/auth_key.
//
// See https://core.telegram.org/api/srp.
func CheckGP(g, p *big.Int) error {
	// Since g is always equal to 2, 3, 4, 5, 6 or 7,
	// this is easily done using quadratic reciprocity law, yielding a simple condition on p mod 4g -- namely,
	var result bool
	switch g.Uint64() {
	case 2:
		// p mod 8 = 7 for g = 2;
		result = checkSubgroup(p, 8, 7)
	case 3:
		// p mod 3 = 2 for g = 3;
		result = checkSubgroup(p, 3, 2)
	case 4:
		// no extra condition for g = 4
		result = true
	case 5:
		// p mod 5 = 1 or 4 for g = 5;
		result = checkSubgroup(p, 5, 1, 4)
	case 6:
		// p mod 24 = 19 or 23 for g = 6;
		result = checkSubgroup(p, 24, 19, 23)
	case 7:
		// and p mod 7 = 3, 5 or 6 for g = 7.
		result = checkSubgroup(p, 7, 3, 5, 6)
	default:
		return fmt.Errorf("unexpected g = %d: g should be equal to 2, 3, 4, 5, 6 or 7", g)
	}

	if !result {
		return errors.New("g should be a quadratic residue mod p")
	}

	return nil
}

func checkSubgroup(p *big.Int, divider int64, expected ...int64) bool {
	rem := new(big.Int).Rem(p, big.NewInt(divider)).Int64()

	for _, e := range expected {
		if rem == e {
			return true
		}
	}

	return false
}

// CheckDHParams checks that g_a, g_b and g params meet key exchange conditions.
//
// https://core.telegram.org/mtproto/auth_key#dh-key-exchange-complete
func CheckDHParams(dhPrime, g, gA, gB *big.Int) error {
	one := big.NewInt(1)
	dhPrimeMinusOne := big.NewInt(0).Sub(dhPrime, one)
	if !InRange(g, one, dhPrimeMinusOne) {
		return errors.New("kex: bad g, g must be 1 < g < dh_prime - 1")
	}
	if !InRange(gA, one, dhPrimeMinusOne) {
		return errors.New("kex: bad g_a, g_a must be 1 < g_a < dh_prime - 1")
	}
	if !InRange(gB, one, dhPrimeMinusOne) {
		return errors.New("kex: bad g_b, g_b must be 1 < g_b < dh_prime - 1")
	}

	// IMPORTANT: Apart from the conditions on the Diffie-Hellman prime
	// dh_prime and generator g, both sides are to check that g, g_a and
	// g_b are greater than 1 and less than dh_prime - 1. We recommend
	// checking that g_a and g_b are between 2^{2048-64} and
	// dh_prime - 2^{2048-64} as well.

	// 2^{2048-64}
	safetyRangeMin := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(2048-64), nil)
	safetyRangeMax := big.NewInt(0).Sub(dhPrime, safetyRangeMin)
	if !InRange(gA, safetyRangeMin, safetyRangeMax) {
		return errors.New("kex: bad g_a, g_a must be 2^{2048-64} < g_a < dh_prime - 2^{2048-64}")
	}
	if !InRange(gB, safetyRangeMin, safetyRangeMax) {
		return errors.New("kex: bad g_b, g_b must be 2^{2048-64} < g_b < dh_prime - 2^{2048-64}")
	}

	return nil
}

// InRange checks whether x is in (min, max) range, i.e. min < x < max.
func InRange(x, min, max *big.Int) bool {
	return x.Cmp(min) > 0 && x.Cmp(max) < 0
}
