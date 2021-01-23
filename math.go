// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

//nolint:gochecknoglobals using it just for simplification and more readable
package mtproto

import (
	"crypto/rsa"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/xelaj/go-dry"
)

var (
	// числа используются только для алгоритмов
	big0  = big.NewInt(0)
	big1  = big.NewInt(1)
	big15 = big.NewInt(15)
	big17 = big.NewInt(17)
)

// doRSAencrypt шифрует ровно 1 блок сообщения длиной 255 байт публичным ключом.
// специфический алгоритм для мтпрото, т.к. документация не указывает, шифрование
// по OAEP или как-то еще
func doRSAencrypt(block []byte, key *rsa.PublicKey) []byte {
	dry.PanicIf(len(block) != math.MaxUint8, "block size isn't equal 255 bytes")
	z := big.NewInt(0).SetBytes(block)
	exponent := big.NewInt(int64(key.E))

	c := big.NewInt(0).Exp(z, exponent, key.N)

	res := make([]byte, 256)
	copy(res, c.Bytes())

	return res
}

// splitPQ раскладывает число на два простых, при том таким образом, что p1 < p2
// Часть алгоритма диффи хеллмана, как работает — без понятия
func splitPQ(pq *big.Int) (p1, p2 *big.Int) {
	rndmax := big.NewInt(0).SetBit(big.NewInt(0), 64, 1)

	what := big.NewInt(0).Set(pq)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint: gosec смысла нет
	g := big.NewInt(0)
	i := 0
	for !(g.Cmp(big1) == 1 && g.Cmp(what) == -1) {
		q := big.NewInt(0).Rand(rnd, rndmax)
		q = q.And(q, big15)
		q = q.Add(q, big17)
		q = q.Mod(q, what)

		x := big.NewInt(0).Rand(rnd, rndmax)
		whatnext := big.NewInt(0).Sub(what, big1)
		x = x.Mod(x, whatnext)
		x = x.Add(x, big1)

		y := big.NewInt(0).Set(x)
		lim := 1 << (uint(i) + 18)
		j := 1
		flag := true

		for j < lim && flag {
			a := big.NewInt(0).Set(x)
			b := big.NewInt(0).Set(x)
			c := big.NewInt(0).Set(q)

			for b.Cmp(big0) == 1 {
				b2 := big.NewInt(0)
				if b2.And(b, big1).Cmp(big0) == 1 {
					c.Add(c, a)
					if c.Cmp(what) >= 0 {
						c.Sub(c, what)
					}
				}
				a.Add(a, a)
				if a.Cmp(what) >= 0 {
					a.Sub(a, what)
				}
				b.Rsh(b, 1)
			}
			x.Set(c)

			z := big.NewInt(0)
			if x.Cmp(y) == -1 {
				z.Add(what, x)
				z.Sub(z, y)
			} else {
				z.Sub(x, y)
			}
			g.GCD(nil, nil, z, what)

			if (j & (j - 1)) == 0 {
				y.Set(x)
			}
			j++

			if g.Cmp(big1) != 0 {
				flag = false
			}
		}
		i++
	}

	p1 = big.NewInt(0).Set(g)
	p2 = big.NewInt(0).Div(what, g)

	if p1.Cmp(p2) == 1 {
		p1, p2 = p2, p1
	}

	return p1, p2
}

func makeGAB(g int32, g_a, dh_prime *big.Int) (b, g_b, g_ab *big.Int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint: gosec зачем
	rndmax := big.NewInt(0).SetBit(big.NewInt(0), 2048, 1)
	b = big.NewInt(0).Rand(rnd, rndmax)
	g_b = big.NewInt(0).Exp(big.NewInt(int64(g)), b, dh_prime)
	g_ab = big.NewInt(0).Exp(g_a, b, dh_prime)

	return
}

func xor(dst, src []byte) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}
