package mtproto

import (
	"context"
	"crypto/rsa"
	"encoding/binary"
	"net"
	"reflect"
	"sync"
	"time"

	bus "github.com/asaskevich/EventBus"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto/serialize"
	"github.com/xelaj/mtproto/utils"
)

type MTProto struct {
	addr         string
	conn         *net.TCPConn
	stopRoutines context.CancelFunc // остановить ping, read, и подобные горутины

	// ключ авторизации. изменять можно только через setAuthKey
	authKey []byte

	// хеш ключа авторизации. изменять можно только через setAuthKey
	authKeyHash []byte

	// соль сессии
	serverSalt int64
	encrypted  bool
	sessionId  int64

	// общий мьютекс
	mutex *sync.Mutex

	// msgsIdDecodeAsVector показывает, что определенный ответ сервера нужно декодировать как
	// слайс. Это костыль, т.к. MTProto ЧЕТКО указывает, что ответы это всегда объекты, но
	// вектор (слайс) это как бы тоже объект. Из-за этого приходится четко указывать, что
	// сообщения с определенным msgID нужно декодировать как слайс, а не объект
	msgsIdDecodeAsVector map[int64]reflect.Type
	msgsIdToResp         map[int64]chan serialize.TL
	idsToAck             map[int64]struct{}
	idsToAckMutex        sync.Mutex

	// каналы, которые ожидают ответа rpc. ответ записывается в канал и удаляется
	responseChannels map[int64]chan serialize.TL

	// идентификаторы сообщений, нужны что бы посылать и принимать сообщения.
	seqNo int32
	msgId int64

	// не знаю что это но как-то используется
	lastSeqNo int32

	// пока непонятно для чего, кажется это нужно клиенту конкретно телеграма
	dclist map[int32]string

	// шина сообщений, используется для разных нотификаций, описанных в константах нотификации
	bus bus.Bus

	// путь до файла токена сессии.
	tokensStorage string

	// один из публичных ключей telegram. нужен только для создания сессии.
	publicKey *rsa.PublicKey

	// serviceChannel нужен только на время создания ключей, т.к. это
	// не RpcResult, поэтому все данные отдаются в один поток без
	// привязки к MsgID
	serviceChannel       chan serialize.TL
	serviceModeActivated bool
}

type Config struct {
	AuthKeyFile string
	ServerHost  string
	PublicKey   *rsa.PublicKey
}

func NewMTProto(c Config) (*MTProto, error) {
	m := new(MTProto)
	m.tokensStorage = c.AuthKeyFile

	err := m.LoadSession()
	if err == nil {
		m.encrypted = true
	} else if errs.IsNotFound(err) {
		m.addr = c.ServerHost
		m.encrypted = false
	} else {
		return nil, errors.Wrap(err, "loading session")
	}

	m.sessionId = utils.GenerateSessionID()
	m.serviceChannel = make(chan serialize.TL)
	m.publicKey = c.PublicKey
	m.responseChannels = make(map[int64]chan serialize.TL)
	m.msgsIdDecodeAsVector = make(map[int64]reflect.Type)
	m.resetAck()

	return m, nil
}

func (m *MTProto) CreateConnection() error {
	// connect
	tcpAddr, err := net.ResolveTCPAddr("tcp", m.addr)
	if err != nil {
		return errors.Wrap(err, "resolving tcp")
	}
	m.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return errors.Wrap(err, "dialing tcp")
	}

	// https://core.telegram.org/mtproto/mtproto-transports#intermediate
	_, err = m.conn.Write([]byte{0xee, 0xee, 0xee, 0xee})
	if err != nil {
		return errors.Wrap(err, "writing first byte")
	}

	ctx, cancelfunc := context.WithCancel(context.Background())
	m.stopRoutines = cancelfunc

	// start reading responses from the server
	m.startReadingResponses(ctx)

	// get new authKey if need
	if !m.encrypted {
		println("not encrypted, creating auth key")
		err = m.makeAuthKey()
		if err != nil {
			return errors.Wrap(err, "making auth key")
		}
	}

	// start goroutines
	m.msgsIdToResp = make(map[int64]chan serialize.TL)
	m.mutex = &sync.Mutex{}

	// start keepalive pinging
	m.startPinging(ctx)

	return nil
}

// отправить запрос
func (m *MTProto) makeRequest(data serialize.TL, as reflect.Type) (serialize.TL, error) {
	println("sending packet " + reflect.TypeOf(data).String())

	resp, err := m.sendPacketNew(data, as)
	if err != nil {
		return nil, errors.Wrap(err, "sending message")
	}
	response := <-resp

	if _, ok := response.(*serialize.ErrorSessionConfigsChanged); ok {
		// если пришел ответ типа badServerSalt, то отправляем данные заново
		return m.makeRequest(data, as)
	}
	if e, ok := response.(*serialize.RpcError); ok {
		return nil, RpcErrorToNative(e)
	}

	return response, nil
}

