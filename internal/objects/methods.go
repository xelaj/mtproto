// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package objects

import (
	"context"
	"fmt"

	"github.com/xelaj/tl"
)

type requester interface {
	MakeRequest(context.Context, []byte) ([]byte, error)
}

func request[IN, OUT any](ctx context.Context, m requester, in *IN, out *OUT) error {
	msg, err := tl.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	respRaw, err := m.MakeRequest(ctx, msg)
	if err != nil {
		return fmt.Errorf("sending: %w", err)
	}
	if err := tl.Unmarshal(respRaw, out); err != nil {
		return fmt.Errorf("got invalid response type: %w", err)
	}

	return nil
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

func ReqDHParams(
	ctx context.Context, m requester,
	nonce, serverNonce [16]byte, p, q []byte, publicKeyFingerprint uint64, encryptedData []byte,
) (ServerDHParams, error) {
	var res ServerDHParams
	if err := request(ctx, m, &ReqDHParamsParams{
		Nonce:                nonce,
		ServerNonce:          serverNonce,
		P:                    p,
		Q:                    q,
		PublicKeyFingerprint: publicKeyFingerprint,
		EncryptedData:        encryptedData,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// rpc_drop_answer
// get_future_salts

type PingParams struct {
	PingID int64
}

func (*PingParams) CRC() uint32 { return 0x7abe77ec }

type PingDelayDisconnectParams struct {
	PingID          int64
	DisconnectDelay int32
}

func (*PingDelayDisconnectParams) CRC() uint32 { return 0xf3427b8c }

// ping_delay_disconnect
// destroy_session
// http_wait

// set_client_DH_params#f5045f1f nonce:int128 server_nonce:int128 encrypted_data:bytes = Set_client_DH_params_answer;

// rpc_drop_answer#58e4a740 req_msg_id:long = RpcDropAnswer;
// get_future_salts#b921bd04 num:int = FutureSalts;
// ping_delay_disconnect#f3427b8c ping_id:long disconnect_delay:int = Pong;
// destroy_session#e7512126 session_id:long = DestroySessionRes;

// http_wait#9299359f max_delay:int wait_after:int max_wait:int = HttpWait;

func must[T any](t T, err error) T { check(err); return t }
func check(err error) {
	if err != nil {
		panic(err)
	}
}
