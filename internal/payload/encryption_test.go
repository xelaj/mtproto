// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload_test

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/xelaj/mtproto/internal/payload"
)

func TestCalcKey(t *testing.T) {
	k := dummyKey()
	m := dummyMsgKey()

	for _, tt := range []struct {
		name    string
		authKey Key
		msgKey  Int128
		side    Side
		wantKey Int256
		wantIV  Int256
	}{{
		name:    "Client",
		authKey: k,
		msgKey:  m,
		side:    SideClient,
		wantKey: [32]byte{
			112, 78, 208, 156, 139, 65, 102, 138, 232, 249, 157, 36, 71, 56, 247, 29,
			189, 220, 68, 70, 155, 107, 189, 74, 168, 87, 61, 208, 66, 189, 5, 158,
		},
		wantIV: [32]byte{
			77, 38, 96, 0, 165, 80, 237, 171, 191, 76, 124, 228, 15, 208, 4, 60, 201, 34, 48,
			24, 76, 211, 23, 165, 204, 156, 36, 130, 253, 59, 147, 24,
		},
	}, {
		name:    "Server",
		authKey: k,
		msgKey:  m,
		side:    SideServer,
		wantKey: [32]byte{
			33, 119, 37, 121, 155, 36, 88, 6, 69, 129, 116, 161, 252, 251, 200, 131, 144, 104,
			7, 177, 80, 51, 253, 208, 234, 43, 77, 105, 207, 156, 54, 78,
		},
		wantIV: [32]byte{
			102, 154, 101, 56, 145, 122, 79, 165, 108, 163, 35, 96, 164, 49, 201, 22, 11, 228,
			173, 136, 113, 64, 152, 13, 171, 145, 206, 123, 220, 71, 255, 188,
		},
	}, {
		name: "real life",
		authKey: makeKey(`
			7e46d98059fa2b277725b8fbc47ad28368269f0ddf3c319d8132ece3c0854812
			b904aaf9bbe28b0ba77560767b6d97000c79a4a47f5f3ead3c0ca55328956f83
			4abc70a1d27ed8b89e78e017cc667e708d17d40b63a739404f03bfe78920745c
			e5143830867b93f5431ebf5abeb6b32343b9e09f44824fd617d9aa4ea712ae7b
			fc965729595164cbb996ac43061ed217e97073b012f255d2211668a8b7db104f
			eaedd78a0a7c99db0eaf64924cfe6c8ac5e30931b3107bb3b3f973c67e36887a
			ddba51c08f10ffbe5dd48fa97aea2e86e9befac361533c09f69315e3acc6c939
			dd3aacde82d4f3349a08e9cf6d9c8d85f9c8bf8f08b99a0d7ff092cb689564c3
		`),
		msgKey: makeInt128(`08a74d1af6d2bb7421dae64bfdc3b76e`),
		side:   SideServer,
		wantKey: makeInt256(`
			c62635c6967cfa6a1ef309a830651ae38a05de9e81b6c3e47ad9c69573e678be
		`),
		wantIV: makeInt256(`
			18b1a92f802b9f26298aa2ec1e3b7e82bd1ae7c377ccc150eff4211cdee5cdd6
		`),
	}} {
		tt := tt // for parallel tests

		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotIV := Keys(&tt.authKey, tt.msgKey, tt.side)
			require.Equal(t, tt.wantKey, gotKey)
			require.Equal(t, tt.wantIV, gotIV)
		})
	}
}