func (m *MTProto) Disconnect() error {
	// stop all routines
	m.stopRoutines()

	err := m.conn.Close()
	if err != nil {
		return errors.Wrap(err, "closing TCP connection")
	}

	// TODO: закрыть каналы

	// возвращаем в false, потому что мы теряем конфигурацию
	// сессии, и можем ее потерять во время отключения.
	m.encrypted = false

	return nil
}

// startPinging пингует сервер что все хорошо, клиент в сети
// нужно просто запустить
func (m *MTProto) startPinging(ctx context.Context) {
	ticker := time.Tick(60 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				_, err := m.Ping(0xCADACADA)
				dry.PanicIfErr(err)
			}
		}
	}()
}

func (m *MTProto) startReadingResponses(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := m.readFromConn(ctx)
				dry.PanicIfErr(err)

				response, err := m.decodeRecievedData(data)
				dry.PanicIfErr(err)

				if m.serviceModeActivated {
					// сервисные сообщения ГАРАНТИРОВАННО в теле содержат TL.
					decoder := serialize.NewDecoder(response.GetMsg())
					obj := decoder.PopObj()
					m.serviceChannel <- obj
				} else {
					err = m.processResponse(int(m.msgId), int(m.seqNo), response)
					dry.PanicIfErr(err)
				}
			}
		}
	}()
}

func (m *MTProto) processResponse(msgId, seqNo int, msg serialize.CommonMessage) error {
	// сначала декодируем исключения

	// TODO: может как-то поопрятней сделать? а то очень кринжово, функция занимается не тем, чем должна
	decoder := serialize.NewDecoder(msg.GetMsg())
	var data serialize.TL
	// если это ответ Rpc, то там может быть слайс вместо объекта, надо проверить указывали ли мы,
	// что ответ с этим MsgId нужно декодировать как слайс, а не объект
	if binary.LittleEndian.Uint32(msg.GetMsg()[:serialize.WordLen]) == serialize.CrcRpcResult {
		_ = decoder.PopCRC() // уже прочитали
		rpc := &serialize.RpcResult{}
		msgID := binary.LittleEndian.Uint64(msg.GetMsg()[serialize.WordLen : serialize.WordLen+serialize.LongLen])
		if typ, ok := m.msgsIdDecodeAsVector[int64(msgID)]; ok {
			rpc.DecodeFromButItsVector(decoder, typ)
			delete(m.msgsIdDecodeAsVector, int64(msgID))
		} else {
			rpc.DecodeFrom(decoder)
		}
		data = rpc
	} else {
		data = decoder.PopObj()
	}

	switch message := data.(type) {
	case *serialize.MessageContainer:
		println("MessageContainer")
		for _, v := range *message {
			err := m.processResponse(int(v.MsgID), int(v.SeqNo), v)
			if err != nil {
				return errors.Wrap(err, "processing item in container")
			}
		}

	case *serialize.BadServerSalt:
		m.serverSalt = message.NewSalt
		err := m.SaveSession()
		dry.PanicIfErr(err)

		m.mutex.Lock()
		for _, v := range m.responseChannels {
			v <- &serialize.ErrorSessionConfigsChanged{}
		}
		m.mutex.Unlock()

	case *serialize.NewSessionCreated:
		println("session created")
		m.serverSalt = message.ServerSalt
		err := m.SaveSession()
		dry.PanicIfErr(err)

	//case *serialize.Ping:
	//	resp, err := m.makeRequest(&TL_Pong{MsgID: int64(msgId), PingID: message.PingID})
	//	pp.Println(resp)
	//	if err != nil {
	//		return errors.Wrap(err, "processing ping")
	//	}

	case *serialize.Pong:
		// игнорим, пришло и пришло, че бубнить то

	case *serialize.MsgsAck:
		for _, id := range message.MsgIds {
			m.gotAck(id)
		}

	case *serialize.BadMsgNotification:
		pp.Println(message)
		panic("bad msg notification")

	case *serialize.RpcResult:
		obj := message.Obj
		if v, ok := obj.(*serialize.GzipPacked); ok {
			obj = v.Obj
		}

		err := m.writeRPCResponse(int(message.ReqMsgID), obj)
		if err != nil {
			return errors.Wrap(err, "writing RPC response")
		}

	default:
		panic("this is not system message: " + reflect.TypeOf(message).String())
	}

	if (seqNo & 1) != 0 {
		_, err := m.makeRequest(&serialize.MsgsAck{[]int64{int64(msgId)}}, nil)
		if err != nil {
			return errors.Wrap(err, "sending ack")
		}
	}

	return nil
}
