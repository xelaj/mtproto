// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package objects

// some types are decoding VEEEEEEERY specific way, so it stored here and only here.

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
	"github.com/xelaj/tl"

	"github.com/xelaj/mtproto/internal/payload"
)

// TYPES

func Unmarshal(data []byte, res any) error {
	return tl.NewDecoder(bytes.NewBuffer(data)).SetRegistry(Registry).Decode(res)
}

var Registry = NewRegistry()

func NewRegistry() *tl.ObjectRegistry {
	r := tl.NewRegistry()
	tl.RegisterObject[*ReqDHParamsParams](r)
	tl.RegisterObject[*PingParams](r)
	tl.RegisterObject[*ServerDHParamsFail](r)
	tl.RegisterObject[*ServerDHParamsOk](r)
	tl.RegisterObject[*ClientDHInnerData](r)
	tl.RegisterObject[*DHGenOk](r)
	tl.RegisterObject[*DHGenRetry](r)
	tl.RegisterObject[*DHGenFail](r)
	tl.RegisterObject[*RpcResult](r)
	tl.RegisterObject[*RpcError](r)
	tl.RegisterObject[*RpcAnswerUnknown](r)
	tl.RegisterObject[*RpcAnswerDroppedRunning](r)
	tl.RegisterObject[*RpcAnswerDropped](r)
	tl.RegisterObject[*FutureSalt](r)
	tl.RegisterObject[*FutureSalts](r)
	tl.RegisterObject[*Pong](r)
	tl.RegisterObject[*NewSessionCreated](r)
	tl.RegisterObject[*MessageContainer](r)
	tl.RegisterObject[*MsgsAck](r)
	tl.RegisterObject[*BadMsgNotification](r)
	tl.RegisterObject[*BadServerSalt](r)
	tl.RegisterObject[*MsgResendReq](r)
	tl.RegisterObject[*MsgsStateReq](r)
	tl.RegisterObject[*MsgsStateInfo](r)
	tl.RegisterObject[*MsgsAllInfo](r)
	tl.RegisterObject[*MsgsDetailedInfo](r)
	tl.RegisterObject[*MsgsNewDetailedInfo](r)

	return r
}

type ServerDHParams interface {
	tl.Object
	_ServerDHParams()
}

type ServerDHParamsFail struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	NewNonceHash *tl.Int128
}

func (*ServerDHParamsFail) _ServerDHParams() {}
func (*ServerDHParamsFail) CRC() uint32      { return 0x79cb045d }

type ServerDHParamsOk struct {
	Nonce           [16]byte
	ServerNonce     [16]byte
	EncryptedAnswer []byte
}

func (*ServerDHParamsOk) _ServerDHParams() {}
func (*ServerDHParamsOk) CRC() uint32      { return 0xd0e8075c }

type ClientDHInnerData struct {
	Nonce       [16]byte
	ServerNonce [16]byte
	Retry       int64
	GB          []byte
}

func (*ClientDHInnerData) CRC() uint32 { return 0x6643b654 }

type DHGenOk struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	NewNonceHash *tl.Int128
}

func (*DHGenOk) _SetClientDHParamsAnswer() {}
func (*DHGenOk) CRC() uint32               { return 0x3bcbf734 }

type SetClientDHParamsAnswer interface {
	tl.Object
	_SetClientDHParamsAnswer()
}

type DHGenRetry struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	NewNonceHash *tl.Int128
}

func (*DHGenRetry) _SetClientDHParamsAnswer() {}
func (*DHGenRetry) CRC() uint32               { return 0x46dc1fb9 }

type DHGenFail struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	NewNonceHash *tl.Int128
}

func (*DHGenFail) _SetClientDHParamsAnswer() {}
func (*DHGenFail) CRC() uint32               { return 0xa69dae02 }

type RpcResult struct {
	ReqMsgID int64
	Obj      []byte // contains tl object
}

func (*RpcResult) CRC() uint32 { return CrcRpcResult }

func (*RpcResult) MarshalTL(e tl.MarshalState) error { panic("forbidden") }

