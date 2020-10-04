// типы, которые описывает mtproto, некоторые декодируются очень специфическим способом, поэтому размещены здесь

package serialize

import (
	"github.com/xelaj/errs"
	"bytes"
	"compress/gzip"
	"fmt"
	"reflect"

	"github.com/k0kubun/pp"

	"github.com/xelaj/go-dry"
)

// TYPES

type ResPQ struct {
	Nonce        *Int128
	ServerNonce  *Int128
	Pq           []byte
	Fingerprints []int64
}

func (_ *ResPQ) CRC() uint32 {
	return 0x05162463
}

func (t *ResPQ) Encode() []byte {
	buf := NewEncoder()
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutMessage(t.Pq)
	buf.PutVector(t.Fingerprints)
	return buf.GetBuffer()
}

func (e *ResPQ) DecodeFrom(d *Decoder) {
	e.Nonce = d.PopInt128()
	e.ServerNonce = d.PopInt128()
	e.Pq = d.PopMessage()
	e.Fingerprints = d.PopVector(int64Type).([]int64)
}

type PQInnerData struct {
	Pq          []byte
	P           []byte
	Q           []byte
	Nonce       *Int128
	ServerNonce *Int128
	NewNonce    *Int256
}

func (_ *PQInnerData) CRC() uint32 {
	return 0x83c95aec
}

func (t *PQInnerData) Encode() []byte {
	buf := NewEncoder()
	buf.PutUint(t.CRC())
	buf.PutMessage(t.Pq)
	buf.PutMessage(t.P)
	buf.PutMessage(t.Q)
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutInt256(t.NewNonce)
	return buf.GetBuffer()
}

func (e *PQInnerData) DecodeFrom(d *Decoder) {
	e.Pq = d.PopMessage()
	e.P = d.PopMessage()
	e.Q = d.PopMessage()
	e.Nonce = d.PopInt128()
	e.ServerNonce = d.PopInt128()
	e.NewNonce = d.PopInt256()
}

type ServerDHParamsFail struct {
	Nonce        *Int128
	ServerNonce  *Int128
	NewNonceHash *Int128
}

func (t *ServerDHParamsFail) ImplementsServerDHParams() {}

func (_ *ServerDHParamsFail) CRC() uint32 {
	return 0x79cb045d
}

func (t *ServerDHParamsFail) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutInt128(t.NewNonceHash)
	return buf.Result()
}

func (t *ServerDHParamsFail) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.NewNonceHash = d.PopInt128()
}

type ServerDHParamsOk struct {
	Nonce           *Int128
	ServerNonce     *Int128
	EncryptedAnswer []byte
}

func (t *ServerDHParamsOk) ImplementsServerDHParams() {}

func (_ *ServerDHParamsOk) CRC() uint32 {
	return 0xd0e8075c
}

func (t *ServerDHParamsOk) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutMessage(t.EncryptedAnswer)
	return buf.Result()
}

func (t *ServerDHParamsOk) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.EncryptedAnswer = d.PopMessage()
}

type ServerDHInnerData struct {
	Nonce       *Int128
	ServerNonce *Int128
	G           int32
	DhPrime     []byte
	GA          []byte
	ServerTime  int32
}

func (_ *ServerDHInnerData) CRC() uint32 {
	return 0xb5890dba
}

func (t *ServerDHInnerData) Encode() []byte {
	panic("not implemented")
}

func (t *ServerDHInnerData) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.G = d.PopInt()
	t.DhPrime = d.PopMessage()
	t.GA = d.PopMessage()
	t.ServerTime = d.PopInt()
}

type ClientDHInnerData struct {
	Nonce       *Int128
	ServerNonce *Int128
	Retry       int64
	GB          []byte
}

func (_ *ClientDHInnerData) CRC() uint32 {
	return 0x6643b654
}

func (t *ClientDHInnerData) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutLong(t.Retry)
	buf.PutMessage(t.GB)
	return buf.Result()
}

func (t *ClientDHInnerData) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.Retry = d.PopLong()
	t.GB = d.PopMessage()
}

type DHGenOk struct {
	Nonce         *Int128
	ServerNonce   *Int128
	NewNonceHash1 *Int128
}

