package transport_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xelaj/mtproto/internal/transport"
)

func TestMode(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	transport, err := transport.NewIntermediateMode(buf)
	require.NoError(t, err)

	require.NoError(t, transport.WriteMsg([]byte("test message")))
	assert.Equal(t, buf.Bytes(), []byte{
		0xee, 0xee, 0xee, 0xee, 0x0c, 0x00, 0x00, 0x00,
		0x74, 0x65, 0x73, 0x74, 0x20, 0x6d, 0x65, 0x73,
		0x73, 0x61, 0x67, 0x65,
	})
}
