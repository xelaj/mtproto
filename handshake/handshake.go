package handshake

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/go-faster/xor"
	"github.com/xelaj/mtproto/internal/transport"
	"github.com/xelaj/tl"
)

// `expiration` param defines how long session will be stored on server side.
// This value guarantees *at most** this duration, but not exactly. or even
// more. According to the documentation, "The server is free to discard its copy
// earlier". Zero value means permanent session, so if keys will leak, someone
// might use them to use authorized session.
func Perform(ctx context.Context, t transport.Transport, keys []*rsa.PublicKey, dc int, expiration time.Duration) ([256]byte, uint64, error) {
	return perform(ctx, t, rand.Reader, keys, dc, expiration)
}

func perform(ctx context.Context, t transport.Transport, rand io.Reader, keys []*rsa.PublicKey, dc int, expiration time.Duration) (authKey [256]byte, salt uint64, err error) {
	// 2. Server sends response of server nonce, fingerprints and PQ
	nonce, serverNonce, pq, key, err := initHandshake(ctx, t, rand, keys)
	if err != nil {
		return authKey, salt, fmt.Errorf("init handshake: %w", err)
	}

	authKey, salt, err = makeProofOfWork(ctx, t, rand, key, pq, nonce, serverNonce, dc, expiration)
	if err != nil {
		return authKey, salt, fmt.Errorf("making PoW: %w", err)
	}

	return authKey, salt, nil
}

func initHandshake(
	ctx context.Context,
	t transport.Transport,
	rand io.Reader,
	keys []*rsa.PublicKey,
) (
	nonce [16]byte,
	serverNonce [16]byte,
	pq uint64,
	key *rsa.PublicKey,
	err error,
) {
	nonce = randInt128(rand)

	resp, err := reqPQMulti(ctx, t, nonce)
	if err != nil {
		return [16]byte{}, [16]byte{}, 0, nil, err
	}

	if resp.Nonce != nonce {
		return [16]byte{}, [16]byte{}, 0, nil, errors.New("server returned wrong client nonce")
	}

	pqRaw := big.NewInt(0).SetBytes(resp.Pq)
	if !pqRaw.IsUint64() {
		// Normally pq is less than or equal to 2^63-1, sometimes it **might**
		// be bigger, but in real world no one was reported about that. In
		// addition, pq between 2^32 and 2^64 is big enough to make secure
		// handshake.
		return [16]byte{}, [16]byte{}, 0, nil, errors.New("server returned wrong pq")
	}

	i := lookupForKeys(keys, resp.Fingerprints)
	if i < 0 {
		return [16]byte{}, [16]byte{}, 0, nil, errors.New("server returned unknown key fingerprints")
	}

	return nonce, resp.ServerNonce, pqRaw.Uint64(), keys[i], nil
}

func makeProofOfWork(
	ctx context.Context,
	t transport.Transport,
	rand io.Reader,
	key *rsa.PublicKey,
	pq uint64,
	nonce [16]byte,
	serverNonce [16]byte,
	dc int,
	expiration time.Duration,
) (
	authKey [256]byte,
	salt uint64,
	err error,
) {
	p, q := DecomposePQ(pq, rand)

	newNonce := randInt256(rand)

	envelope := serializePQInnerData(pq, p, q, nonce, serverNonce, newNonce, dc, expiration)

	encrypted, err := RSAPad(envelope, key, rand)
	if err != nil {
		return [256]byte{}, 0, err
	}

	pBytes := big.NewInt(0).SetUint64(p).Bytes()
	qBytes := big.NewInt(0).SetUint64(q).Bytes()

	resp, err := reqDHParams(ctx, t, nonce, serverNonce, pBytes, qBytes, rsaFingerprint(key), encrypted)
	if err != nil {
		return [256]byte{}, 0, err
	}

	if resp.Nonce != nonce {
		return [256]byte{}, 0, errors.New("server returned wrong client nonce")
	}

	if resp.ServerNonce != serverNonce {
		return [256]byte{}, 0, errors.New("server returned wrong server nonce")
	}

	tempKey, tempIV := TempAESKeys(newNonce, serverNonce)

	decrypted, err := decryptHandshake(tempKey, tempIV, resp.EncryptedAnswer)

	var inner ServerDHInnerData
	if err := tl.NewDecoder(bytes.NewBuffer(decrypted)).SetRegistry(Registry).Decode(&inner); err != nil {
		return [256]byte{}, 0, fmt.Errorf("got invalid response type: %w", err)
	}

	dhPrime := big.NewInt(0).SetBytes(inner.DhPrime)
	gA := big.NewInt(0).SetBytes(inner.GA)

	gB, b, err := MakeGAB(rand, big.NewInt(int64(inner.G)), gA, dhPrime)
	if err != nil {
		return [256]byte{}, 0, fmt.Errorf("make gab: %w", err)
	}

	clientDH, err := tl.Marshal(&ClientDHInnerData{
		ServerNonce: serverNonce,
		Nonce:       nonce,
		Retry:       0, // first attempt
		GB:          gB.Bytes(),
	})
	if err != nil {
		panic(err)
	}

	if encrypted, err = encryptHandshake(rand, tempKey, tempIV, clientDH); err != nil {
		panic(err)
	}

	dhGen, err := setClientDHParams(ctx, t, nonce, serverNonce, encrypted)
	if err != nil {
		return [256]byte{}, 0, err
	}

	dhGenOk, ok := dhGen.(*DHGenOk)
	if !ok {
		return [256]byte{}, 0, errors.New("server failed to generate DH")
	}

	if dhGenOk.Nonce != nonce {
		return [256]byte{}, 0, errors.New("server returned wrong client nonce")
	}

	if dhGenOk.ServerNonce != serverNonce {
		return [256]byte{}, 0, errors.New("server returned wrong server nonce")
	}

	// 7. Computing auth_key using formula (g_a)^b mod dh_prime
	big.NewInt(0).Exp(gA, b, dhPrime).FillBytes(authKey[:])

	if dhGenOk.NewNonceHash != nonceHash(newNonce, authKey) {
		return [256]byte{}, 0, errors.New("server returned wrong newNonce hash")
	}

	return authKey, ServerSalt(newNonce, serverNonce), nil
}

func ServerSalt(newNonce [32]byte, serverNonce [16]byte) (salt uint64) {
	var serverSalt [8]byte
	copy(serverSalt[:], newNonce[:8])
	xor.Bytes(serverSalt[:], serverSalt[:], serverNonce[:8])
	return binary.LittleEndian.Uint64(serverSalt[:])
}

// nonceHash computes nonce_hash_1.
// See https://core.telegram.org/mtproto/auth_key#dh-key-exchange-complete.
func nonceHash(newNonce [32]byte, key [256]byte) (r [16]byte) {
	var buf []byte
	buf = append(buf, newNonce[:]...)
	buf = append(buf, 1)
	buf = append(buf, sha(key[:])[0:8]...)
	buf = sha(buf)[4:20]
	copy(r[:], buf)
	return
}

func sha(v []byte) []byte { h := sha1.Sum(v); return h[:] }
