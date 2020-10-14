package mtproto

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/serialize"
)

type ReqPQParams struct {
	Nonce *serialize.Int128
}

func (_ *ReqPQParams) CRC() uint32 {
	return 0x60469778
}

func (t *ReqPQParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	return buf.Result()
}

func (t *ReqPQParams) DecodeFrom(d *serialize.Decoder) {
	t.Nonce = d.PopInt128()
}

func (m *MTProto) ReqPQ(nonce *serialize.Int128) (*serialize.ResPQ, error) {
	data, err := m.MakeRequest(&ReqPQParams{Nonce: nonce})
	if err != nil {
		return nil, errors.Wrap(err, "sending ReqPQ")
	}

	resp, ok := data.(*serialize.ResPQ)
	if !ok {
		return nil, errors.New("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}

type ReqDHParamsParams struct {
	Nonce                *serialize.Int128
	ServerNonce          *serialize.Int128
	P                    []byte
	Q                    []byte
	PublicKeyFingerprint int64
	EncryptedData        []byte
}

func (_ *ReqDHParamsParams) CRC() uint32 {
	return 0xd712e4be
}

func (t *ReqDHParamsParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutMessage(t.P)
	buf.PutMessage(t.Q)
	buf.PutLong(t.PublicKeyFingerprint)
	buf.PutMessage(t.EncryptedData)
	return buf.Result()
}

func (t *ReqDHParamsParams) DecodeFrom(d *serialize.Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.P = d.PopMessage()
	t.Q = d.PopMessage()
	t.PublicKeyFingerprint = d.PopLong()
	t.EncryptedData = d.PopMessage()
}

func (m *MTProto) ReqDHParams(nonce, serverNonce *serialize.Int128, p, q []byte, publicKeyFingerprint int64, encryptedData []byte) (serialize.ServerDHParams, error) {
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

	resp, ok := data.(serialize.ServerDHParams)
	if !ok {
		return nil, errors.New("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}

type SetClientDHParamsParams struct {
	Nonce         *serialize.Int128
	ServerNonce   *serialize.Int128
	EncryptedData []byte
}

func (_ *SetClientDHParamsParams) CRC() uint32 {
	return 0xf5045f1f
}

func (t *SetClientDHParamsParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutMessage(t.EncryptedData)
	return buf.Result()
}

func (t *SetClientDHParamsParams) DecodeFrom(d *serialize.Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.EncryptedData = d.PopMessage()
}

func (m *MTProto) SetClientDHParams(nonce, serverNonce *serialize.Int128, encryptedData []byte) (serialize.SetClientDHParamsAnswer, error) {
	data, err := m.MakeRequest(&SetClientDHParamsParams{
		Nonce:         nonce,
		ServerNonce:   serverNonce,
		EncryptedData: encryptedData,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending Ping")
	}

	resp, ok := data.(serialize.SetClientDHParamsAnswer)
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

func (_ *PingParams) CRC() uint32 {
	return 0x7abe77ec
}

func (t *PingParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutLong(t.PingID)
	return buf.Result()
}

func (t *PingParams) DecodeFrom(d *serialize.Decoder) {
	t.PingID = d.PopLong()
}

func (m *MTProto) Ping(pingID int64) (*serialize.Pong, error) {
	data, err := m.MakeRequest(&PingParams{
		PingID: pingID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending Ping")
	}

	resp, ok := data.(*serialize.Pong)
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
