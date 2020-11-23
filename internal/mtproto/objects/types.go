// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package objects

// some types are decoding VEEEEEEERY specific way, so it stored here and only here.

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"reflect"

	"github.com/k0kubun/pp"
	"github.com/xelaj/mtproto/encoding/tl"
)

// TYPES

type ResPQ struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	Pq           []byte
	Fingerprints []int64
}

func (*ResPQ) CRC() uint32 {
	return 0x05162463
}

type PQInnerData struct {
	Pq          []byte
	P           []byte
	Q           []byte
	Nonce       *tl.Int128
	ServerNonce *tl.Int128
	NewNonce    *tl.Int256
}

func (*PQInnerData) CRC() uint32 {
	return 0x83c95aec
}

type ServerDHParams interface {
	tl.Object
	ImplementsServerDHParams()
}


type ServerDHParamsFail struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	NewNonceHash *tl.Int128
}

func (*ServerDHParamsFail) ImplementsServerDHParams() {}

func (*ServerDHParamsFail) CRC() uint32 {
	return 0x79cb045d
}

type ServerDHParamsOk struct {
	Nonce           *tl.Int128
	ServerNonce     *tl.Int128
	EncryptedAnswer []byte
}

func (*ServerDHParamsOk) ImplementsServerDHParams() {}

func (*ServerDHParamsOk) CRC() uint32 {
	return 0xd0e8075c
}

type ServerDHInnerData struct {
	Nonce       *tl.Int128
	ServerNonce *tl.Int128
	G           int32
	DhPrime     []byte
	GA          []byte
	ServerTime  int32
}

func (*ServerDHInnerData) CRC() uint32 {
	return 0xb5890dba
}

type ClientDHInnerData struct {
	Nonce       *tl.Int128
	ServerNonce *tl.Int128
	Retry       int64
	GB          []byte
}

func (*ClientDHInnerData) CRC() uint32 {
	return 0x6643b654
}

type DHGenOk struct {
	Nonce         *tl.Int128
	ServerNonce   *tl.Int128
	NewNonceHash1 *tl.Int128
}

func (t *DHGenOk) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenOk) CRC() uint32 {
	return 0x3bcbf734
}

type SetClientDHParamsAnswer interface {
	tl.Object
	ImplementsSetClientDHParamsAnswer()
}

type DHGenRetry struct {
	Nonce         *tl.Int128
	ServerNonce   *tl.Int128
	NewNonceHash2 *tl.Int128
}

func (*DHGenRetry) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenRetry) CRC() uint32 {
	return 0x46dc1fb9
}

type DHGenFail struct {
	Nonce         *tl.Int128
	ServerNonce   *tl.Int128
	NewNonceHash3 *tl.Int128
}

func (*DHGenFail) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenFail) CRC() uint32 {
	return 0xa69dae02
}

type RpcResult struct {
	ReqMsgID int64
	Obj      tl.Object
}

