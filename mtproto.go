// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"context"
	"crypto/rsa"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"

	"github.com/xelaj/mtproto/internal/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/messages"
	"github.com/xelaj/mtproto/internal/mtproto/objects"
	"github.com/xelaj/mtproto/internal/utils"
)

type MTProto struct {
	addr         string
	conn         *net.TCPConn
	stopRoutines context.CancelFunc // остановить ping, read, и подобные горутины
	routineswg   sync.WaitGroup     // WaitGroup что бы быть уверенным, что все рутины остановились

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

	msgsIdToResp  map[int64]chan tl.Object
	idsToAck      map[int64]null
	idsToAckMutex sync.Mutex

	// каналы, которые ожидают ответа rpc. ответ записывается в канал и удаляется
	responseChannels map[int64]chan tl.Object
	expectedTypes    map[int64][]reflect.Type // uses for parcing bool values in rpc result for example

	// идентификаторы сообщений, нужны что бы посылать и принимать сообщения.
	seqNo int32
	msgId int64

	// не знаю что это но как-то используется
	lastSeqNo int32

	// айдишники DC для КОНКРЕТНОГО Приложения и клиента. Может меняться, но фиксирована для
	// связки приложение+клиент
	dclist map[int]string

	// путь до файла токена сессии.
	tokensStorage string

	// один из публичных ключей telegram. нужен только для создания сессии.
	publicKey *rsa.PublicKey

	// serviceChannel нужен только на время создания ключей, т.к. это
	// не RpcResult, поэтому все данные отдаются в один поток без
	// привязки к MsgID
	serviceChannel       chan tl.Object
	serviceModeActivated bool

	//! DEPRECATED RecoverFunc используется только до того момента, когда из пакета будут убраны все паники
	RecoverFunc func(i any)
	// если задан, то в канал пишутся ошибки
	Warnings chan error

	serverRequestHandlers []customHandlerFunc
}

type customHandlerFunc = func(i any) bool

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
	m.serviceChannel = make(chan tl.Object)
	m.publicKey = c.PublicKey
	m.responseChannels = make(map[int64]chan tl.Object)
	m.expectedTypes = make(map[int64][]reflect.Type)
	m.serverRequestHandlers = make([]customHandlerFunc, 0)
	// копируем мапу, т.к. все таки дефолтный список нельзя менять, вдруг его использует несколько клиентов
	m.SetDCStorages(defaultDCList)

	m.resetAck()

	return m, nil
}

func (m *MTProto) SetDCStorages(in map[int]string) {
	if m.dclist == nil {
		m.dclist = make(map[int]string)
	}
	for k, v := range defaultDCList {
		m.dclist[k] = v
	}
}

// Stop останавливает текущее соединение
func (m *MTProto) Stop() error {
	m.stopRoutines()
	m.routineswg.Wait()

	err := m.conn.Close()
	if err != nil {
		return errors.Wrap(err, "closing connection")
	}

	// все остановили, погнали пересоздаваться
	return nil
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
	_, err = m.conn.Write(transportModeIntermediate[:])
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
	m.msgsIdToResp = make(map[int64]chan tl.Object)
	m.mutex = &sync.Mutex{}

	// start keepalive pinging
	m.startPinging(ctx)

	return nil
}

// отправить запрос
func (m *MTProto) makeRequest(data tl.Object, expectedTypes ...reflect.Type) (any, error) {
	resp, err := m.sendPacketNew(data, expectedTypes...)
	if err != nil {
		return nil, errors.Wrap(err, "sending message")
	}

	response := <-resp

	switch r := response.(type) {
	case *objects.RpcError:
		realErr := RpcErrorToNative(r)

		err = m.tryToProcessErr(realErr.(*ErrResponseCode))
		if err != nil {
			return nil, err
		}

		return m.makeRequest(data, expectedTypes...)

	case *errorSessionConfigsChanged:
		return m.makeRequest(data, expectedTypes...)

	}

	return tl.UnwrapNativeTypes(response), nil
}

// Disconnect is closing current TCP connection and stopping all routines like pinging, reading etc.
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
	m.routineswg.Add(1)
	ticker := time.Tick(time.Minute)
	go func() {
		defer m.recoverGoroutine()
		for {
			select {
			case <-ctx.Done():
				m.routineswg.Done()
				return
			case <-ticker:
				_, err := m.ping(0xCADACADA) //nolint:gomnd not magic
				if err != nil {
					m.warnError(errors.Wrap(err, "ping unsuccsesful"))
				}
			}
		}
	}()
}