func (t *DHGenOk) ImplementsSetClientDHParamsAnswer() {}

func (_ *DHGenOk) CRC() uint32 {
	return 0x3bcbf734
}

func (t *DHGenOk) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutInt128(t.NewNonceHash1)
	return buf.Result()
}

func (t *DHGenOk) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.NewNonceHash1 = d.PopInt128()
}

type DHGenRetry struct {
	Nonce         *Int128
	ServerNonce   *Int128
	NewNonceHash2 *Int128
}

func (t *DHGenRetry) ImplementsSetClientDHParamsAnswer() {}

func (_ *DHGenRetry) CRC() uint32 {
	return 0x46dc1fb9
}

func (t *DHGenRetry) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutInt128(t.Nonce)
	buf.PutInt128(t.ServerNonce)
	buf.PutInt128(t.NewNonceHash2)
	return buf.Result()
}

func (t *DHGenRetry) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.NewNonceHash2 = d.PopInt128()
}

type DHGenFail struct {
	Nonce         *Int128
	ServerNonce   *Int128
	NewNonceHash3 *Int128
}

func (t *DHGenFail) ImplementsSetClientDHParamsAnswer() {}

func (_ *DHGenFail) CRC() uint32 {
	return 0xa69dae02
}

func (t *DHGenFail) Encode() []byte {
	panic("not implemented")
}

func (t *DHGenFail) DecodeFrom(d *Decoder) {
	t.Nonce = d.PopInt128()
	t.ServerNonce = d.PopInt128()
	t.NewNonceHash3 = d.PopInt128()
}

type RpcResult struct {
	ReqMsgID int64
	Obj      TL
}

func (_ *RpcResult) CRC() uint32 {
	return 0xf35c6d01
}

func (t *RpcResult) Encode() []byte {
	panic("not implemented")
}

func (t *RpcResult) DecodeFrom(d *Decoder) {
	t.ReqMsgID = d.PopLong()
	t.Obj = d.PopObj()
}

type RpcError struct {
	ErrorCode    int32
	ErrorMessage string
}

func (_ *RpcError) CRC() uint32 {
	return 0x2144ca19
}

func (t *RpcError) Encode() []byte {
	panic("makes no sense")
}

func (t *RpcError) DecodeFrom(d *Decoder) {
	t.ErrorCode = d.PopInt()
	t.ErrorMessage = d.PopString()
}

type RpcAnswerUnknown struct{}

func (t *RpcAnswerUnknown) ImplementsRpcDropAnswer() {}

func (_ *RpcAnswerUnknown) CRC() uint32 {
	return 0x5e2ad36e
}

func (t *RpcAnswerUnknown) Encode() []byte {
	panic("makes no sense")
}

func (t *RpcAnswerUnknown) DecodeFrom(d *Decoder) {
}

type RpcAnswerDroppedRunning struct{}

func (t *RpcAnswerDroppedRunning) ImplementsRpcDropAnswer() {}

func (_ *RpcAnswerDroppedRunning) CRC() uint32 {
	return 0xcd78e586
}

func (t *RpcAnswerDroppedRunning) Encode() []byte {
	panic("makes no sense")
}

func (t *RpcAnswerDroppedRunning) DecodeFrom(d *Decoder) {
}

type RpcAnswerDropped struct {
	MsgID int64
	SewNo int32
	Bytes int32
}

func (t *RpcAnswerDropped) ImplementsRpcDropAnswer() {}

func (_ *RpcAnswerDropped) CRC() uint32 {
	return 0xa43ad8b7
}

func (t *RpcAnswerDropped) Encode() []byte {
	panic("makes no sense")
}

func (t *RpcAnswerDropped) DecodeFrom(d *Decoder) {
	t.MsgID = d.PopLong()
	t.SewNo = d.PopInt()
	t.Bytes = d.PopInt()
}

type FutureSalt struct {
	ValidSince int32
	ValidUntil int32
	Salt       int64
}

func (_ *FutureSalt) CRC() uint32 {
	return 0x0949d9dc
}

func (t *FutureSalt) Encode() []byte {
	panic("makes no sense")
}

