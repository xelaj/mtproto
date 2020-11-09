package tl_test

import (
	"testing"

	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/telegram"
)

func BenchmarkEncoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tl.Encode(&telegram.AccountInstallThemeParams{
			Dark:   true,
			Format: "abc",
			Theme: &telegram.InputThemeObj{
				Id:         123,
				AccessHash: 321,
			},
		})
	}
}
