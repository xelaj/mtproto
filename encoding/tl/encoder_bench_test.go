// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

import (
	"testing"

	"github.com/xelaj/mtproto/encoding/tl"
)

func BenchmarkEncoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tl.Marshal(&AccountInstallThemeParams{
			Dark:   true,
			Format: "abc",
			Theme: &InputThemeObj{
				ID:         123,
				AccessHash: 321,
			},
		})
	}
}