func (t *FutureSalt) DecodeFrom(d *Decoder) {
	t.ValidSince = d.PopInt()
	t.ValidUntil = d.PopInt()
	t.Salt = d.PopLong()
}

type FutureSalts struct {
	ReqMsgID int64
	Now      int32
	Salts    []*FutureSalt
}

func (_ *FutureSalts) CRC() uint32 {
	return 0xae500895
}

func (t *FutureSalts) Encode() []byte {
	panic("makes no sense")
}

func (t *FutureSalts) DecodeFrom(d *Decoder) {
	t.ReqMsgID = d.PopLong()
	t.Now = d.PopInt()
	t.Salts = d.PopVector(reflect.TypeOf(&FutureSalt{})).([]*FutureSalt)
}

type Pong struct {
	MsgID  int64
	PingID int64
}

func (_ *Pong) CRC() uint32 {
	return 0x347773c5
}

func (t *Pong) Encode() []byte {
	panic("not implemented")
}

func (t *Pong) DecodeFrom(d *Decoder) {
	t.MsgID = d.PopLong()
	t.PingID = d.PopLong()
}

// destroy_session_ok#e22045fc session_id:long = DestroySessionRes;
// destroy_session_none#62d350c9 session_id:long = DestroySessionRes;

type NewSessionCreated struct {
	FirstMsgID int64
	UniqueID   int64
	ServerSalt int64
}

func (_ *NewSessionCreated) CRC() uint32 {
	return 0x9ec20908
}

func (t *NewSessionCreated) Encode() []byte {
	panic("not implemented")
}

func (t *NewSessionCreated) DecodeFrom(d *Decoder) {
	t.FirstMsgID = d.PopLong()
	t.UniqueID = d.PopLong()
	t.ServerSalt = d.PopLong()
}

//! исключение из правил: это оказывается почти-вектор, т.к.
//  записан как `msg_container#73f1f8dc messages:vector<%Message> = MessageContainer;`
//  судя по всему, <%Type> означает, что может это неявный вектор???
//! возможно разработчики в этот момент поехаи кукухой, я не знаю правда
type MessageContainer []*EncryptedMessage

func (_ *MessageContainer) CRC() uint32 {
	return 0x73f1f8dc
}

func (t *MessageContainer) Encode() []byte {
	buf := NewEncoder()
	buf.PutUint(t.CRC())

	buf.PutInt(int32(len(*t)))
	for _, msg := range *t {
		encoded := msg.Msg.Encode()

		buf.PutLong(msg.MsgID)
		buf.PutInt(msg.SeqNo)
		//         msgID     seqNo     len             object
		buf.PutInt(LongLen + WordLen + WordLen + int32(len(encoded)))
		buf.PutRawBytes(encoded)
	}
	return buf.GetBuffer()
}

func (t *MessageContainer) DecodeFrom(d *Decoder) {
	count := int(d.PopInt())
	arr := make([]*EncryptedMessage, count)
	for i := 0; i < count; i++ {
		msg := new(EncryptedMessage)
		msg.MsgID = d.PopLong()
		msg.SeqNo = d.PopInt()
		_ = d.PopInt() // size, но нам нахуй не нужен
		msg.Msg = d.PopObj().(TL)
		arr[i] = msg
	}
	*t = arr
}

type Message struct {
	MsgID int64
	SeqNo int32
	Bytes int32
	Body  TL
}

type MsgCopy struct {
	OrigMessage *Message
}

func (_ *MsgCopy) CRC() uint32 {
	return 0xe06046b2
}

func (t *MsgCopy) Encode() []byte {
	panic("makes no sense")
}

func (t *MsgCopy) DecodeFrom(d *Decoder) {
	pp.Println(d)
	panic("очень специфичный конструктор Message, надо сначала посмотреть, как это что это")
}

type GzipPacked struct {
	Obj TL
}

func (_ *GzipPacked) CRC() uint32 {
	return 0x3072cfa1
}

func (t *GzipPacked) Encode() []byte {
	panic("not implemented")
}

