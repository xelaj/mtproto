package mtproto

import (
	"context"
	"crypto/rsa"
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

const (
	appId   = 124100
	appHash = "3ecccc5a1ec554722c3c5bbd35eb14ec"
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

	msgsIdToResp  map[int64]chan serialize.TL
	idsToAck      map[int64]struct{}
	idsToAckMutex sync.Mutex

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
		m.addr = "149.154.167.50:443"
		m.encrypted = false
	} else {
		return nil, errors.Wrap(err, "loading session")
	}

	m.sessionId = utils.GenerateSessionID()
	m.serviceChannel = make(chan serialize.TL)
	m.publicKey = c.PublicKey
	m.responseChannels = make(map[int64]chan serialize.TL)
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
func (m *MTProto) makeRequest(data serialize.TL) (serialize.TL, error) {
	println("sending packet " + reflect.TypeOf(data).String())

	resp, err := m.sendPacketNew(data)
	if err != nil {
		return nil, errors.Wrap(err, "sending message")
	}
	response := <-resp

	if _, ok := response.(*serialize.ErrorSessionConfigsChanged); ok {
		// если пришел ответ типа badServerSalt, то отправляем данные заново
		return m.makeRequest(data)
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

				pp.Println("got", response)

				if m.serviceModeActivated {
					m.serviceChannel <- response
				} else {
					err = m.processResponse(int(m.msgId), int(m.seqNo), response)
					dry.PanicIfErr(err)
				}
			}
		}
	}()
}

func (m *MTProto) processResponse(msgId, seqNo int, data serialize.TL) error {
	switch message := data.(type) {
	case *serialize.MessageContainer:
		pp.Println("MessageContainer")
		for _, v := range *message {
			err := m.processResponse(int(v.MsgID), int(v.SeqNo), v.Msg)
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
		pp.Println("session created")
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
		resp, err := m.makeRequest(&serialize.MsgsAck{[]int64{int64(msgId)}})
		pp.Println(resp)
		if err != nil {
			return errors.Wrap(err, "sending ack")
		}
	}

	return nil
}

/*
//! DEPRECATED
func (m *MTProto) process(msgId int64, seqNo int32, data interface{}) interface{} {
	pp.Println("process:", data)
	switch converted := data.(type) {
	case TL_msg_container:
		data := data.(TL_msg_container).items
		for _, v := range data {
			m.process(v.msg_id, v.seq_no, v.data)
		}
		os.Exit(1)

	case TL_bad_server_salt:
		m.serverSalt = int64(binary.LittleEndian.Uint64(converted.new_server_salt))
		err := m.SaveSession()
		dry.PanicIfErr(err)
		m.mutex.Lock() // что это?
		for k, v := range m.msgsIdToAck {
			delete(m.msgsIdToAck, k)
			m.queueSend <- v
		}
		m.mutex.Unlock() // что это?

	case *TL_BadServerSalt:
		m.serverSalt = converted.NewSalt // СОЛЬ ЭТО LONG
		err := m.SaveSession()
		dry.PanicIfErr(err)
		m.mutex.Lock() // что это?
		for k, v := range m.msgsIdToAck {
			delete(m.msgsIdToAck, k)
			m.queueSend <- v
		}
		m.mutex.Unlock() // что это?

	case TL_new_session_created:
		m.serverSalt = int64(binary.LittleEndian.Uint64(converted.server_salt))
		err := m.SaveSession()
		dry.PanicIfErr(err)
		os.Exit(1)

	case TL_ping:
		data := data.(TL_ping)
		m.queueSend <- packetToSend{TL_pong{msgId, data.ping_id}, nil}
		os.Exit(1)

	case TL_pong:
		// (ignore)

	case TL_msgs_ack:
		data := data.(TL_msgs_ack)
		m.mutex.Lock() // что это?
		for _, v := range data.msgIds {
			delete(m.msgsIdToAck, v)
		}
		m.mutex.Unlock() // что это?
		os.Exit(1)

	case TL_rpc_result:
		x := m.process(msgId, seqNo, converted.obj)
		m.mutex.Lock() // что это?
		v, ok := m.msgsIdToResp[converted.req_msg_id]
		if ok {
			v <- x.(TL)
			close(v)
			delete(m.msgsIdToResp, converted.req_msg_id)
		}
		delete(m.msgsIdToAck, converted.req_msg_id)
		m.mutex.Unlock() // что это?
		os.Exit(1)

	case *TL_RpcResult:
		pp.Println("RPC result!")
		x := m.process(msgId, seqNo, converted.Obj)
		m.mutex.Lock() // что это?
		v, ok := m.msgsIdToResp[converted.ReqMsgID]
		if ok {
			v <- x.(TL)
			close(v)
			delete(m.msgsIdToResp, converted.ReqMsgID)
		} else {
			pp.Println(converted.ReqMsgID, m.msgsIdToResp)
			panic("msgID not found")
		}

		delete(m.msgsIdToAck, converted.ReqMsgID)
		m.mutex.Unlock() // что это?
		os.Exit(1)

	default:
		fmt.Println(FullStack())
		return data

	}
	pp.Println("processed")

	if (seqNo & 1) == 1 {
		m.queueSend <- packetToSend{TL_msgs_ack{[]int64{msgId}}, nil}
	}

	return nil
}
*/

/*


//! DEPRECATED
func (m *MTProto) GetContacts() error {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{TL_contacts_getContacts{""}, resp}
	x := <-resp
	list, ok := x.(TL_contacts_contacts)
	if !ok {
		return fmt.Errorf("RPC: %#v", x)
	}

	contacts := make(map[int32]TL_userContact)
	for _, v := range list.users {
		if v, ok := v.(TL_userContact); ok {
			contacts[v.id] = v
		}
	}
	fmt.Printf(
		"\033[33m\033[1m%10s    %10s    %-30s    %-20s\033[0m\n",
		"id", "mutual", "name", "username",
	)
	for _, v := range list.contacts {
		v := v.(TL_contact)
		fmt.Printf(
			"%10d    %10t    %-30s    %-20s\n",
			v.user_id,
			v.mutual,
			fmt.Sprintf("%s %s", contacts[v.user_id].first_name, contacts[v.user_id].last_name),
			contacts[v.user_id].username,
		)
	}

	return nil
}

//! DEPRECATED
func (m *MTProto) TestAnything() error {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{TL_contacts_getContacts{""}, resp}
	x := <-resp
	pp.Println(x)
	return nil
}

//! DEPRECATED
func (m *MTProto) SendMessage(user_id int32, msg string) error {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_sendMessage{
			TL_inputPeerContact{user_id},
			msg,
			rand.Int63(),
		},
		resp,
	}
	x := <-resp
	_, ok := x.(TL_messages_sentMessage)
	if !ok {
		return fmt.Errorf("RPC: %#v", x)
	}

	return nil
}
*/
