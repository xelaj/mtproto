package handshake_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"testing"

	"github.com/quenbyako/ext/slices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/xelaj/mtproto/handshake"
)

func TestRSAPad(t *testing.T) {
	a := require.New(t)

	key := parseKeys(`
-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEA6LszBcC1LGzyr992NzE0ieY+BSaOW622Aa9Bd4ZHLl+TuFQ4lo4g
5nKaMBwK/BIb9xUfg0Q29/2mgIR6Zr9krM7HjuIcCzFvDtr+L0GQjae9H0pRB2OO
62cECs5HKhT5DZ98K33vmWiLowc621dQuwKWSQKjWf50XYFw42h21P2KXUGyp2y/
+aEyZ+uVgLLQbRA1dEjSDZ2iGRy12Mk5gpYc397aYp438fsJoHIgJ2lgMv5h7WY9
t6N/byY9Nw9p21Og3AoXSL2q/2IJ1WRUhebgAdGVMlV1fkuOQoEzR7EdpqtQD9Cs
5+bfo3Nhmcyvk5ftB0WkJ9z6bNZ7yxrP8wIDAQAB
-----END RSA PUBLIC KEY-----`)[0]
	data := bytes.Repeat([]byte{'a'}, 144)

	encrypted, err := RSAPad(data, key, nopReader{})
	a.NoError(err)
	a.Len(encrypted, 256)

	hexResult := "bf68719e836806b040cd261ecaf66eb3c4ba19f3bbea3031b2e6cf29167bab647201d101b291dc" +
		"5b716a42e789a38d947fe59e9bcce8f30ef46a946743ea8b6babbce7fc0afc46b802aa453e83471d82a4dfad83f971f35" +
		"0b4b4fb474cd1c48fdf427e4b5fecce9ec3178ae7dac3985856fdefa21d6fdc5e0e0fd8a57bc4f51580d637d372be8d87" +
		"c9aa3fde8e6f8287bcb3be846aadcdd59465375479e248f62ed438f9804fbe36d41ca906243a5f740f3937949aa149ba8" +
		"a8b8e68b3f3e1e3cd3f946387520e21eee55845e1f015a919a22f6a72bfaecd2cae946c91983b41f9ffabe97963bbde8f" +
		"30eaf5fd3c5b8cecab8711bd269e441b6084f385726ff0"
	expected, err := hex.DecodeString(hexResult)
	a.NoError(err)
	a.Equal(expected, encrypted)
}

