// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

// ! IMPORTANT! if someone can write normal tests and rewrite the mishmash from
// 2020 - please write tests on cyclic encryption (both ways with one key),
// because I failed.

package payload_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/xelaj/mtproto/v2/internal/payload"
)

func TestSerializeUnencryptedMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *Unencrypted
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "usualMessage",
			msg: &Unencrypted{
				Msg: []byte("hello mtproto messages!"),
				ID:  123,
			}, //        |   authKeyHash||    msgID     ||mlen ||rest of data                 >>
			want: Hexed("00000000000000007b000000000000001700000068656c6c6f206d7470726f746f206d65" +
				"73736167657321"),
		},
	}
	for _, tt := range tests {
		tt.wantErr = noErrAsDefault(tt.wantErr)

		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.Serialize()
			assert.Equal(t, tt.want, got)

			back, err := DeserializeUnencrypted(got)
			require.NoError(t, err)
			require.Equal(t, tt.msg, back)
		})
	}
}

func TestEncryptedMessage(t *testing.T) {
	for _, tt := range []struct {
		name         string
		skip         bool
		msg          Encrypted
		key          Key
		side         Side
		rand         io.Reader
		wantEnvelope []byte
		want         []byte
		wantErr      assert.ErrorAssertionFunc
	}{{
		skip: true,
		msg: Encrypted{
			Salt:      567,
			SessionID: 123,
			ID:        7323362318062205344,
			SeqNo:     4,
			Msg:       []byte("hello mtproto messages!"),
		},
		key:  val(new([256]byte)),
		side: SideClient,
		rand: nopReader{},
		want: Hexed("b1080bfd570d9b91bcaf69c99a389f94069be336c4d4f7228a85861a" +
			"dfcf48b0d2db155be20418a4cf949aaa4129a940e775e6bd270504606ef2ce49" +
			"8867c53e90dac6a016c81f853547c1255bbaa74e3f3b0b2a3dc0ce60"),
	}, {
		name: "real life",
		msg: Encrypted{
			Salt:      asUint(-1061050580851877637),
			SessionID: asUint(-956258382667033759),
			ID:        7324708596786199632,
			SeqNo:     7,
			Msg: []byte{
				0x4f, 0x24, 0x77, 0xa6, 0x0b, 0x2b, 0x34, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
				0x39, 0x30, 0x00, 0x00, 0x06, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x00, 0x78, 0x3d, 0x25, 0xad,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		rand: nopReader{},
		key: makeKey(`
			7e46d98059fa2b277725b8fbc47ad28368269f0ddf3c319d8132ece3c0854812
			b904aaf9bbe28b0ba77560767b6d97000c79a4a47f5f3ead3c0ca55328956f83
			4abc70a1d27ed8b89e78e017cc667e708d17d40b63a739404f03bfe78920745c
			e5143830867b93f5431ebf5abeb6b32343b9e09f44824fd617d9aa4ea712ae7b
			fc965729595164cbb996ac43061ed217e97073b012f255d2211668a8b7db104f
			eaedd78a0a7c99db0eaf64924cfe6c8ac5e30931b3107bb3b3f973c67e36887a
			ddba51c08f10ffbe5dd48fa97aea2e86e9befac361533c09f69315e3acc6c939
			dd3aacde82d4f3349a08e9cf6d9c8d85f9c8bf8f08b99a0d7ff092cb689564c3
		`),
		side: SideClient,
		want: []byte{
			0x1d, 0x51, 0x23, 0xfd, 0x04, 0xf8, 0xef, 0x8b, 0xb5, 0x21, 0x2d, 0xd3, 0x07, 0xb2, 0x2d, 0x01,
			0x87, 0xb8, 0x8e, 0x4a, 0xab, 0x39, 0x70, 0xc0, 0x76, 0x55, 0x15, 0xaa, 0xbb, 0xd5, 0xc8, 0xb4,
			0xb7, 0xdc, 0xba, 0x7b, 0xf3, 0xd7, 0x7f, 0x64, 0xb8, 0xda, 0xa1, 0x9b, 0xc6, 0xeb, 0x6e, 0x18,
			0x03, 0xdb, 0xd2, 0x4c, 0x4a, 0xd8, 0x3d, 0x4d, 0xe6, 0x84, 0xdd, 0x0e, 0xec, 0x8a, 0x3a, 0x44,
			0x53, 0xf5, 0x2a, 0x26, 0x41, 0x67, 0x23, 0x7f, 0x55, 0x06, 0x0b, 0x43, 0xe2, 0xbd, 0xc5, 0xed,
			0xa3, 0x1f, 0xa6, 0xf7, 0x41, 0x8f, 0xce, 0xe3, 0xab, 0x3c, 0x27, 0x87, 0x22, 0xd5, 0x91, 0xa5,
			0xc3, 0xaf, 0xc8, 0xcb, 0xd1, 0xdd, 0xb6, 0x24, 0x18, 0xbf, 0x12, 0x42, 0xe5, 0xbe, 0xc7, 0x38,
			0x6e, 0x2a, 0x22, 0x73, 0x85, 0x50, 0xa2, 0x8d,
		},
	}} {
		tt.wantErr = noErrAsDefault(tt.wantErr)

		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}

			cipher := NewClientCipher(tt.rand, tt.key)

			got, err := cipher.Encrypt([8]byte{}, BuildEnvelope(tt.msg.Salt, tt.msg.SessionID, tt.msg.SeqNo, tt.msg.ID, tt.msg.Msg, tt.rand))
			require.NoError(t, err)
			require.Equal(t, tt.want, got, "want")

			decrypted, err := cipher.Decrypt(got)
			require.NoError(t, err)

			e, err := DeserializeEnvelope(decrypted)
			require.NoError(t, err)
			require.Equal(t, tt.msg.Msg, e.Msg)
		})
	}
}

func val[T any](value *T) T { return *value }

func dummyKey() (res [256]byte) {
	for i := range res {
		res[i] = byte(i)
	}

	return
}

func dummyMsgKey() (res [16]byte) {
	for i := range res {
		res[i] = byte(i)
	}

	return
}

type nopReader struct{}

func (nopReader) Read(p []byte) (n int, err error) { clear(p); return len(p), nil }

type iterReader struct{ i byte }

func (r *iterReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = r.i + byte(i)
	}
	r.i += byte(len(p))
	return len(p), nil
}

func asUint(i int) uint64 { return uint64(i) }
