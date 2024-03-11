package handshake

import "io"

func randInt128(r io.Reader) [16]byte {
	var nonce [16]byte
	if _, err := r.Read(nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func randInt256(r io.Reader) [32]byte {
	var nonce [32]byte
	if _, err := r.Read(nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}
