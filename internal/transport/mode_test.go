package transport_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xelaj/mtproto/internal/transport"
)

func TestMode(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	tr, err := transport.NewIntermediateMode(DummyConn(buf))
	require.NoError(t, err)

	require.NoError(t, tr.WriteMsg([]byte("test message")))
	require.Equal(t, buf.Bytes(), []byte{
		0xee, 0xee, 0xee, 0xee, 0x0c, 0x00, 0x00, 0x00,
		0x74, 0x65, 0x73, 0x74, 0x20, 0x6d, 0x65, 0x73,
		0x73, 0x61, 0x67, 0x65,
	})

	tr, err = transport.NewIntermediateMode(DummyConn(bytes.NewBuffer([]byte{
		0x0c, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74,
		0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	})))
	require.NoError(t, err)

	res, err := tr.ReadMsg()
	require.NoError(t, err)
	require.Equal(t, []byte("test message"), res)
}

func DummyConn(rw io.ReadWriter) transport.Conn { return nopCloser{ReadWriter: rw} }

type nopCloser struct {
	io.ReadWriter
}

func (nopCloser) Close() error { return nil }