// checking that serializing and deserializing again got same result.
func TestEncryptAndBack(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		key  [256]byte
		rand io.Reader
		raw  []byte
		want []byte
	}{{
		key: makeKey(`
			4f322ecf1f8f93faca05939860d47ca41ab6950a4f5f49f7ccdd588b20aa510c
			1b1553fb88becb7625e2d61756dd18384ee2be5a7ffce6412e654190f5547c98
			f200bb981e9359c8f4529f15501d418bbbc9719cf6fba2c7f790720222923936
			2fc2640af48a8e2d4360269d256ba1aaa96916c57cdb913529dfbc14e824f59a
			eb2ff025cfe8eac1bc6813b1c56155f49074f75dd2322cea92513609fef8abce
			a16ab9c0f073c9bf39f38e250d3d5e6f5958027cb0ef763deadab7915d2ec458
			907c36f4366080be3f18f372afc51f84cb6c78e8d316731933b0af6e1c7f7ff2
			0df7963831c32123e99f025b7ab132735495fecab16f0f6f08e97360210aa8f9
		`),
		rand: nopReader{},
		raw: Hexed(`
			8fd5e78ad27041b0a7c6b5e51509ab86187570258d87a5650500000010000000
			48a5910d15c4b51c010000003fb1c1f77d71e37a1637b18c1a2c184a7467975c
			b96446283e1faa74bfdc3ae1741e65d8
		`),
		want: Hexed(`
			14b88ec458f153409aa13e82b05f42da
			60b56c2c7b5fffda55f53b9f672f0dd6edcb16336d704f78b6d9c6dbcfdf8071
			5bbfbde57e417a38d200d0c75a531cc10907270a611a4b62e77459f27a62b485
			9626594020f0a9fe5475098494948c5a08373d655cdc91ee1e55edde64928928
		`),
	}, {
		name: "real life",
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
		rand: nopReader{},
		raw: Hexed(`
			fb189a4f1e6446f16153a44311b0baf250fcb02aa997a6650700000024000000
			4f2477a60b2b34313233343536373839393000000641424344454600783d25ad
			00000000
		`),
		want: Hexed(`
			b5212dd307b22d0187b88e4aab3970c0
			765515aabbd5c8b4b7dcba7bf3d77f64b8daa19bc6eb6e1803dbd24c4ad83d4d
			e684dd0eec8a3a4453f52a264167237f55060b43e2bdc5eda31fa6f7418fcee3
			ab3c278722d591a5c3afc8cbd1ddb62418bf1242e5bec7386e2a22738550a28d
		`),
	}} {
		tt := tt // for parallel tests

		c := NewClientCipher(tt.rand, tt.key)

		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Encrypt([8]byte{}, tt.raw)
			require.NoError(t, err)

			require.Equal(t, tt.want, got)

			// decrypting back
			decrypted, err := c.Decrypt(got)
			require.NoError(t, err)

			// cutting to size of len, because encrypted message will be always
			// padded to aes block. Real length of message is encoded right
			// inside encrypted message, so padding won't go till tl parser, but
			// transport should handle this problem.
			require.Equal(t, tt.raw, decrypted[:len(tt.raw)])
		})
	}
}

