package keys

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/big"

	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/encoding/tl"
)

// RSAFingerprint вычисляет отпечаток ключа
// т.к. rsa ключ в понятиях MTProto это TL объект, то используется буффер
// подробнее https://core.telegram.org/mtproto/auth_key
func RSAFingerprint(key *rsa.PublicKey) ([]byte, error) {
	dry.PanicIf(key == nil, "key can't be nil")
	exponentAsBigInt := (big.NewInt(0)).SetInt64(int64(key.E))

	buf := bytes.NewBuffer(nil)
	w := tl.NewWriteCursor(buf)
	if err := w.PutMessage(key.N.Bytes()); err != nil {
		return nil, err
	}
	if err := w.PutMessage(exponentAsBigInt.Bytes()); err != nil {
		return nil, err
	}

	fingerprint := dry.Sha1(string(buf.Bytes()))
	return []byte(fingerprint)[12:], nil // последние 8 байт это и есть отпечаток
}

func ReadFromFile(path string) ([]*rsa.PublicKey, error) {
	if !dry.FileExists(path) {
		return nil, errs.NotFound("file", path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file  keys")
	}
	keys := make([]*rsa.PublicKey, 0)
	for {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}

		key, err := pemBytesToRsa(block.Bytes)
		if err != nil {
			const offset = 1 // +1 потому что считаем с 0
			return nil, errors.Wrapf(err, "decoding key №%d", len(keys)+offset)
		}

		keys = append(keys, key)
		data = rest
	}

	return keys, nil
}

func pemBytesToRsa(data []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKCS1PublicKey(data)
	if err == nil {
		return key, nil
	}

	if err.Error() == "x509: failed to parse public key (use ParsePKIXPublicKey instead for this key format)" {
		var k interface{}
		k, err = x509.ParsePKIXPublicKey(data)
		if err == nil {
			return k.(*rsa.PublicKey), nil
		}
	}

	return nil, err
}

func SaveRsaKey(key *rsa.PublicKey) string {
	data := x509.MarshalPKCS1PublicKey(key)
	buf := bytes.NewBufferString("")
	err := pem.Encode(buf, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: data,
	})
	dry.PanicIfErr(err)

	return buf.String()
}