func (t *GzipPacked) DecodeFrom(d *Decoder) {
	// TODO: СТАНДАРТНЫЙ СУКА ПАКЕТ gzip пишет "gzip: invalid header". при этом как я разобрался, в сам гзип попадает кусок, который находится за миллиард бит от реального сообщения
	//       например: сообщение начинается с 0x1f 0x8b 0x08 0x00 ..., но при этом в сам гзип отдается кусок, который дальше начала сообщения за 500+ байт
	//! вот ЭТОТ кусок работает. так что наверное не будем трогать, дай бог чтоб работал

	obj := make([]byte, 0, 4096)

	var buf bytes.Buffer
	_, _ = buf.Write(d.PopMessage())
	gz, err := gzip.NewReader(&buf)
	dry.PanicIfErr(err)
	b := make([]byte, 4096)
	for {
		n, _ := gz.Read(b)

		obj = append(obj, b[0:n]...)
		if n <= 0 {
			break
		}
	}

	decoder := NewDecoder(obj)
	t.Obj = decoder.PopObj()

	//? это то что я пытался сделать
	// data := d.PopMessage()
	// gz, err := gzip.NewReader(bytes.NewBuffer(data))
	// dry.PanicIfErr(err)

	// decompressed, err := ioutil.ReadAll(gz)
	// dry.PanicIfErr(err)

	// decoder := NewDecoder(decompressed)
	// t.Obj = decoder.PopObj()

	// return
}

type MsgsAck struct {
	MsgIds []int64
}

func (_ *MsgsAck) CRC() uint32 {
	return 0x62d6b459
}

func (t *MsgsAck) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutVector(t.MsgIds)
	return buf.Result()
}

func (t *MsgsAck) DecodeFrom(d *Decoder) {
	t.MsgIds = d.PopVector(int64Type).([]int64)
}

type BadMsgNotification struct {
	BadMsgID    int64
	MadMsgSeqNo int32
	ErrorCode   int32
}

func (t *BadMsgNotification) ImplementsBadMsgNotification() {}

func (_ *BadMsgNotification) CRC() uint32 {
	return 0xa7eff811
}

func (t *BadMsgNotification) Encode() []byte {
	buf := NewEncoder()
	buf.PutCRC(t.CRC())
	buf.PutLong(t.BadMsgID)
	buf.PutInt(t.MadMsgSeqNo)
	buf.PutInt(t.ErrorCode)
	return buf.Result()
}

func (t *BadMsgNotification) DecodeFrom(d *Decoder) {
	t.BadMsgID = d.PopLong()
	t.MadMsgSeqNo = d.PopInt()
	t.ErrorCode = d.PopInt()
}

type BadServerSalt struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	ErrorCode   int32
	NewSalt     int64
}

func (t *BadServerSalt) ImplementsBadMsgNotification() {}

func (_ *BadServerSalt) CRC() uint32 {
	return 0xedab447b
}

func (t *BadServerSalt) Encode() []byte {
	panic("makes no sense")
}

func (t *BadServerSalt) DecodeFrom(d *Decoder) {
	t.BadMsgID = d.PopLong()
	t.BadMsgSeqNo = d.PopInt()
	t.ErrorCode = d.PopInt()
	t.NewSalt = d.PopLong()
}

// msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;

type MsgResendReq struct {
	MsgIds []int64
}

func (_ *MsgResendReq) CRC() uint32 {
	return 0x7d861a08

}

func (t *MsgResendReq) Encode() []byte {
	panic("not implemented")
}

func (t *MsgResendReq) DecodeFrom(d *Decoder) {
	panic("not implemented")
}

type MsgsStateReq struct {
	MsgIds []int64
}

func (_ *MsgsStateReq) CRC() uint32 {
	return 0xda69fb52

}

func (t *MsgsStateReq) Encode() []byte {
	panic("not implemented")
}

func (t *MsgsStateReq) DecodeFrom(d *Decoder) {
	panic("not implemented")
}

type MsgsStateInfo struct {
	ReqMsgId int64
	Info     []byte
}

func (_ *MsgsStateInfo) CRC() uint32 {
	return 0x04deb57d

}

func (t *MsgsStateInfo) Encode() []byte {
	panic("not implemented")
}

