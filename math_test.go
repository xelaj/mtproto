// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"math/big"
	"testing"
)

func TestSplitPQ(t *testing.T) {
	cases := []struct {
		pq, p, q *big.Int
	}{
		{big.NewInt(1724114033281923457), big.NewInt(1229739323), big.NewInt(1402015859)},
		{big.NewInt(378221), big.NewInt(613), big.NewInt(617)},
		{big.NewInt(15), big.NewInt(3), big.NewInt(5)},
	}

	for _, c := range cases {
		r1, r2 := splitPQ(c.pq)
		if c.p.Cmp(r1) != 0 || c.q.Cmp(r2) != 0 {
			t.Errorf("PQ mismatch: %v %v, want %v %v", r1, r2, c.p, c.q)
		}
	}
}
