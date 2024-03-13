package handshake

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xelaj/tl"

	"github.com/xelaj/mtproto/v2/internal/payload"
	"github.com/xelaj/mtproto/v2/internal/transport"
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

func reqPQMulti(ctx context.Context, t transport.Transport, nonce Int128) (*ResPQ, error) {
	var res ResPQ
	if err := request(ctx, t, &ReqPQMultiParams{Nonce: nonce}, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type ReqPQMultiParams struct {
	Nonce Int128
}

func (*ReqPQMultiParams) CRC() uint32 { return 0xbe7e8ef1 }

type ResPQ struct {
	Nonce        Int128
	ServerNonce  Int128
	Pq           []byte
	Fingerprints []uint64
}

func (*ResPQ) CRC() uint32 { return 0x05162463 }

func reqDHParams(ctx context.Context, t transport.Transport, nonce, serverNonce Int128, p, q []byte, publicKeyFingerprint uint64, encryptedData []byte) (*ServerDHParamsOk, error) {
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
	Nonce                Int128
	ServerNonce          Int128
	P                    []byte
	Q                    []byte
	PublicKeyFingerprint uint64
	EncryptedData        []byte
}

func (*ReqDHParamsParams) CRC() uint32 { return 0xd712e4be }

type ServerDHParamsOk struct {
	Nonce           Int128
	ServerNonce     Int128
	EncryptedAnswer []byte
}

func (*ServerDHParamsOk) _ServerDHParams() {}
func (*ServerDHParamsOk) CRC() uint32      { return 0xd0e8075c }

func setClientDHParams(ctx context.Context, t transport.Transport, nonce, serverNonce Int128, encryptedData []byte) (SetClientDHParamsAnswer, error) {
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
	Nonce         Int128
	ServerNonce   Int128
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
	Nonce        Int128
	ServerNonce  Int128
	NewNonceHash Int128
}

func (*DHGenOk) _SetClientDHParamsAnswer() {}
func (*DHGenOk) CRC() uint32               { return 0x3bcbf734 }

type DHGenRetry struct {
	Nonce        Int128
	ServerNonce  Int128
	NewNonceHash Int128
}

func (*DHGenRetry) _SetClientDHParamsAnswer() {}
func (*DHGenRetry) CRC() uint32               { return 0x46dc1fb9 }

type DHGenFail struct {
	Nonce        Int128
	ServerNonce  Int128
	NewNonceHash Int128
}

func (*DHGenFail) _SetClientDHParamsAnswer() {}
func (*DHGenFail) CRC() uint32               { return 0xa69dae02 }
