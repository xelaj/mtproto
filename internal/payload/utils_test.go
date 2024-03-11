// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultAllocRuns = 10

// MaxAlloc checks that f does not allocate more than n.
func MaxAlloc(t *testing.T, n int, f func()) {
	t.Helper()
	if true {
		t.Skip("Skipped (race detector conflicts with allocation tests)")
	}
	avg := testing.AllocsPerRun(defaultAllocRuns, f)
	if avg > float64(n) {
		t.Errorf("Allocated %f bytes per run, expected less than %d", avg, n)
	}
}

// ZeroAlloc checks that f does not allocate.
func ZeroAlloc(t *testing.T, f func()) {
	t.Helper()
	MaxAlloc(t, 0, f)
}

func randSeed(data []byte) int64 {
	if len(data) == 0 {
		return 0
	}

	seedBuf := make([]byte, 64/8)
	copy(seedBuf, data)

	return int64(binary.BigEndian.Uint64(seedBuf))
}

// Rand returns a new rand.Rand with source deterministically initialized
// from seed byte slice.
//
// Zero length seed (or nil) is valid input.
func predictedRand(seed []byte) *rand.Rand {
	return rand.New(rand.NewSource(randSeed(seed))) // #nosec
}

func noErrAsDefault(e assert.ErrorAssertionFunc) assert.ErrorAssertionFunc {
	if e == nil {
		return assert.NoError
	}

	return e
}