func (t *MsgsStateInfo) DecodeFrom(d *Decoder) {
	panic("not implemented")
}

type MsgsAllInfo struct {
	MsgIds []int64
	Info   []byte
}

func (_ *MsgsAllInfo) CRC() uint32 {
	return 0x8cc0d131

}

func (t *MsgsAllInfo) Encode() []byte {
	panic("not implemented")
}

func (t *MsgsAllInfo) DecodeFrom(d *Decoder) {
	panic("not implemented")
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

func (t *MsgsDetailedInfo) Encode() []byte {
	panic("not implemented")
}

func (t *MsgsDetailedInfo) DecodeFrom(d *Decoder) {
	panic("not implemented")
}

type MsgsNewDetailedInfo struct {
	AnswerMsgId int64
	Bytes       int32
	Status      int32
}

func (_ *MsgsNewDetailedInfo) CRC() uint32 {
	return 0x809db6df

}

func (t *MsgsNewDetailedInfo) Encode() []byte {
	panic("not implemented")
}

func (t *MsgsNewDetailedInfo) DecodeFrom(d *Decoder) {
	panic("not implemented")
}

type ServerDHParams interface {
	TL
	ImplementsServerDHParams()
}

type SetClientDHParamsAnswer interface {
	TL
	ImplementsSetClientDHParamsAnswer()
}

func GenerateCommonObject(constructorID uint32) (obj TL, isEnum bool, err error) {
	switch constructorID {
	case 0x05162463:
		return &ResPQ{}, false, nil
	case 0x83c95aec:
		return &PQInnerData{}, false, nil
	case 0x79cb045d:
		return &ServerDHParamsFail{}, false, nil
	case 0xd0e8075c:
		return &ServerDHParamsOk{}, false, nil
	case 0xb5890dba:
		return &ServerDHInnerData{}, false, nil
	case 0x6643b654:
		return &ClientDHInnerData{}, false, nil
	case 0x3bcbf734:
		return &DHGenOk{}, false, nil
	case 0x46dc1fb9:
		return &DHGenRetry{}, false, nil
	case 0xa69dae02:
		return &DHGenFail{}, false, nil
	case 0xf35c6d01:
		return &RpcResult{}, false, nil
	case 0x2144ca19:
		return &RpcError{}, false, nil
	case 0x5e2ad36e:
		return &RpcAnswerUnknown{}, false, nil
	case 0xcd78e586:
		return &RpcAnswerDroppedRunning{}, false, nil
	case 0xa43ad8b7:
		return &RpcAnswerDropped{}, false, nil
	case 0x0949d9dc:
		return &FutureSalt{}, false, nil
	case 0xae500895:
		return &FutureSalts{}, false, nil
	case 0x347773c5:
		return &Pong{}, false, nil
	//case 0xe22045fc:
	//	return &destroy_session_ok{}, false, nil
	//case 0x62d350c9:
	//	return &destroy_session_none{}, false, nil
	case 0x9ec20908:
		return &NewSessionCreated{}, false, nil
	case 0x73f1f8dc: //! SPECIFIC
		return &MessageContainer{}, false, nil
	case 0xe06046b2:
		return &MsgCopy{}, false, nil
	case 0x3072cfa1:
		return &GzipPacked{}, false, nil
	case 0x62d6b459:
		return &MsgsAck{}, false, nil
	case 0xa7eff811:
		return &BadMsgNotification{}, false, nil
	case 0xedab447b:
		return &BadServerSalt{}, false, nil
	case 0x7d861a08:
		return &MsgResendReq{}, false, nil
	case 0xda69fb52:
		return &MsgsStateReq{}, false, nil
	case 0x04deb57d:
		return &MsgsStateInfo{}, false, nil
	case 0x8cc0d131:
		return &MsgsAllInfo{}, false, nil
	case 0x276d3ec6:
		return &MsgsDetailedInfo{}, false, nil
	case 0x809db6df:
		return &MsgsNewDetailedInfo{}, false, nil
	default:
		return nil, false, errs.NotFound("constructorID", fmt.Sprintf("%#v", constructorID))
	}
}