func (*RpcResult) CRC() uint32 {
	return CrcRpcResult
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

func (*RpcError) CRC() uint32 {
	return 0x2144ca19
}

type RpcDropAnswer interface{
	tl.Object
	ImplementsRpcDropAnswer()
}

type RpcAnswerUnknown struct{}

func (*RpcAnswerUnknown) ImplementsRpcDropAnswer() {}

func (*RpcAnswerUnknown) CRC() uint32 {
	return 0x5e2ad36e
}

type RpcAnswerDroppedRunning struct{}

func (*RpcAnswerDroppedRunning) ImplementsRpcDropAnswer() {}

func (*RpcAnswerDroppedRunning) CRC() uint32 {
	return 0xcd78e586
}

type RpcAnswerDropped struct {
	MsgID int64
	SewNo int32
	Bytes int32
}

func (*RpcAnswerDropped) ImplementsRpcDropAnswer() {}

func (*RpcAnswerDropped) CRC() uint32 {
	return 0xa43ad8b7
}

type FutureSalt struct {
	ValidSince int32
	ValidUntil int32
	Salt       int64
}

func (*FutureSalt) CRC() uint32 {
	return 0x0949d9dc
}

type FutureSalts struct {
	ReqMsgID int64
	Now      int32
	Salts    []*FutureSalt
}

func (*FutureSalts) CRC() uint32 {
	return 0xae500895
}

type Pong struct {
	MsgID  int64
	PingID int64
}

func (*Pong) CRC() uint32 {
	return 0x347773c5
}

// destroy_session_ok#e22045fc session_id:long = DestroySessionRes;
// destroy_session_none#62d350c9 session_id:long = DestroySessionRes;

type NewSessionCreated struct {
	FirstMsgID int64
	UniqueID   int64
	ServerSalt int64
}

func (*NewSessionCreated) CRC() uint32 {
	return 0x9ec20908
}

//! исключение из правил: это оказывается почти-вектор, т.к.
//  записан как `msg_container#73f1f8dc messages:vector<%Message> = MessageContainer;`
//  судя по всему, <%Type> означает, что может это неявный вектор???
//! возможно разработчики в этот момент поехаи кукухой, я не знаю правда
type MessageContainer []*EncryptedMessage

func (_ *MessageContainer) CRC() uint32 {
	return 0x73f1f8dc
}

func (t *MessageContainer) MarshalTL(e *tl.Encoder) error {
	e.PutUint(t.CRC())
	e.PutInt(int32(len(*t)))
	if e.CheckErr() != nil {
		return err
	}

	for _, msg := range *t {
		e.PutLong(msg.MsgID)
		e.PutInt(msg.SeqNo)
		//       msgID     seqNo     len             object
		e.PutInt(LongLen + WordLen + WordLen + int32(len(msg.Msg)))
		e.PutRawBytes(msg.Msg)
	}
	return buf.GetBuffer()
}

func (t *MessageContainer) UnmarshalTL(d *tl.Decoder) error {
	crc := d.PopCRC()
	if crc != t.CRC() {
		return errors.New("wrong CRC code, want %#v, got %#v", t.CRC(), crc)
	}

	count := int(d.PopInt())
	arr := make([]*EncryptedMessage, count)
	for i := 0; i < count; i++ {
		msg := new(EncryptedMessage)
		msg.MsgID = d.PopLong()
		msg.SeqNo = d.PopInt()
		size := d.PopInt()
		msg.Msg = d.PopRawBytes(int(size))
		arr[i] = msg
	}
	*t = arr

	return nil
}

type Message struct {
	MsgID int64
	SeqNo int32
	Bytes int32
	Body  tl.Object
}

type MsgCopy struct {
	OrigMessage *Message
}

func (*MsgCopy) CRC() uint32 {
	return 0xe06046b2
}

type GzipPacked struct {
	Obj tl.Object
}

func (*GzipPacked) CRC() uint32 {
	return CrcGzipPacked
}

func (*GzipPacked) MarshalTL(e *tl.Encoder) error {
	panic("not implemented")
}

func (t *GzipPacked) UnmarshalTL(d *tl.Decoder) error {
	crc := d.PopCRC()
	if crc != t.CRC() {
		return errors.New("wrong CRC code, want %#v, got %#v", t.CRC(), crc)
	}

	obj := t.popMessageAsBytes(d)
	innerDecoder := NewDecoder(obj)
	t.Obj = innerDecoder.PopObj()
}

func (*GzipPacked) popMessageAsBytes(d *tl.Decoder) ([]byte, error) {
	// TODO: СТАНДАРТНЫЙ СУКА ПАКЕТ gzip пишет "gzip: invalid header". при этом как я разобрался, в
	//       сам гзип попадает кусок, который находится за миллиард бит от реального сообщения
	//       например: сообщение начинается с 0x1f 0x8b 0x08 0x00 ..., но при этом в сам гзип
	//       отдается кусок, который дальше начала сообщения за 500+ байт
	//! вот ЭТОТ кусок работает. так что наверное не будем трогать, дай бог чтоб работал

	decompressed := make([]byte, 0, 4096)

	var buf bytes.Buffer
	_, _ = buf.Write(d.PopMessage())
	gz, err := gzip.NewReader(&buf)

	b := make([]byte, 4096)
	for {
		n, _ := gz.Read(b)

		decompressed = append(decompressed, b[0:n]...)
		if n <= 0 {
			break
		}
	}

	return decompressed
	//? это то что я пытался сделать
	// data := d.PopMessage()
	// gz, err := gzip.NewReader(bytes.NewBuffer(data))
	// dry.PanicIfErr(err)

	// decompressed, err := ioutil.ReadAll(gz)
	// dry.PanicIfErr(err)

	// return decompressed
}

type MsgsAck struct {
	MsgIds []int64
}

func (*MsgsAck) CRC() uint32 {
	return 0x62d6b459
}

type BadMsgNotification struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	Code        int32
}

func (*BadMsgNotification) ImplementsBadMsgNotification() {}

func (*BadMsgNotification) CRC() uint32 {
	return 0xa7eff811
}

type BadServerSalt struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	ErrorCode   int32
	NewSalt     int64
}

func (*BadServerSalt) ImplementsBadMsgNotification() {}

func (*BadServerSalt) CRC() uint32 {
	return 0xedab447b
}

// msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;

type MsgResendReq struct {
	MsgIds []int64
}

func (*MsgResendReq) CRC() uint32 {
	return 0x7d861a08

}

type MsgsStateReq struct {
	MsgIds []int64
}

func (*MsgsStateReq) CRC() uint32 {
	return 0xda69fb52

}


type MsgsStateInfo struct {
	ReqMsgId int64
	Info     []byte
}

func (*MsgsStateInfo) CRC() uint32 {
	return 0x04deb57d

}

type MsgsAllInfo struct {
	MsgIds []int64
	Info   []byte
}

func (*MsgsAllInfo) CRC() uint32 {
	return 0x8cc0d131

}

type MsgsDetailedInfo struct {
	MsgId       int64
	AnswerMsgId int64
	Bytes       int32
	Status      int32
}

func (_ *MsgsDetailedInfo) CRC() uint32 {
	return 0x276d3ec6

}

type MsgsNewDetailedInfo struct {
	AnswerMsgId int64
	Bytes       int32
	Status      int32
}

func (*MsgsNewDetailedInfo) CRC() uint32 {
	return 0x809db6df

}
