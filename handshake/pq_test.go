package handshake_test

import (
	"crypto/rand"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/xelaj/mtproto/v2/handshake"
)

func TestDecomposePQ(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		rand io.Reader
		pq   uint64
		p, q uint64
	}{
		{pq: 1724114033281923457, p: 1229739323, q: 1402015859},
		{pq: 378221, p: 613, q: 617},
		{pq: 15, p: 3, q: 5},

		// Testing vector taken from telegram docs.
		// * https://core.telegram.org/mtproto/samples-auth_key#4-encrypted-data-generation
		{pq: 0x17ED48941A08F981, p: 0x494C553B, q: 0x53911073},
		{pq: 2090522174869285481, p: 1112973847, q: 1878321023},
	} {
		tt := tt // for parallel tests
		if tt.rand == nil {
			tt.rand = rand.Reader
		}

		t.Run(tt.name, func(t *testing.T) {
			p, q := DecomposePQ(tt.pq, tt.rand)
			require.Equal(t, tt.p, p)
			require.Equal(t, tt.q, q)
		})
	}
}

func TestDecomposeRaw(t *testing.T) {
	for _, tt := range []struct {
		pq, p, q []byte
	}{
		{pq: Hexed("1ae945fd86042ea9"), p: Hexed("47625cd9"), q: Hexed("60827e51")},
	} {
		tt := tt // for parallel tests
		t.Run("", func(t *testing.T) {
			pq := big.NewInt(0).SetBytes(tt.pq).Uint64()

			p, q := DecomposePQ(pq, rand.Reader)

			require.Equal(t, tt.p, big.NewInt(0).SetUint64(p).Bytes())
			require.Equal(t, tt.q, big.NewInt(0).SetUint64(q).Bytes())

		})
	}
}