func (t *RpcResult) UnmarshalTL(d tl.UnmarshalState) (err error) {
	if t.ReqMsgID, err = d.PopLong(); err != nil {
		return err
	} else if t.Obj, err = d.ReadAll(); err != nil {
		return err
	}

	return nil
}

// DecodeFromButItsVector
// декодирует ТАК ЖЕ как DecodeFrom, но за тем исключением, что достает не объект, а слайс.
// проблема в том, что вектор (слайс) в понятиях MTProto это как-бы объект, но вот как бы и нет
// технически, эта функция — костыль, т.к. нет никакого внятного способа передать декодеру
// информацию, что нужно доставать вектор (ведь RPC Result это всегда объекты, но вектор тоже
// объект, кто бы мог подумать)
// другими словами:
// т.к. telegram отсылает на реквесты сообщения (messages, TL в рамках этого пакета)
// НО! иногда на некоторые запросы приходят ответы в виде вектора. Просто потому что.
// поэтому этот кусочек возвращает корявое апи к его же описанию — ответы это всегда объекты.
//func (t *RpcResult) DecodeFromButItsVector(d *Decoder, as reflect.Type) {
//	t.ReqMsgID = d.PopLong()
//	crc := binary.LittleEndian.Uint32(d.GetRestOfMessage()[:WordLen])
//	if crc == CrcGzipPacked {
//		_ = d.PopCRC()
//		gz := &GzipPacked{}
//		gz.DecodeFromButItsVector(d, as)
//		t.Obj = gz.Obj.(*InnerVectorObject)
//	} else {
//		vector := d.PopVector(as)
//		t.Obj = &InnerVectorObject{I: vector}
//	}
//}

type RpcError struct {
	ErrorCode    int32
	ErrorMessage string
}

func (*RpcError) CRC() uint32 { return 0x2144ca19 }

type RpcDropAnswer interface {
	tl.Object
	_RpcDropAnswer()
}

type RpcAnswerUnknown null

func (*RpcAnswerUnknown) _RpcDropAnswer() {}
func (*RpcAnswerUnknown) CRC() uint32     { return 0x5e2ad36e }

type RpcAnswerDroppedRunning null

func (*RpcAnswerDroppedRunning) _RpcDropAnswer() {}
func (*RpcAnswerDroppedRunning) CRC() uint32     { return 0xcd78e586 }

type RpcAnswerDropped struct {
	MsgID int64
	SewNo int32
	Bytes int32
}

func (*RpcAnswerDropped) _RpcDropAnswer() {}
func (*RpcAnswerDropped) CRC() uint32     { return 0xa43ad8b7 }

type FutureSalt struct {
	ValidSince int32
	ValidUntil int32
	Salt       int64
}

func (*FutureSalt) CRC() uint32 { return 0x0949d9dc }

type FutureSalts struct {
	ReqMsgID int64
	Now      int32
	Salts    []*FutureSalt
}

func (*FutureSalts) CRC() uint32 { return 0xae500895 }

type Pong struct {
	MsgID  int64
	PingID int64
}

func (*Pong) CRC() uint32 { return 0x347773c5 }

// destroy_session_ok#e22045fc session_id:long = DestroySessionRes;
// destroy_session_none#62d350c9 session_id:long = DestroySessionRes;

type NewSessionCreated struct {
	FirstMsgID int64
	UniqueID   int64
	ServerSalt int64
}

func (*NewSessionCreated) CRC() uint32 { return 0x9ec20908 }

// ! исключение из правил: это оказывается почти-вектор, т.к.
//
//	записан как `msg_container#73f1f8dc messages:vector<%Message> = MessageContainer;`
//	судя по всему, <%Type> означает, что может это неявный вектор???
//
// ! Perhaps the devs have gone nuts at this point, I don't know really.
type MessageContainer struct {
	Content []*payload.Encrypted
}

func (*MessageContainer) CRC() uint32 { return 0x73f1f8dc }

func (t *MessageContainer) MarshalTL(e tl.MarshalState) (err error) {
	panic("forbidden")
}

