package handshake

import (
	crand "crypto/rand"
	"io"
	"math/big"
)

// DecomposePQ decomposes pq into prime factors such that p < q.
func DecomposePQ(pq uint64, rand io.Reader) (p, q uint64) { // nolint:gocognit
	var (
		value0  = big.NewInt(0)
		value1  = big.NewInt(1)
		value15 = big.NewInt(15)
		value17 = big.NewInt(17)
		rndMax  = big.NewInt(0).SetBit(big.NewInt(0), 64, 1)

		y        = big.NewInt(0)
		whatNext = big.NewInt(0)

		a = big.NewInt(0)
		b = big.NewInt(0)
		c = big.NewInt(0)

		b2 = big.NewInt(0)

		z = big.NewInt(0)
	)

	what := big.NewInt(0).SetUint64(pq)
	g := big.NewInt(0)
	i := 0
	for !(g.Cmp(value1) == 1 && g.Cmp(what) == -1) {
		v, err := crand.Int(rand, rndMax)
		if err != nil {
			panic(err)
		}
		v = v.And(v, value15)
		v = v.Add(v, value17)
		v = v.Mod(v, what)

		x, err := crand.Int(rand, rndMax)
		if err != nil {
			panic(err)
		}
		whatNext.Sub(what, value1)
		x = x.Mod(x, whatNext)
		x = x.Add(x, value1)

		y.Set(x)
		lim := 1 << (uint(i) + 18)
		j := 1
		flag := true

		for j < lim && flag {
			a.Set(x)
			b.Set(x)
			c.Set(v)

			for b.Cmp(value0) == 1 {
				b2.SetInt64(0)
				if b2.And(b, value1).Cmp(value0) == 1 {
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

			z.SetInt64(0)
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

			if g.Cmp(value1) != 0 {
				flag = false
			}
		}
		i++
	}

	p = g.Uint64()
	q = big.NewInt(0).Div(what, g).Uint64()

	if p > q {
		p, q = q, p
	}

	return p, q
}
