package session_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/xelaj/mtproto/internal/session"
)

func TestDecodeBool(t *testing.T) {
	for _, tt := range []struct {
		name string
		key  [256]byte
		want uint64
	}{{
		key:  dummyKey(),
		want: binary.LittleEndian.Uint64([]byte{0x32, 0xd1, 0x58, 0x6e, 0xa4, 0x57, 0xdf, 0xc8}),
	}, {
		key:  makeKey(`0396efa1b10e7b7e99a9acbd476886c07015a64b3a1b9f16f2409af008f515bf0854942e6e4a596079115fc35980af13495621dc2614530cd817a9e16abef7b8efb3f6dbceb0b4600e96abb3fe88dd390d3b92dd31fa4150f38de82b8c5380664df2c03d94aa380f441e89158c64247085b5167df1d7fbb137a52c8d8a2e07a59424332cf52df60a9e67c90a00d53a31434dd9e80303a4d915f45afae95a68ab6a3804dab6a9b6d4857036bc22cd8f79ac6dfa932159b78e5bf880f4e9c3cce20dd1104b4decee6bbf14eaf08335f85c31af7fb896f68f48335bdc7ed50ec97338d013180c980f4647baeb88c5447f4c7ba485aefd9719ee2bcdca61faaf3ee6`),
		want: binary.LittleEndian.Uint64(Hexed(`6e31bd4029daf387`)),
	}} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := Session{Key: tt.key}

			require.Equal(t, tt.want, s.KeyID())
		})
	}
}

func dummyKey() (res [256]byte) {
	for i := range res {
		res[i] = byte(i)
	}

	return
}

func makeKey(in string) (res [256]byte) { copy(res[:], Hexed(in)); return res }

func Hexed(in string) []byte {
	reader := bytes.NewReader([]byte(in))
	buf := []rune{}
	for {
		r, ok := readByte(reader)
		if !ok {
			break
		}
		if r != 0 {
			buf = append(buf, r)
		}
	}

	return checkFunc(hex.DecodeString(string(buf)))
}

func readByte(reader *bytes.Reader) (rune, bool) {
	r, ok := readAndCheck(reader)
	if !ok {
		return 0, false
	}
	switch r {
	case ' ', '\n', '\t':
		return 0, true

	case '/':
		if r, ok := readAndCheck(reader); !ok || r != '/' {
			panic("expected comment")
		}
		skipComment(reader)

		return 0, true

	default:
		return r, true
	}
}

func skipComment(reader *bytes.Reader) {
	for {
		r, ok := readAndCheck(reader)
		if !ok || r == '\n' {
			break
		}
	}
}

func readAndCheck(reader *bytes.Reader) (r rune, ok bool) {
	r, _, err := reader.ReadRune()
	if err == io.EOF {
		return 0, false
	}
	check(err)

	return r, true
}

func checkFunc[T any](res T, err error) T {
	check(err)

	return res
}
