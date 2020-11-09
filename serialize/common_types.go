// типы, которые описывает mtproto, некоторые декодируются очень специфическим способом, поэтому размещены здесь

package serialize

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/xelaj/mtproto/encoding/tl"
)

// TYPES

type ResPQ struct {
	Nonce        *Int128
	ServerNonce  *Int128
	Pq           []byte
	Fingerprints []int64
}

func (*ResPQ) CRC() uint32 { return 0x05162463 }

type PQInnerData struct {
	Pq          []byte
	P           []byte
	Q           []byte
	Nonce       *Int128
	ServerNonce *Int128
	NewNonce    *Int256
}

func (*PQInnerData) CRC() uint32 { return 0x83c95aec }

type ServerDHParamsFail struct {
	Nonce        *Int128
	ServerNonce  *Int128
	NewNonceHash *Int128
}

func (t *ServerDHParamsFail) ImplementsServerDHParams() {}

func (_ *ServerDHParamsFail) CRC() uint32 { return 0x79cb045d }

type ServerDHParamsOk struct {
	Nonce           *Int128
	ServerNonce     *Int128
	EncryptedAnswer []byte
}

func (t *ServerDHParamsOk) ImplementsServerDHParams() {}

func (_ *ServerDHParamsOk) CRC() uint32 { return 0xd0e8075c }

type ServerDHInnerData struct {
	Nonce       *Int128
	ServerNonce *Int128
	G           int32
	DhPrime     []byte
	GA          []byte
	ServerTime  int32
}

func (*ServerDHInnerData) CRC() uint32 { return 0xb5890dba }

type ClientDHInnerData struct {
	Nonce       *Int128
	ServerNonce *Int128
	Retry       int64
	GB          []byte
}

func (*ClientDHInnerData) CRC() uint32 { return 0x6643b654 }

type DHGenOk struct {
	Nonce         *Int128
	ServerNonce   *Int128
	NewNonceHash1 *Int128
}

func (t *DHGenOk) ImplementsSetClientDHParamsAnswer() {}

func (_ *DHGenOk) CRC() uint32 { return 0x3bcbf734 }

type DHGenRetry struct {
	Nonce         *Int128
	ServerNonce   *Int128
	NewNonceHash2 *Int128
}

func (*DHGenRetry) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenRetry) CRC() uint32 { return 0x46dc1fb9 }

type DHGenFail struct {
	Nonce         *Int128
	ServerNonce   *Int128
	NewNonceHash3 *Int128
}

func (*DHGenFail) ImplementsSetClientDHParamsAnswer() {}

func (*DHGenFail) CRC() uint32 { return 0xa69dae02 }

type RpcResult struct {
	ReqMsgID int64
	Payload  []byte
}

func (*RpcResult) CRC() uint32 { return 0xf35c6d01 } // CrcRpcResult

func (rpc *RpcResult) UnmarshalTL(r *tl.ReadCursor) (err error) {
	rpc.ReqMsgID, err = r.PopLong()
	if err != nil {
		return err
	}

	rpc.Payload, err = r.GetRestOfMessage()
	return err
}

func (rpc *RpcResult) MarshalTL(w *tl.WriteCursor) error {
	panic("don't use me!")
}

type GzipPacked struct {
	PackedData []byte
}

func (*GzipPacked) CRC() uint32 { return 0x3072cfa1 } // CrcGzipPacked

func (g *GzipPacked) UnmarshalTL(r *tl.ReadCursor) (err error) {
	data, err := r.PopMessage()
	if err != nil {
		panic(err)
	}

	g.PackedData, err = decompressData(data)
	if err != nil {
		panic(err)
	}

	return
}

func (*GzipPacked) MarshalTL(w *tl.WriteCursor) error {
	panic("don't use me")
}

// func (*RpcResult) UnmarshalTL(*tl.ReadCursor) error {
// 	panic("don't use me")
// }

// func (*RpcResult) MarshalTL(*tl.WriteCursor) error {
// 	panic("don't use me")
// }

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
// func (t *RpcResult) DecodeFromButItsVector(r *tl.ReadCursor, as reflect.Type) error {
// 	var err error
// 	t.ReqMsgID, err = r.PopLong()
// 	if err != nil {
// 		return err
// 	}

