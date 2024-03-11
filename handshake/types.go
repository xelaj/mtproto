package handshake

import (
	"math/big"
	"time"

	"github.com/xelaj/tl"
)

func serializePQInnerData(pq, p, q uint64, nonce, serverNonce [16]byte, newNonce [32]byte, dc int, expiration time.Duration) (res []byte) {
	//! IMPORTANT: values here are big endian, unlike in whole mtproto
	//  protocol. This is described in handshake documentation.
	//
	// See more https://core.telegram.org/mtproto/auth_key (search by "endian"
	// keyword).

	var err error
	if expiration != 0 {
		if res, err = tl.Marshal(&PQInnerDataTempDC{
			// ! IMPORTANT: see comment above describing reason of using big
			// endian
			Pq:          big.NewInt(0).SetUint64(pq).Bytes(),
			P:           big.NewInt(0).SetUint64(p).Bytes(),
			Q:           big.NewInt(0).SetUint64(q).Bytes(),
			Nonce:       nonce,
			ServerNonce: serverNonce,
			NewNonce:    newNonce,
			DC:          int32(dc),
			ExpiresIn:   int32(expiration.Seconds()),
		}); err != nil {
			panic(err)
		}

		return res
	}

	if res, err = tl.Marshal(&PQInnerDataDC{
		// ! IMPORTANT: see comment above describing reason of using big endian
		Pq:          big.NewInt(0).SetUint64(pq).Bytes(),
		P:           big.NewInt(0).SetUint64(p).Bytes(),
		Q:           big.NewInt(0).SetUint64(q).Bytes(),
		Nonce:       nonce,
		ServerNonce: serverNonce,
		NewNonce:    newNonce,
		// DC:          int32(dc),
	}); err != nil {
		panic(err)
	}

	return res
}

type PQInnerData interface {
	tl.Object
	_PQInnerData()
}

type PQInnerDataObj struct {
	Pq          []byte
	P           []byte
	Q           []byte
	Nonce       [16]byte
	ServerNonce [16]byte
	NewNonce    [32]byte
}

func (*PQInnerDataObj) CRC() uint32 { return 0x83c95aec }

type PQInnerDataDC struct {
	Pq          []byte
	P           []byte
	Q           []byte
	Nonce       [16]byte
	ServerNonce [16]byte
	NewNonce    [32]byte
	DC          int32
}

func (*PQInnerDataDC) _PQInnerData() {}
func (*PQInnerDataDC) CRC() uint32   { return 0xa9f55f95 }

type PQInnerDataTempDC struct {
	Pq          []byte
	P           []byte
	Q           []byte
	Nonce       [16]byte
	ServerNonce [16]byte
	NewNonce    [32]byte
	DC          int32
	ExpiresIn   int32
}

func (*PQInnerDataTempDC) _PQInnerData() {}
func (*PQInnerDataTempDC) CRC() uint32   { return 0x56fddf88 }

type ServerDHInnerData struct {
	Nonce       [16]byte
	ServerNonce [16]byte
	G           int32
	DhPrime     []byte
	GA          []byte
	ServerTime  int32
}

func (*ServerDHInnerData) CRC() uint32 { return 0xb5890dba }

type ClientDHInnerData struct {
	Nonce       [16]byte
	ServerNonce [16]byte
	Retry       int64
	GB          []byte
}

func (*ClientDHInnerData) CRC() uint32 { return 0x6643b654 }