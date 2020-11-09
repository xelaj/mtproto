package mtproto

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	bus "github.com/asaskevich/EventBus"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto/encoding/tl"
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
	sessionID  int64

	// общий мьютекс
	mutex         *sync.Mutex
	pending       map[int64]pendingRequest
	idsToAck      map[int64]struct{}
	idsToAckMutex sync.Mutex

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
	serviceChannel       chan []byte
	serviceModeActivated bool

	serverRequestHandlers []customHandlerFunc
}

type pendingRequest struct {
	// в response хранится тип который ожидаем получить
	// если он nil, то эта структурка не попадет в мапу pendingRequests
	response interface{}
	echan    chan error
}

type customHandlerFunc = func(i interface{}) bool

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

	m.sessionID = utils.GenerateSessionID()
	m.serviceChannel = make(chan []byte)
	m.publicKey = c.PublicKey
	m.serverRequestHandlers = make([]customHandlerFunc, 0)
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
	go m.startReadingResponses(ctx)

	// get new authKey if need
	if !m.encrypted {
		println("not encrypted, creating auth key")
		err = m.makeAuthKey()
		fmt.Println("authkey status:", err)
		if err != nil {
			return errors.Wrap(err, "making auth key")
		}
	}

	// start goroutines
	m.mutex = &sync.Mutex{}
	m.pending = make(map[int64]pendingRequest)
	// start keepalive pinging
	go m.startPinging(ctx)

	return nil
}

func (m *MTProto) makeRequest(req tl.Object, resp interface{}) error {
	err := m.sendPacket(req, resp)
	// если пришел ответ типа badServerSalt, то отправляем данные заново
	if errors.As(err, &serialize.ErrorSessionConfigsChanged{}) {
		return m.makeRequest(req, resp)
	}

	return err
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
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			_, err := m.Ping(0xCADACADA)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (m *MTProto) startReadingResponses(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := m.readFromConn(ctx)
			if err != nil {
				panic(err)
			}

			response, err := m.decodeRecievedData(data)
			if err != nil {
				panic(err)
			}

			if m.serviceModeActivated {
				fmt.Println("got service message")
				m.serviceChannel <- response.GetMsg()
				fmt.Println("service message pushed")
				continue
			}

			// NOTE:
			// Зачем сюда передавать m.msgId, m.seqNo
			// если у сообщения есть методы GetMsgID() и GetSeqNo()?
			if err := m.processResponse(
				atomic.LoadInt64(&m.msgId),
				atomic.LoadInt32(&m.seqNo),
				response.GetMsg(),
			); err != nil {
				panic(err)
			}
		}
	}
}

func (m *MTProto) processResponse(msgID int64, seqNo int32, data []byte) error {
	object, err := tl.DecodeRegistered(data)
	if err != nil {
		return fmt.Errorf("decode base message: %w", err)
	}

	switch message := object.(type) {
	case *serialize.RpcResult:
		m.mutex.Lock()
		req, found := m.pending[message.ReqMsgID]
		if !found {
			m.mutex.Unlock()
			fmt.Printf("pending request for messageID %d not found\n", message.ReqMsgID)
			break
		}
		delete(m.pending, message.ReqMsgID)
		m.mutex.Unlock()

		rpcMessageObject, err := tl.DecodeRegistered(message.Payload)
		if err != nil {
			// если не смогли заанмаршалить в зареганный тип
			// пробуем анмаршалить в тип прокинутый юзером
			//
			// Такое случается потому что DecodeRegistered (в отличие от Decode) не умеет
			// анмаршалить CrcVector, но можно его научить
			req.echan <- tl.Decode(message.Payload, req.response)
			break
		}

		// джедайские трюки
		switch rpcMessage := rpcMessageObject.(type) {
		case *serialize.GzipPacked:
			req.echan <- tl.Decode(rpcMessage.PackedData, req.response)
		case *serialize.RpcError:
			req.echan <- rpcMessage
		default: // если в rpc хз что, то анмаршалим его пейлоад в тот тип который запросил юзер

			// NOTE:
			// мб сделать свитч чисто по CRC чтобы убрать повторный анмаршал?
			// или установить значение rpcMessageObject в req.response через reflect?
			req.echan <- tl.Decode(message.Payload, req.response)
		}

	case *serialize.MessageContainer:
		println("MessageContainer")
		for _, v := range *message {
			// NOTE:
			// Зачем сюда передавать m.msgId, m.seqNo
			// если у сообщения есть методы GetMsgID() и GetSeqNo()?
			err := m.processResponse(v.MsgID, v.SeqNo, v.GetMsg())
			if err != nil {
				return errors.Wrap(err, "processing item in container")
			}
		}

	case *serialize.BadServerSalt:
		m.serverSalt = message.NewSalt
		err := m.SaveSession()
		dry.PanicIfErr(err)

		m.mutex.Lock()
		for _, v := range m.pending {
			v.echan <- &serialize.ErrorSessionConfigsChanged{}
		}
		m.mutex.Unlock()

	case *serialize.NewSessionCreated:
		println("session created")
		m.serverSalt = message.ServerSalt
		err := m.SaveSession()
		if err != nil {
			panic(err)
		}

	case *serialize.Pong:
		// игнорим, пришло и пришло, че бубнить то

	case *serialize.MsgsAck:
		for _, id := range message.MsgIds {
			m.gotAck(id)
		}

	case *serialize.BadMsgNotification:
		// NOTE:
		// что-то сделать с этим)
		panic(message)
		// return BadMsgErrorFromNative(message)

	default:
		panic(fmt.Sprintf("type %T not handled", message))
	}

	if (seqNo & 1) != 0 {
		// NOTE:
		// похоже MsgsAck можно кидать Ack на несколько сообщений сразу
		// Мб отправлять их батчами для меньшего жора сети?
		err = m.MakeRequest(&serialize.MsgsAck{MsgIds: []int64{int64(msgID)}}, nil)
		if err != nil {
			return errors.Wrap(err, "sending ack")
		}
	}

	return nil
}
