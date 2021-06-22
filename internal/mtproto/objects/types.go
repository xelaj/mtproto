// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package objects

// some types are decoding VEEEEEEERY specific way, so it stored here and only here.

import (
	"bytes"
	"compress/gzip"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/internal/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/messages"
)

// TYPES

// Null это пустой объект, который нужен для передачи в каналы TL с информацией, что ответа можно не ждать
type Null struct {
}

func (*Null) CRC() uint32 {
	panic("makes no sense")
}

type ResPQ struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	Pq           []byte
	Fingerprints []int64
}

func (*ResPQ) CRC() uint32 {
	return 0x05162463 //nolint:gomnd not magic
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
	return 0x83c95aec //nolint:gomnd not magic
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
	return 0x79cb045d //nolint:gomnd not magic
}

type ServerDHParamsOk struct {
	Nonce           *tl.Int128
	ServerNonce     *tl.Int128
	EncryptedAnswer []byte
}

func (*ServerDHParamsOk) ImplementsServerDHParams() {}

func (*ServerDHParamsOk) CRC() uint32 {
	return 0xd0e8075c //nolint:gomnd not magic
}

type ServerDHInnerData struct { //nolint:maligned telegram require to fix
	Nonce       *tl.Int128
	ServerNonce *tl.Int128
	G           int32
	DhPrime     []byte
	GA          []byte
	ServerTime  int32
}

func (*ServerDHInnerData) CRC() uint32 {
	return 0xb5890dba //nolint:gomnd not magic
}

type ClientDHInnerData struct {
	Nonce       *tl.Int128
	ServerNonce *tl.Int128
	Retry       int64
	GB          []byte
}

func (*ClientDHInnerData) CRC() uint32 {
	return 0x6643b654 //nolint:gomnd not magic
}

type DHGenOk struct {
	Nonce         *tl.Int128
	ServerNonce   *tl.Int128
	NewNonceHash1 *tl.Int128
}

func (t *DHGenOk) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenOk) CRC() uint32 {
	return 0x3bcbf734 //nolint:gomnd not magic
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
	return 0x46dc1fb9 //nolint:gomnd not magic
}

type DHGenFail struct {
	Nonce         *tl.Int128
	ServerNonce   *tl.Int128
	NewNonceHash3 *tl.Int128
}

func (*DHGenFail) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenFail) CRC() uint32 {
	return 0xa69dae02 //nolint:gomnd not magic
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
	return 0x2144ca19 //nolint:gomnd not magic
}

type RpcDropAnswer interface {
	tl.Object
	ImplementsRpcDropAnswer()
}

type RpcAnswerUnknown null

func (*RpcAnswerUnknown) ImplementsRpcDropAnswer() {}

func (*RpcAnswerUnknown) CRC() uint32 {
	return 0x5e2ad36e //nolint:gomnd not magic
}

type RpcAnswerDroppedRunning null

func (*RpcAnswerDroppedRunning) ImplementsRpcDropAnswer() {}

func (*RpcAnswerDroppedRunning) CRC() uint32 {
	return 0xcd78e586 //nolint:gomnd not magic
}

type RpcAnswerDropped struct {
	MsgID int64
	SewNo int32
	Bytes int32
}

func (*RpcAnswerDropped) ImplementsRpcDropAnswer() {}

func (*RpcAnswerDropped) CRC() uint32 {
	return 0xa43ad8b7 //nolint:gomnd not magic
}

type FutureSalt struct {
	ValidSince int32
	ValidUntil int32
	Salt       int64
}

func (*FutureSalt) CRC() uint32 {
	return 0x0949d9dc //nolint:gomnd not magic
}

type FutureSalts struct {
	ReqMsgID int64
	Now      int32
	Salts    []*FutureSalt
}

func (*FutureSalts) CRC() uint32 {
	return 0xae500895 //nolint:gomnd not magic
}

type Pong struct {
	MsgID  int64
	PingID int64
}

func (*Pong) CRC() uint32 {
	return 0x347773c5 //nolint:gomnd not magic
}

// destroy_session_ok#e22045fc session_id:long = DestroySessionRes;
// destroy_session_none#62d350c9 session_id:long = DestroySessionRes;

type NewSessionCreated struct {
	FirstMsgID int64
	UniqueID   int64
	ServerSalt int64
}