func (m *MTProto) startReadingResponses(ctx context.Context) {
	m.routineswg.Add(1)
	go func() {
		defer m.recoverGoroutine()
		for {
			select {
			case <-ctx.Done():
				m.routineswg.Done()
				return
			default:
				data, err := m.readFromConn(ctx)
				if err != nil {
					m.warnError(errors.Wrap(err, "reading from connection"))
					break // select
				}

				response, err := m.decodeRecievedData(data)
				if err != nil {
					m.warnError(errors.Wrap(err, "decoding received data"))
					break // select
				}

				if m.serviceModeActivated {
					var obj tl.Object
					// сервисные сообщения ГАРАНТИРОВАННО в теле содержат TL.
					obj, err = tl.DecodeUnknownObject(response.GetMsg())
					if err != nil {
						m.warnError(errors.Wrap(err, "parsing object"))
						break
					}
					m.serviceChannel <- obj
					break
				}

				err = m.processResponse(response)
				if err != nil {
					m.warnError(errors.Wrap(err, "processing response"))
				}
			}
		}
	}()
}

func (m *MTProto) processResponse(msg messages.Common) error {
	var data tl.Object
	var err error
	if et, ok := m.expectedTypes[int64(msg.GetMsgID())]; ok && len(et) > 0 {
		data, err = tl.DecodeUnknownObject(msg.GetMsg(), et...)
	} else {
		data, err = tl.DecodeUnknownObject(msg.GetMsg())
	}
	if err != nil {
		return errors.Wrap(err, "unmarshaling response")
	}
	switch message := data.(type) {
	case *objects.MessageContainer:
		for _, v := range *message {
			err := m.processResponse(v)
			if err != nil {
				return errors.Wrap(err, "processing item in container")
			}
		}

	case *objects.BadServerSalt:
		m.serverSalt = message.NewSalt
		err := m.SaveSession()
		check(err)

		m.mutex.Lock()
		for _, v := range m.responseChannels {
			v <- &errorSessionConfigsChanged{}
		}
		m.mutex.Unlock()

	case *objects.NewSessionCreated:
		m.serverSalt = message.ServerSalt
		err := m.SaveSession()
		if err != nil {
			m.warnError(errors.Wrap(err, "saving session"))
		}

	case *objects.Pong:
		// игнорим, пришло и пришло, че бубнить то

	case *objects.MsgsAck:
		for _, id := range message.MsgIDs {
			m.gotAck(id)
		}

	case *objects.BadMsgNotification:
		pp.Println(message)
		panic(message) // for debug, looks like this message is important
		return BadMsgErrorFromNative(message)

	case *objects.RpcResult:
		obj := message.Obj
		if v, ok := obj.(*objects.GzipPacked); ok {
			obj = v.Obj
		}

		err := m.writeRPCResponse(int(message.ReqMsgID), obj)
		if err != nil {
			return errors.Wrap(err, "writing RPC response")
		}

	default:
		processed := false
		for _, f := range m.serverRequestHandlers {
			processed = f(message)
			if processed {
				break
			}
		}
		if !processed {
			m.warnError(errors.New("got nonsystem message from server: " + reflect.TypeOf(message).String()))
		}
	}

	if (msg.GetSeqNo() & 1) != 0 {
		_, err := m.MakeRequest(&objects.MsgsAck{MsgIDs: []int64{int64(msg.GetMsgID())}})
		if err != nil {
			return errors.Wrap(err, "sending ack")
		}
	}

	return nil
}

// tryToProcessErr пытается автоматически решить ошибку полученную от сервера. в случае успеха вернет nil,
// в случае если нет способа решить эту проблему, возвращается сама ошибка
// если в процессе решения появлиась еще одна ошибка, то она оборачивается в errors.Wrap, основная
// игнорируется (потому что гарантируется, что обработка ошибки надежна, и параллельная ошибка это что-то из
// ряда вон выходящее)
func (m *MTProto) tryToProcessErr(e *ErrResponseCode) error {
	switch e.Message {
	case "PHONE_MIGRATE_X":
		newIP, found := m.dclist[e.AdditionalInfo.(int)]
		if !found {
			return errors.Wrapf(e, "DC with id %v not found", e.AdditionalInfo)
		}
		err := m.Stop()
		if err != nil {
			return errors.Wrap(err, "stopping session")
		}

		m.addr = newIP

		err = m.CreateConnection()
		if err != nil {
			return errors.Wrap(err, "recreating session")
		}

		return nil

	default:
		return e
	}
}