func TestMessageKey2(t *testing.T) {
	for _, tt := range []struct {
		name string
		key  Key
		raw  []byte
		side Side
		want [16]byte
	}{{
		name: "real life",
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
		raw: Hexed(`
			c3348dce8560b514048e6fb2b295d2860190481f56e0a5650900000024000000
			016d5cf3d00fcf0056e0a56519ca4421900100000e4150495f49445f494e5641
			4c4944000709c48459dbda36f96a3da5
		`),
		side: SideServer,
		want: makeInt128(`e5e98757016becc21c66f0a0bdbdce53`),
	}, {
		name: "real life",
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
		raw: Hexed(`
			c3348dce8560b514048e6fb2b295d286d00fcf0056e0a5650700000024000000
			4f2477a60b2b34313233343536373839393000000641424344454600783d25ad
			0000000038c06a129b598b30654ed048586439169824a1ff70f7cd28a3145274
		`),
		side: SideClient,
		want: makeInt128(`056d7ea2e861b705a517a7508dc90d6e`),
	}} {
		tt := tt // for parallel tests

		t.Run(tt.name, func(t *testing.T) {
			got := getMsgKey(&tt.key, tt.side, tt.raw)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEnvelope(t *testing.T) {
	for _, tt := range []struct {
		name      string
		salt      uint64
		sessionID uint64
		msgID     MsgID
		seqno     uint32
		msg       []byte

		want []byte
	}{{
		salt:      0x7cf29f37e736cef7, // 9003433667518189303,
		sessionID: 0xd944c5317987dc6a, // 15655855020929703018,
		msgID:     0x65a5ec5e280b9640, // 7324520258130908736,
		seqno:     0x7,
		msg: Hexed(`
			4f2477a60b2b34313233343536373839393000000641424344454600783d25ad
			00000000
		`),
		want: Hexed(`
			f7ce36e7379ff27c6adc877931c544d940960b285eeca5650700000024000000
			4f2477a60b2b34313233343536373839393000000641424344454600783d25ad
			00000000000102030405060708090a0b
		`),
	}, {
		salt:      1034,
		sessionID: 2345512351,
		msgID:     3401235566,
		seqno:     1,
		msg:       []byte{1, 2, 3, 100, 112},
		want: Hexed(`
			0a040000000000009fadcd8b000000006ebcbaca000000000100000005000000
			0102036470000102030405060708090a
		`),
	}} {
		tt := tt // for parallel tests

		t.Run(tt.name, func(t *testing.T) {
			got := BuildEnvelope(tt.salt, tt.sessionID, tt.seqno, tt.msgID, tt.msg, &iterReader{})
			require.Equal(t, tt.want, got)

			deserialized, err := DeserializeEnvelope(got)
			require.NoError(t, err)
			require.Equal(t, &Envelope{
				Salt:      tt.salt,
				SessionID: tt.sessionID,
				MsgID:     tt.msgID,
				SeqNo:     tt.seqno,
				Msg:       tt.msg,
			}, deserialized)
		})
	}
}

func BenchmarkKeys(b *testing.B) {
	k, m := genMessageAndAuthKeys()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Keys(&k, m, SideClient)
	}
}

func TestKeys(t *testing.T) {
	k, m := genMessageAndAuthKeys()

	ZeroAlloc(t, func() {
		_, _ = Keys(&k, m, SideClient)
	})
}

func BenchmarkMessageKey(b *testing.B) {
	k, _ := genMessageAndAuthKeys()
	payload := make([]byte, 1024)
	if _, err := io.ReadFull(rand.Reader, payload); err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = MessageKey(&k, payload, SideClient)
	}
}

func TestMessageKey(t *testing.T) {
	k, _ := genMessageAndAuthKeys()
	payload := make([]byte, 1024)
	if _, err := io.ReadFull(rand.Reader, payload); err != nil {
		t.Error(err)
	}

	ZeroAlloc(t, func() {
		_ = MessageKey(&k, payload, SideClient)
	})
}

func noErrAsDefault(e assert.ErrorAssertionFunc) assert.ErrorAssertionFunc {
	if e == nil {
		return assert.NoError
	}

	return e
}

func makeKey(in string) (res [256]byte)   { copy(res[:], Hexed(in)); return res }
func makeInt128(in string) (res [16]byte) { copy(res[:], Hexed(in)); return res }
func makeInt256(in string) (res [32]byte) { copy(res[:], Hexed(in)); return res }

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

func readByte(reader io.RuneReader) (rune, bool) {
	r, ok := readAndCheck(reader)
	if !ok {
		return 0, false
	}
	switch r {
	case ' ', '\n', '\t':
		return 0, true

	case '/':
		if r2, ok := readAndCheck(reader); !ok || r2 != '/' {
			panic("expected comment")
		}
		skipComment(reader)

		return 0, true

	default:
		return r, true
	}
}

func skipComment(reader io.RuneReader) {
	for {
		r, ok := readAndCheck(reader)
		if !ok || r == '\n' {
			break
		}
	}
}

func readAndCheck(reader io.RuneReader) (r rune, ok bool) {
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

func check(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

// msg_key = substr(msg_key_large, 8, 16)
func getMsgKey(authKey *Key, s Side, dataPadded []byte) (key Int128) {
	msgKeyLarge := msgKeyLarge(authKey, s, dataPadded)
	copy(key[:], msgKeyLarge[8:8+16])

	return key
}

// msg_key_large = SHA256(substr(auth_key, 88+x, 32) + plaintext + random_padding);
func msgKeyLarge(authKey *Key, s Side, dataPadded []byte) (hash Int256) {
	x := s.X()

	h := sha256.New()
	h.Write(authKey[88+x : 32+88+x])
	h.Write(dataPadded)
	h.Sum(hash[0:0])

	return hash
}

func genMessageAndAuthKeys() (key Key, msgID Int128) {
	for i := 0; i < 256; i++ {
		key[i] = byte(i)
	}
	for i := 0; i < 16; i++ {
		msgID[i] = byte(i)
	}

	return key, msgID
}