func (*NewSessionCreated) CRC() uint32 {
	return 0x9ec20908 //nolint:gomnd not magic
}

//! исключение из правил: это оказывается почти-вектор, т.к.
//  записан как `msg_container#73f1f8dc messages:vector<%Message> = MessageContainer;`
//  судя по всему, <%Type> означает, что может это неявный вектор???
//! возможно разработчики в этот момент поехаи кукухой, я не знаю правда
type MessageContainer []*messages.Encrypted

func (*MessageContainer) CRC() uint32 {
	return 0x73f1f8dc //nolint:gomnd not magic
}

func (t *MessageContainer) MarshalTL(e *tl.Encoder) error {
	e.PutUint(t.CRC())
	e.PutInt(int32(len(*t)))
	if err := e.CheckErr(); err != nil {
		return err
	}

	for _, msg := range *t {
		e.PutLong(msg.MsgID)
		e.PutInt(msg.SeqNo)
		//       msgID        seqNo        len                object
		e.PutInt(tl.LongLen + tl.WordLen + tl.WordLen + int32(len(msg.Msg)))
		e.PutRawBytes(msg.Msg)
	}
	return e.CheckErr()
}

func (t *MessageContainer) UnmarshalTL(d *tl.Decoder) error {
	count := int(d.PopInt())
	arr := make([]*messages.Encrypted, count)
	for i := 0; i < count; i++ {
		msg := new(messages.Encrypted)
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
	return 0xe06046b2 //nolint:gomnd not magic
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
	obj, err := t.popMessageAsBytes(d)
	if err != nil {
		return err
	}

	t.Obj, err = tl.DecodeUnknownObject(obj)
	if err != nil {
		return errors.Wrap(err, "parsing gzipped object")
	}

	return nil
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
	if err != nil {
		return nil, errors.Wrap(err, "creating gzip reader")
	}

	b := make([]byte, 4096)
	for {
		n, _ := gz.Read(b)

		decompressed = append(decompressed, b[0:n]...)
		if n <= 0 {
			break
		}
	}

	return decompressed, nil
	//? это то что я пытался сделать
	// data := d.PopMessage()
	// gz, err := gzip.NewReader(bytes.NewBuffer(data))
	// check(err)

	// decompressed, err := ioutil.ReadAll(gz)
	// check(err)

	// return decompressed
}

type MsgsAck struct {
	MsgIDs []int64
}

func (*MsgsAck) CRC() uint32 {
	return 0x62d6b459 //nolint:gomnd not magic
}

type BadMsgNotification struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	Code        int32
}

func (*BadMsgNotification) ImplementsBadMsgNotification() {}

func (*BadMsgNotification) CRC() uint32 {
	return 0xa7eff811 //nolint:gomnd not magic
}

type BadServerSalt struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	ErrorCode   int32
	NewSalt     int64
}

func (*BadServerSalt) ImplementsBadMsgNotification() {}

func (*BadServerSalt) CRC() uint32 {
	return 0xedab447b //nolint:gomnd not magic
}

// msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;

type MsgResendReq struct {
	MsgIDs []int64
}

func (*MsgResendReq) CRC() uint32 {
	return 0x7d861a08 //nolint:gomnd not magic
}

type MsgsStateReq struct {
	MsgIDs []int64
}

func (*MsgsStateReq) CRC() uint32 {
	return 0xda69fb52 //nolint:gomnd not magic
}

type MsgsStateInfo struct {
	ReqMsgID int64
	Info     []byte
}

func (*MsgsStateInfo) CRC() uint32 {
	return 0x04deb57d //nolint:gomnd not magic
}

type MsgsAllInfo struct {
	MsgIDs []int64
	Info   []byte
}

func (*MsgsAllInfo) CRC() uint32 {
	return 0x8cc0d131 //nolint:gomnd not magic
}

type MsgsDetailedInfo struct {
	MsgID       int64
	AnswerMsgID int64
	Bytes       int32
	Status      int32
}

func (*MsgsDetailedInfo) CRC() uint32 {
	return 0x276d3ec6 //nolint:gomnd not magic
}

type MsgsNewDetailedInfo struct {
	AnswerMsgID int64
	Bytes       int32
	Status      int32
}

func (*MsgsNewDetailedInfo) CRC() uint32 {
	return 0x809db6df //nolint:gomnd not magic
}