// 	msg, err := r.GetRestOfMessage()
// 	if err != nil {
// 		return err
// 	}

// 	crc := binary.LittleEndian.Uint32(msg[:tl.WordLen])
// 	if crc == CrcGzipPacked {
// 		if _, err := r.PopCRC(); err != nil {
// 			return err
// 		}

// 		gz := &GzipPacked{}
// 		gz.DecodeFromButItsVector(r, as)
// 		t.Obj = gz.Obj.(*InnerVectorObject)
// 	} else {
// 		vector, err := r.PopVector(as)
// 		if err != nil {
// 			return err
// 		}

// 		t.Obj = &InnerVectorObject{I: vector}
// 	}

// 	return nil
// }

type RpcError struct {
	ErrorCode    int32
	ErrorMessage string
}

func (*RpcError) CRC() uint32 { return 0x2144ca19 }

func (e *RpcError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.ErrorCode, e.ErrorMessage)
}

type RpcAnswerUnknown struct{}

func (*RpcAnswerUnknown) ImplementsRpcDropAnswer() {}

func (*RpcAnswerUnknown) CRC() uint32 { return 0x5e2ad36e }

type RpcAnswerDroppedRunning struct{}

func (*RpcAnswerDroppedRunning) ImplementsRpcDropAnswer() {}

func (*RpcAnswerDroppedRunning) CRC() uint32 { return 0xcd78e586 }

type RpcAnswerDropped struct {
	MsgID int64
	SewNo int32
	Bytes int32
}

func (*RpcAnswerDropped) ImplementsRpcDropAnswer() {}

func (*RpcAnswerDropped) CRC() uint32 { return 0xa43ad8b7 }

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

//! исключение из правил: это оказывается почти-вектор, т.к.
//  записан как `msg_container#73f1f8dc messages:vector<%Message> = MessageContainer;`
//  судя по всему, <%Type> означает, что может это неявный вектор???
//! возможно разработчики в этот момент поехаи кукухой, я не знаю правда
type MessageContainer []*EncryptedMessage

func (*MessageContainer) CRC() uint32 { return 0x73f1f8dc }

func (t *MessageContainer) MarshalTL(w *tl.WriteCursor) error {
	if err := w.PutUint(t.CRC()); err != nil {
		return err
	}

	if err := w.PutUint(uint32(len(*t))); err != nil {
		return err
	}

	for _, msg := range *t {
		if err := w.PutLong(msg.MsgID); err != nil {
			return err
		}

		if err := w.PutUint(uint32(msg.SeqNo)); err != nil {
			return err
		}

		//                            msgID     seqNo     len             object
		if err := w.PutUint(uint32(tl.LongLen + tl.WordLen + tl.WordLen + int32(len(msg.Msg)))); err != nil {
			return err
		}

		if err := w.PutRawBytes(msg.Msg); err != nil {
			return err
		}
	}

	return nil
}

