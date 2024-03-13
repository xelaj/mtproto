package handshake

import "io"

func randInt128(r io.Reader) Int128 {
	var nonce Int128
	if _, err := r.Read(nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func randInt256(r io.Reader) Int256 {
	var nonce Int256
	if _, err := r.Read(nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}
