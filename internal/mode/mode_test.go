package mode_test

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	mode "github.com/xelaj/mtproto/internal/mode"
)

func TestModeEncode(t *testing.T) {
	randomBigByteset := make([]byte, 0x5d0)
	rand.Read(randomBigByteset)

	for _, tt := range []struct {
		name   string
		in     []byte
		mode   mode.Variant
		expect []byte
	}{
		{
			name: "intermediate, main mode",
			in:   []byte("test message"),
			mode: mode.Intermediate,
			expect: []byte{
				0xee, 0xee, 0xee, 0xee, 0x0c, 0x00, 0x00, 0x00,
				0x74, 0x65, 0x73, 0x74, 0x20, 0x6d, 0x65, 0x73,
				0x73, 0x61, 0x67, 0x65,
			},
		},
		{
			name: "arbiged, most unstable",
			in:   []byte("test message"),
			mode: mode.Abridged,
			expect: []byte{
				0xef, 0x03, 0x74, 0x65, 0x73, 0x74, 0x20, 0x6d,
				0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
			},
		},
		{
			name: "arbiged, but huge message",
			in:   randomBigByteset,
			mode: mode.Abridged,
			expect: append([]byte{
				0xef, 0x7f, 0x74, 0x01, 0x00}, randomBigByteset...),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			m, err := tt.mode(NopCloser(buf), true)
			require.NoError(t, err)

			err = m.WriteMsg(tt.in)
			require.NoError(t, err)

			require.Equal(t, tt.expect, buf.Bytes())
		})
	}
}

func TestModeDecode(t *testing.T) {
	randomBigByteset := make([]byte, 0x5d0)
	rand.Read(randomBigByteset)

	for _, tt := range []struct {
		name   string
		in     []byte
		mode   mode.Variant
		expect []byte
	}{
		{
			name: "intermediate, main mode",
			in: []byte{
				0xee, 0xee, 0xee, 0xee, 0x0c, 0x00, 0x00, 0x00,
				0x74, 0x65, 0x73, 0x74, 0x20, 0x6d, 0x65, 0x73,
				0x73, 0x61, 0x67, 0x65,
			},
			mode:   mode.Intermediate,
			expect: []byte("test message"),
		},
		{
			name: "arbiged, most unstable",
			in: []byte{
				0xef, 0x03, 0x74, 0x65, 0x73, 0x74, 0x20, 0x6d,
				0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
			},
			mode:   mode.Abridged,
			expect: []byte("test message"),
		},
		{
			name: "arbiged, but huge message",
			in: append([]byte{
				0xef, 0x7f, 0x74, 0x01, 0x00}, randomBigByteset...),
			mode:   mode.Abridged,
			expect: randomBigByteset,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(tt.in)

			m, err := mode.Detect(NopCloser(buf))
			require.NoError(t, err)

			got, err := m.ReadMsg(context.Background())
			require.NoError(t, err)
			require.Equal(t, tt.expect, got)
		})
	}
}

func NopCloser(r io.ReadWriter) io.ReadWriteCloser {
	return nopCloser{r}
}

type nopCloser struct {
	io.ReadWriter
}

func (nopCloser) Close() error { return nil }