func TestRSARealLife(t *testing.T) {
	for _, tt := range []struct {
		name       string
		data       []byte
		key        string
		randSource [][]byte
		want       []byte
		wantErr    assert.ErrorAssertionFunc
	}{{
		data: Hexed(`
			955ff5a9081761da96a18b5385000000044b9c1699000000044f2af3cd000000
			655c51687cd8faa965ac85a6e996e5ad8548684598806ca7f6e326f02c638ef4
			8ed9757a8400e9c72f628ce63cf0a5ef5a01aed3aef41cab39a4495c38150c9a
			02000000
		`),
		key: `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEA6LszBcC1LGzyr992NzE0ieY+BSaOW622Aa9Bd4ZHLl+TuFQ4lo4g
5nKaMBwK/BIb9xUfg0Q29/2mgIR6Zr9krM7HjuIcCzFvDtr+L0GQjae9H0pRB2OO
62cECs5HKhT5DZ98K33vmWiLowc621dQuwKWSQKjWf50XYFw42h21P2KXUGyp2y/
+aEyZ+uVgLLQbRA1dEjSDZ2iGRy12Mk5gpYc397aYp438fsJoHIgJ2lgMv5h7WY9
t6N/byY9Nw9p21Og3AoXSL2q/2IJ1WRUhebgAdGVMlV1fkuOQoEzR7EdpqtQD9Cs
5+bfo3Nhmcyvk5ftB0WkJ9z6bNZ7yxrP8wIDAQAB
-----END RSA PUBLIC KEY-----`,
		randSource: [][]byte{
			Hexed(`
			dd36c292116f79bf8ca8c40039fc85aa03540025276b184636a94f0094d7a53d
			ec311731fd44670e80f6750db788c680817c2e16862970d5fcea7c43e53c4180
			7cac41b5e0fccb7f6223ce30aaeea2e1efb9ac8c93e05875cc0e9300785ec587
			9a9c7cd6fb86a1b3c19bcc8ab107a611f97a7cd75f079e76e0bfba19582679aa
			d4cda9e64adc389061f504c85cea97a91c811c2d76f295f2457edc02af19d4a0
			4bf79189c84f4ad020f41f3e79f156d79a11b4cd8ed095e59f1c1fb604f3fcc8
			`), // pad
			Hexed(`
			805289c27e8e599e1f554259ce9cd79891b97cd7219319c3d122c8002538b1c9
			`), // temp aes key
		},
	}} {
		tt := tt // for parallel tests
		tt.wantErr = noErrAsDefault(tt.wantErr)

		t.Run(tt.name, func(t *testing.T) {
			rand := bytes.NewBuffer(slices.Concat(tt.randSource...))

			got, err := RSAPad(tt.data, parseKeys(tt.key)[0], rand)
			if !tt.wantErr(t, err) || err == nil {
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeRSAPad(t *testing.T) {
	a := require.New(t)
	r := rand.Reader

	key, err := rsa.GenerateKey(r, 2048)
	a.NoError(err)
	size := 144

	data := make([]byte, size)
	_, err = io.ReadFull(r, data)
	a.NoError(err)

	encrypted, err := RSAPad(data, &key.PublicKey, r)
	a.NoError(err)
	a.Len(encrypted, 256)

	decrypted, err := DecodeRSAPad(encrypted, key)
	a.NoError(err)
	a.Equal(data, decrypted[:size])
}

func BenchmarkRSAPad(b *testing.B) {
	data := make([]byte, 144)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := RSAPad(data, &testPrivateKey.PublicKey, rand.Reader); err != nil {
			b.Fatal(err)
		}
	}
}

func parseKeys(k string) []*rsa.PublicKey {
	var keys []*rsa.PublicKey

	data := []byte(k)
	for {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}

		key, err := parseRSA(block.Bytes)
		if err != nil {
			panic(fmt.Errorf("parse RSA from PEM: %w", err))
		}

		keys = append(keys, key)
		data = rest
	}

	return keys
}

func parseRSA(data []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKCS1PublicKey(data)
	if err == nil {
		return key, nil
	}
	k, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	publicKey, ok := k.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("parsed unexpected key type %T", k)
	}

	return publicKey, nil
}

type nopReader struct{}

func (nopReader) Read(p []byte) (n int, err error) { clear(p); return len(p), nil }

var testPrivateKey = func() *rsa.PrivateKey {
	rawKey, err := base64.StdEncoding.DecodeString("MIIEowIBAAKCAQEAvPXkA3h/InI+o9Q1B01ysoRZGDlazlSwOMX5Q8KgLsOSKyUDYCI" +
		"07AWV60da4eAUgbI6BNE6B/vw2jH8gInEKb+0DyOKTPGv2t0mPq5a+I+C8xbXIVLwTHBm0mWaiiDQbcaQLBSzxhZ8BTa8VyMK8RO/XIGPoNSnJhf" +
		"LcKg5pmrIenzKDnlPDE2vEPWe8E84cknnmZQVRbbyae/Vqnu+XbadKXuhOro2r1Yz4n49jLHZfUVuyoSbLbYBoTBkdiadO5wCAU9edKl9Bt4LtBu" +
		"SC24MXK4WGWCuX5P7ujENLAMl8Evn5qabD5mMOlFWJUtlBp2iZQhOzJdHtsshMq25pwIDAQABAoIBAEV+zbA1FdTuXXlVZ3dbFY7wO/A7z9jIrtM" +
		"ChK1WHCF2zgBOKZKmof4YA843PQaLqh8VFF+HL6eWEju9XJdNk7ajCa7zrD6mOL3uzc0JxO1bopaS1OYtobELOdWxhoe8j8t/1rBPoNp+lHg6bER" +
		"D4BdP4vY7tD47V4ocADdbt3ArfbfQrEhpYh6kF/bju6PdsjfmkTihG8N8d4CqUSfxr930HFUdNXF72ga7XRG+pBFRAVgZQgNJJkXPx+41WBnFqmp" +
		"Sw44/fT6MeOzy1IoMibDcSZjA/PNSIWoeMxEDKV+6VnkbsiEkwAPotDFzvPm8qROra4JRfGEB+iU3FS08+9ECgYEA54BWA9IAKgeNbzKZkExkq9e" +
		"qrOt0PUA9DrfWZEr2GO+OR7yu7Mi6uhS2ZeUM/3+OTUPfQULZBg2YxPLJ/VuFe/8gPMczT/sZr3arKDgHCDuI/Ft+HQOoEvs+IdrvfC6alfUOnoW" +
		"62QjwfPzEp9UE52yDRsNWQsX4+qJTe6aWmWkCgYEA0PUSB8R1YQx1qWvNxYH9KA73tlFw+WA7GQMbunGEDhgjO4dFscA1YiFLonlqK7WLxqvCtSJ" +
		"an2g1paOQR0V6M5mpDKSeCvLAVhfE1p+z2MPXDx9l7mWRz5z4mJJIXtEqAIn2t7ZOG4MkebcTo3Qq+S92RVnzO1ZpKYS9jOyUyI8CgYA18koZCc7" +
		"P/IKQ7xGp9qNfCBrVwOiNfXK9A0oKhQ1kMi7NuMJqmzwoMLtwczfcMjVO/AoCgzlfl7uJ6an4SGOKyaERiLoEYVdS9Cxeau/4kycQ56Ez0a5Q/gs" +
		"0iHhWT+XmG/0UI8Wu3c5s0dph4doKs9bDnrFzTf7/KOSbY+6kQQKBgBImx8su8LdeerYd7EEU+qXJLxGCX5r6FgglMfpvM/Z5eE4KgS5gsQJ2O/j" +
		"ALU3gtmSqtP5BHrgsOETMQZM/YM8ssPetMSFoVvbjl7DBLMFOudbRdmxQHGt5ikrOokTCTLDBS1JIHt7a9IcyNR2E0NrWmaKKnstvxTDbHBAq2P3" +
		"XAoGBALoQXxKH/gwnri5ioL5LPiHb+SstmSEePS/FcQsuyvgV4a9r5yl+orZVQ0FVaTSKYXSx10Cugja/CqOVop2R7oKLi7HlOKeM4fL2GXID8qp" +
		"SxHZMoDAjdG9Ph1WgU7NI5Sxm70wtDos+vbpmDHvuYHmQ56ljX+5mD3T+ZjuYk7TM")
	if err != nil {
		panic(err)
	}
	k, err := x509.ParsePKCS1PrivateKey(rawKey)
	if err != nil {
		panic(err)
	}

	return k
}()

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

func check(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func noErrAsDefault(e assert.ErrorAssertionFunc) assert.ErrorAssertionFunc {
	if e == nil {
		return assert.NoError
	}

	return e
}