func (t *MessageContainer) UnmarshalTL(r *tl.ReadCursor) error {
	count, err := r.PopUint()
	if err != nil {
		return err
	}

	arr := make([]*EncryptedMessage, count)
	for i := 0; i < int(count); i++ {
		msg := new(EncryptedMessage)
		msg.MsgID, err = r.PopLong()
		if err != nil {
			return err
		}

		seqNo, err := r.PopUint()
		if err != nil {
			return err
		}

		size, err := r.PopUint()
		if err != nil {
			return err
		}

		msg.SeqNo = int32(seqNo)
		msg.Msg, err = r.PopRawBytes(int(size)) // или size * wordLen?
		if err != nil {
			return err
		}

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

func (*MsgCopy) CRC() uint32 { return 0xe06046b2 }

func (t *MsgCopy) UnmarshalTL(r *tl.ReadCursor) error {
	// pp.Println(r)
	panic("очень специфичный конструктор Message, надо сначала посмотреть, как это что это")
}

func (*GzipPacked) popMessageAsBytes(r *tl.ReadCursor) ([]byte, error) {
	// TODO: СТАНДАРТНЫЙ СУКА ПАКЕТ gzip пишет "gzip: invalid header". при этом как я разобрался, в
	//       сам гзип попадает кусок, который находится за миллиард бит от реального сообщения
	//       например: сообщение начинается с 0x1f 0x8b 0x08 0x00 ..., но при этом в сам гзип
	//       отдается кусок, который дальше начала сообщения за 500+ байт
	//! вот ЭТОТ кусок работает. так что наверное не будем трогать, дай бог чтоб работал

	decompressed := make([]byte, 0, 4096)
	var buf bytes.Buffer

	msg, err := r.PopMessage()
	if err != nil {
		return nil, err
	}

	n, err := buf.Write(msg)
	if err != nil {
		return nil, err
	}

	if n != len(msg) {
		return nil, fmt.Errorf("can't write omg")
	}

	gz, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, err
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
	// dry.PanicIfErr(err)

	// decompressed, err := ioutil.ReadAll(gz)
	// dry.PanicIfErr(err)

	// return decompressed
}

type MsgsAck struct {
	MsgIds []int64
}

func (*MsgsAck) CRC() uint32 { return 0x62d6b459 }

type BadMsgNotification struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	Code        int32
}

func (*BadMsgNotification) ImplementsBadMsgNotification() {}

func (*BadMsgNotification) CRC() uint32 { return 0xa7eff811 }

type BadServerSalt struct {
	BadMsgID    int64
	BadMsgSeqNo int32
	ErrorCode   int32
	NewSalt     int64
}

func (*BadServerSalt) ImplementsBadMsgNotification() {}

func (*BadServerSalt) CRC() uint32 { return 0xedab447b }

// msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;

type MsgResendReq struct {
	MsgIds []int64
}

func (*MsgResendReq) CRC() uint32 { return 0x7d861a08 }

type MsgsStateReq struct {
	MsgIds []int64
}

func (*MsgsStateReq) CRC() uint32 { return 0xda69fb52 }

type MsgsStateInfo struct {
	ReqMsgId int64
	Info     []byte
}

func (*MsgsStateInfo) CRC() uint32 { return 0x04deb57d }

type MsgsAllInfo struct {
	MsgIds []int64
	Info   []byte
}

func (*MsgsAllInfo) CRC() uint32 { return 0x8cc0d131 }

type MsgsDetailedInfo struct {
	MsgId       int64
	AnswerMsgId int64
	Bytes       int32
	Status      int32
}

func (*MsgsDetailedInfo) CRC() uint32 { return 0x276d3ec6 }

type MsgsNewDetailedInfo struct {
	AnswerMsgId int64
	Bytes       int32
	Status      int32
}

func (*MsgsNewDetailedInfo) CRC() uint32 { return 0x809db6df }

type ServerDHParams interface {
	tl.Object
	ImplementsServerDHParams()
}

type SetClientDHParamsAnswer interface {
	tl.Object
	ImplementsSetClientDHParamsAnswer()
}

func init() {
	tl.RegisterObjects(
		&ResPQ{},
		&PQInnerData{},
		&ServerDHParamsFail{},
		&ServerDHParamsOk{},
		&ServerDHInnerData{},
		&ClientDHInnerData{},
		&DHGenOk{},
		&DHGenRetry{},
		&DHGenFail{},
		&RpcResult{},
		&RpcError{},
		&RpcAnswerUnknown{},
		&RpcAnswerDroppedRunning{},
		&RpcAnswerDropped{},
		&FutureSalt{},
		&FutureSalts{},
		&Pong{},
		// &destroy_session_ok{}
		// &destroy_session_none{}
		&NewSessionCreated{},
		&MessageContainer{},
		&MsgCopy{},
		&GzipPacked{},
		&MsgsAck{},
		&BadMsgNotification{},
		&BadServerSalt{},
		&MsgResendReq{},
		&MsgsStateReq{},
		&MsgsStateInfo{},
		&MsgsAllInfo{},
		&MsgsDetailedInfo{},
		&MsgsNewDetailedInfo{},
	)
}
