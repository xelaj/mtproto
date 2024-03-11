package handshake

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xelaj/mtproto/internal/payload"
	"github.com/xelaj/mtproto/internal/transport"
	"github.com/xelaj/tl"
)

func request[IN, OUT any](ctx context.Context, t transport.Transport, in *IN, out *OUT) error {
	msg, err := tl.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	_, err = t.WriteMsg(ctx, msg, payload.InitiatorClient)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	_, resp, _, err := t.ReadMsg(ctx)
	if err != nil {
		return fmt.Errorf("receiving response: %w", err)
	}

	if err := tl.NewDecoder(bytes.NewBuffer(resp)).SetRegistry(Registry).Decode(out); err != nil {
		return fmt.Errorf("got invalid response type: %w", err)
	}

	return nil
}

var Registry = NewRegistry()

func NewRegistry() *tl.ObjectRegistry {
	r := tl.NewRegistry()
	tl.RegisterObject[*ResPQ](r)
	tl.RegisterObject[*ServerDHParamsOk](r)
	tl.RegisterObject[*ServerDHInnerData](r)
	tl.RegisterObject[*DHGenOk](r)
	tl.RegisterObject[*DHGenRetry](r)
	tl.RegisterObject[*DHGenFail](r)

	return r
}

// func reqPQ(ctx context.Context, t transport.Transport, nonce [16]byte) (*ResPQ, error) {
// 	var res ResPQ
// 	if err := request(ctx, t, &ReqPQParams{Nonce: nonce}, &res); err != nil {
// 		return nil, err
// 	}
//
// 	return &res, nil
// }
//
// type ReqPQParams struct {
// 	Nonce [16]byte
// }
//
// func (*ReqPQParams) CRC() uint32 { return 0x60469778 }

func reqPQMulti(ctx context.Context, t transport.Transport, nonce [16]byte) (*ResPQ, error) {
	var res ResPQ
	if err := request(ctx, t, &ReqPQMultiParams{Nonce: nonce}, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type ReqPQMultiParams struct {
	Nonce [16]byte
}

func (*ReqPQMultiParams) CRC() uint32 { return 0xbe7e8ef1 }

type ResPQ struct {
	Nonce        [16]byte
	ServerNonce  [16]byte
	Pq           []byte
	Fingerprints []uint64
}

func (*ResPQ) CRC() uint32 { return 0x05162463 }

func reqDHParams(ctx context.Context, t transport.Transport, nonce, serverNonce [16]byte, p, q []byte, publicKeyFingerprint uint64, encryptedData []byte) (*ServerDHParamsOk, error) {
	var res ServerDHParamsOk
	if err := request(ctx, t, &ReqDHParamsParams{
		Nonce:                nonce,
		ServerNonce:          serverNonce,
		P:                    p,
		Q:                    q,
		PublicKeyFingerprint: publicKeyFingerprint,
		EncryptedData:        encryptedData,
	}, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type ReqDHParamsParams struct {
	Nonce                [16]byte
	ServerNonce          [16]byte
	P                    []byte
	Q                    []byte
	PublicKeyFingerprint uint64
	EncryptedData        []byte
}

func (*ReqDHParamsParams) CRC() uint32 { return 0xd712e4be }

type ServerDHParamsOk struct {
	Nonce           [16]byte
	ServerNonce     [16]byte
	EncryptedAnswer []byte
}

func (*ServerDHParamsOk) _ServerDHParams() {}
func (*ServerDHParamsOk) CRC() uint32      { return 0xd0e8075c }

func setClientDHParams(ctx context.Context, t transport.Transport, nonce, serverNonce [16]byte, encryptedData []byte) (SetClientDHParamsAnswer, error) {
	var res SetClientDHParamsAnswer
	if err := request(ctx, t, &SetClientDHParamsParams{
		Nonce:         nonce,
		ServerNonce:   serverNonce,
		EncryptedData: encryptedData,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

type SetClientDHParamsParams struct {
	Nonce         [16]byte
	ServerNonce   [16]byte
	EncryptedData []byte
}

func (*SetClientDHParamsParams) CRC() uint32 { return 0xf5045f1f }

type SetClientDHParamsAnswer interface {
	tl.Object
	_SetClientDHParamsAnswer()
}

var (
	_ SetClientDHParamsAnswer = (*DHGenOk)(nil)
	_ SetClientDHParamsAnswer = (*DHGenRetry)(nil)
	_ SetClientDHParamsAnswer = (*DHGenFail)(nil)
)

type DHGenOk struct {
	Nonce        [16]byte
	ServerNonce  [16]byte
	NewNonceHash [16]byte
}

func (*DHGenOk) _SetClientDHParamsAnswer() {}
func (*DHGenOk) CRC() uint32               { return 0x3bcbf734 }

type DHGenRetry struct {
	Nonce        [16]byte
	ServerNonce  [16]byte
	NewNonceHash [16]byte
}

func (*DHGenRetry) _SetClientDHParamsAnswer() {}
func (*DHGenRetry) CRC() uint32               { return 0x46dc1fb9 }

type DHGenFail struct {
	Nonce        [16]byte
	ServerNonce  [16]byte
	NewNonceHash [16]byte
}

func (*DHGenFail) _SetClientDHParamsAnswer() {}
func (*DHGenFail) CRC() uint32               { return 0xa69dae02 }