func (t *MessageContainer) UnmarshalTL(d tl.UnmarshalState) error {
	count, err := d.PopInt()
	if err != nil {
		return err
	}
	arr := make([]*payload.Encrypted, count)
	for i := int32(0); i < count; i++ {
		msg := new(payload.Encrypted)
		id, err := d.PopLong()
		if err != nil {
			return err
		}
		msg.ID = payload.MsgID(id)
		msg.SeqNo, err = d.PopCRC()
		if err != nil {
			return err
		}
		size, err := d.PopInt()
		if err != nil {
			return err
		}

		msg.Msg = make([]byte, size)
		if _, err = d.Read(msg.Msg); err != nil {
			return err
		}
		arr[i] = msg
	}
	t.Content = arr

	return nil
}

// type Message struct {
// 	MsgID int64
// 	SeqNo int32
// 	Bytes int32
// 	Body  tl.Object
// }
//
// type MsgCopy struct {
// 	OrigMessage *Message
// }
//
// func (*MsgCopy) CRC() uint32 { return 0xe06046b2 }

// type GzipPacked struct {
// 	Orig []byte
// 	Obj  []byte
// }

// func (*GzipPacked) CRC() uint32 { return CrcGzipPacked }

// func (*GzipPacked) MarshalTL(e tl.MarshalState) error { panic("not implemented") }

// func (t *GzipPacked) UnmarshalTL(d tl.UnmarshalState) (err error) {
// 	if t.Orig, err = t.popMessageAsBytes(d); err != nil {
// 		return err
// 	} else if err = tl.Unmarshal(t.Orig, &t.Obj); err != nil {
// 		return errors.Wrap(err, "parsing gzipped object")
// 	}

// 	return nil
// }

// func (*GzipPacked) popMessageAsBytes(d tl.UnmarshalState) ([]byte, error) {
// 	data, err := io.ReadAll(d)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	return UnzipObject(data, false)
// }

func UnzipObject(in []byte, readCRC bool) ([]byte, error) {
	if readCRC {
		if len(in) < tl.WordLen {
			return nil, errors.New("invalid data")
		}
		crc := binary.LittleEndian.Uint32(in)
		if crc != CrcGzipPacked {
			return nil, errors.New("not a gzip")
		}
		in = in[tl.WordLen:]
	}

	// reading via tl cause input is typelang byte object with length and padding
	var b []byte
	if err := tl.Unmarshal(in, &b); err != nil {
		return nil, errors.Wrap(err, "reading byte array from message")
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, errors.Wrap(err, "creating gzip reader")
	}

	return io.ReadAll(gz)
}

type MsgsAck struct {
	MsgIDs []payload.MsgID
}

func (*MsgsAck) CRC() uint32 { return 0x62d6b459 }

type BadMsgNotification struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	Code        int32
}

func (*BadMsgNotification) _BadMsgNotification() {}
func (*BadMsgNotification) CRC() uint32          { return 0xa7eff811 }

type BadServerSalt struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	ErrorCode   int32
	NewSalt     int64
}

func (*BadServerSalt) _BadMsgNotification() {}
func (*BadServerSalt) CRC() uint32          { return 0xedab447b }

// msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;

type MsgResendReq struct {
	MsgIDs []int64
}

func (*MsgResendReq) CRC() uint32 { return 0x7d861a08 }

type MsgsStateReq struct {
	MsgIDs []int64
}

func (*MsgsStateReq) CRC() uint32 { return 0xda69fb52 }

type MsgsStateInfo struct {
	ReqMsgID int64
	Info     []byte
}

func (*MsgsStateInfo) CRC() uint32 { return 0x04deb57d }

type MsgsAllInfo struct {
	MsgIDs []int64
	Info   []byte
}

func (*MsgsAllInfo) CRC() uint32 { return 0x8cc0d131 }

type MsgsDetailedInfo struct {
	MsgID       int64
	AnswerMsgID int64
	Bytes       int32
	Status      int32
}

func (*MsgsDetailedInfo) CRC() uint32 { return 0x276d3ec6 }

type MsgsNewDetailedInfo struct {
	AnswerMsgID int64
	Bytes       int32
	Status      int32
}

func (*MsgsNewDetailedInfo) CRC() uint32 { return 0x809db6df }
