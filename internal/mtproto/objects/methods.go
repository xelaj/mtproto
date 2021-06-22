// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package objects

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/internal/encoding/tl"
)

type requester interface {
	MakeRequest(tl.Object) (any, error)
}

type ReqPQParams struct {
	Nonce *tl.Int128
}

func (*ReqPQParams) CRC() uint32 {
	return 0x60469778 //nolint:gomnd not magic
}

func ReqPQ(m requester, nonce *tl.Int128) (*ResPQ, error) {
	data, err := m.MakeRequest(&ReqPQParams{Nonce: nonce})
	if err != nil {
		return nil, errors.Wrap(err, "sending ReqPQ")
	}

	resp, ok := data.(*ResPQ)
	if !ok {
		return nil, errors.New("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}

type ReqDHParamsParams struct {
	Nonce                *tl.Int128
	ServerNonce          *tl.Int128
	P                    []byte
	Q                    []byte
	PublicKeyFingerprint int64
	EncryptedData        []byte
}

func (*ReqDHParamsParams) CRC() uint32 {
	return 0xd712e4be //nolint:gomnd not magic
}

func ReqDHParams(
	m requester,
	nonce, serverNonce *tl.Int128, p, q []byte, publicKeyFingerprint int64, encryptedData []byte,
) (ServerDHParams, error) {
	data, err := m.MakeRequest(&ReqDHParamsParams{
		Nonce:                nonce,
		ServerNonce:          serverNonce,
		P:                    p,
		Q:                    q,
		PublicKeyFingerprint: publicKeyFingerprint,
		EncryptedData:        encryptedData,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending ReqDHParams")
	}

	resp, ok := data.(ServerDHParams)
	if !ok {
		return nil, errors.New("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}

type SetClientDHParamsParams struct {
	Nonce         *tl.Int128
	ServerNonce   *tl.Int128
	EncryptedData []byte
}

func (*SetClientDHParamsParams) CRC() uint32 {
	return 0xf5045f1f //nolint:gomnd not magic
}

func SetClientDHParams(m requester, nonce, serverNonce *tl.Int128, encryptedData []byte) (SetClientDHParamsAnswer, error) {
	data, err := m.MakeRequest(&SetClientDHParamsParams{
		Nonce:         nonce,
		ServerNonce:   serverNonce,
		EncryptedData: encryptedData,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending Ping")
	}

	resp, ok := data.(SetClientDHParamsAnswer)
	if !ok {
		return nil, errors.New("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}

// rpc_drop_answer
// get_future_salts

type PingParams struct {
	PingID int64
}

func (*PingParams) CRC() uint32 {
	return 0x7abe77ec //nolint:gomnd not magic
}

func Ping(m requester, pingID int64) (*Pong, error) {
	data, err := m.MakeRequest(&PingParams{
		PingID: pingID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending Ping")
	}

	resp, ok := data.(*Pong)
	if !ok {
		return nil, errors.New("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}

// ping_delay_disconnect
// destroy_session
// http_wait

// set_client_DH_params#f5045f1f nonce:int128 server_nonce:int128 encrypted_data:bytes = Set_client_DH_params_answer;

// rpc_drop_answer#58e4a740 req_msg_id:long = RpcDropAnswer;
// get_future_salts#b921bd04 num:int = FutureSalts;
// ping_delay_disconnect#f3427b8c ping_id:long disconnect_delay:int = Pong;
// destroy_session#e7512126 session_id:long = DestroySessionRes;

// http_wait#9299359f max_delay:int wait_after:int max_wait:int = HttpWait;
