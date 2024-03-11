// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/quenbyako/ext/slices"
	"github.com/xelaj/tl"

	"github.com/xelaj/mtproto/internal/mode"
	"github.com/xelaj/mtproto/internal/objects"
	"github.com/xelaj/mtproto/internal/payload"
	"github.com/xelaj/mtproto/internal/transport"
	"github.com/xelaj/mtproto/internal/utils"
)

type responseChanMsg struct {
	data []byte
	err  error
}

type MTProto struct {
	transport transport.Transport

	// каналы, которые ожидают ответа rpc. ответ записывается в канал и удаляется
	//
	// не используется sync.Map, поскольку при отмене контекста или ошибке нужно
	// отправить всем запросам что клиент умер.
	chanMux          sync.Mutex
	responseChannels map[payload.MsgID]chan<- responseChanMsg
}

type customHandlerFunc = func(i []byte) error

func NewUnencryptedTransport(m mode.Mode) transport.Transport {
	return transport.NewUnencrypted(m)
}

func NewEncryptedTransport(m mode.Mode, key [256]byte, salt uint64) transport.Transport {
	return transport.NewTransport(m, key, salt)
}

func NewMode(host string) (mode.Mode, error) {
	t, err := transport.NewTCP(host)
	if err != nil {
		return nil, err
	}

	return mode.Intermediate(t, true)
}

func New(t transport.Transport) (m *MTProto, err error) {
	m = &MTProto{
		transport:        t,
		responseChannels: make(map[payload.MsgID]chan<- responseChanMsg),
	}

	return m, nil
}

// NewAndRun initiates and spawns goroutine that calls client.Run method.
//
// IMPORTANT: This function is highly not recommended to use in real code, if
// you want to use it — probably there is something wrong with your architecture
// design. Think about this function like syntax sugar or context.TODO: it's
// worth it to manually call client.Run, cause you will have more control and
// code will look less implicit (obscurity was a giant problem of xelaj/mtproto
// 1.0, don't repeat our mistakes please 🙏)
func NewAndRun(ctx context.Context, t transport.Transport, serverRequestHandler customHandlerFunc, finalize func(err error)) (m *MTProto, err error) {
	client, err := New(t)
	if err != nil {
		return nil, err
	}

	go func() {
		err := client.Run(ctx, serverRequestHandler)
		if errors.Is(err, context.Canceled) {
			// resetting error, cause cancellation in this case is not an error
			err = nil
		}

		finalize(err)
	}()

	return client, nil
}

func (m *MTProto) Run(ctx context.Context, serverRequestHandler customHandlerFunc) error {
	// child context to cancell all jobs, if one of them fails.
	jobsCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := []func(context.Context) error{
		m.startReadingResponses(serverRequestHandler), // reading responses from the server
		m.startPinging, // keepalive pinging
	}

	var errs []error
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for i, job := range jobs {
		go func(i int, job func(context.Context) error) {
			err := omitContextErr(omitEOFErr(job(jobsCtx)))
			errs = append(errs, err)
			cancel()
			wg.Done()
		}(i, job)
	}
	wg.Wait()

	m.rejectAllRequests(errConnectionClosed)

	if err := errors.Join(errs...); err != nil {
		return err
	}

	return ctx.Err()
}

func (m *MTProto) MakeRequest(ctx context.Context, msg []byte) ([]byte, error) {
	return m.makeRequest(ctx, msg)
}

const defaultTimeout = 60 * time.Second // 60 seconds is maximum timeouts without pings

// startPinging pings the server that everything is fine, the client is online
// you just need to run and forget about it
func (m *MTProto) startPinging(ctx context.Context) error {
	ticker := time.NewTicker(defaultTimeout / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.pingDelay(ctx, utils.GenerateMessageId(), int32(defaultTimeout.Seconds())); err != nil {
				return fmt.Errorf("ping unsuccessful: %w", err)
			}
		}
	}
}

func (m *MTProto) startReadingResponses(serverRequestHandler customHandlerFunc) func(context.Context) error {
	if serverRequestHandler == nil {
		serverRequestHandler = defaultUnknownMsgHandler
	}

	return func(ctx context.Context) error {
		for {
			id, response, ack, err := m.transport.ReadMsg(ctx)
			if err == nil {
				// continue
			} else if e := new(transport.ErrCode); errors.As(err, e) {
				return err
			} else if errors.Is(err, transport.ErrCancelledAllRequests) {
				m.rejectAllRequests(errRetryRequest)
				continue
			} else if errorsIsAny(err,
				io.EOF,
				context.Canceled,
				context.DeadlineExceeded,
			) {
				return err
			} else {
				return fmt.Errorf("reading message: %w", err)
			}

			ids, err := m.processResponse(id, response, ack, serverRequestHandler)
			if err != nil {
				return fmt.Errorf("processing response: %w", err)
			}

			if len(ids) > 0 {
				if err := m.transport.Ack(ctx, ids); err != nil {
					return fmt.Errorf("sending ack: %w", err)
				}
			}
		}
	}
}

var errRetryRequest = errors.New("retry request")
var errConnectionClosed = errors.New("connection closed")

func (m *MTProto) processResponse(msgID payload.MsgID, msg []byte, needAck bool, serverRequestHandler customHandlerFunc) (processed []payload.MsgID, err error) {
	var data tl.Object
	err = objects.Unmarshal(msg, &data)
	if err != nil {
		err := fmt.Errorf("unmarshaling: %w", err)
		if len(msg) < 4 {
			return nil, err
		}

		// sometimes there may be messages that are not directly related to
		// MTProto, but to api (this is a valid TL object, but not described in
		// the TL schema of the protocol.
		//
		// We process these messages outside the connection, because we leave it
		// to Durov's conscience that it is a valid TL object and not a json or
		// xml.
		//
		//!IMPORTANT: we do NOT use reconciliation on tl.ErrObjectNotRegistered,
		// since only the difference of the original crc code is allowed, and no
		// others. (errors inside the object are hendled here as well)
		if _, ok := objects.Registry.ConstructObject(binary.LittleEndian.Uint32(msg)); !ok {
			if err := serverRequestHandler(msg); err != nil {
				return nil, err
			}

			return []payload.MsgID{msgID}, nil
		}

		return nil, err
	}

	switch data := data.(type) {
	case *objects.MessageContainer:
		ids := make([]payload.MsgID, 0, len(data.Content))
		for _, v := range data.Content {
			ack, err := m.processResponse(v.ID, v.Msg, v.NeedToAck(), serverRequestHandler)
			if err != nil {
				return nil, err
			}
			if v.NeedToAck() {
				ids = append(ids, ack...)
			}
		}
		return append([]payload.MsgID{msgID}, ids...), nil

	case *objects.Pong, *objects.MsgsAck:
		// игнорим, пришло и пришло, че бубнить то
		if needAck {
			return []payload.MsgID{msgID}, nil
		}
		return []payload.MsgID{}, nil

	case *objects.BadMsgNotification:
		return nil, BadMsgErrorFromNative(data)

	case *objects.RpcResult:
		// if firstCRC := binary.LittleEndian.Uint32(data.Obj); firstCRC == objects.CrcGzipPacked {
		// 	if data.Obj, err = objects.UnzipObject(data.Obj, true); err != nil {
		// 		return nil, fmt.Errorf("unzipping message %v: %w", data.ReqMsgID, err)
		// 	}
		// }

		found := m.writeRPCResponse(payload.MsgID(data.ReqMsgID), data.Obj)
		if !found {
			return nil, fmt.Errorf("msgID not found: %v", data.ReqMsgID)
		}

		return []payload.MsgID{msgID}, nil

	case *objects.NewSessionCreated:
		// pp.Println("uniq id", data.UniqueID)
		// s, err := m.session.Load()
		// if err != nil {
		// return nil, fmt.Errorf("loading session: %w", err)
		// }
		// s.Salt = data.ServerSalt
		// if err := m.session.Store(s); err != nil {
		// return nil, fmt.Errorf("storing session: %w", err)
		// }
		//
		return ackIfNecessary(nil, msgID, needAck), nil

	default:
		if err := serverRequestHandler(msg); err != nil {
			return nil, err
		}

		return ackIfNecessary(nil, msgID, needAck), nil
	}
}

func handleReadMsg(err error) error {
	if err == nil {
		return nil
	} else if e := new(transport.ErrCode); errors.As(err, e) {
		return ErrResponseCode{Code: int(*e)}
	} else if errorsIsAny(err,
		io.EOF,
		context.Canceled,
		context.DeadlineExceeded,
	) {
		return err
	}

	return fmt.Errorf("reading message: %w", err)
}

func errorsIsAny(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (m *MTProto) ping(ctx context.Context, pingID int64) error {
	msg, err := tl.Marshal(&objects.PingParams{PingID: pingID})
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	if _, err := m.sendPacket(ctx, msg, false); err != nil {
		return fmt.Errorf("sending: %w", err)
	}

	return nil
}

func (m *MTProto) pingDelay(ctx context.Context, pingID int64, delay int32) error {
	msg, err := tl.Marshal(&objects.PingDelayDisconnectParams{PingID: pingID, DisconnectDelay: delay})
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	if _, err := m.sendPacket(ctx, msg, false); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

// ackIfNecessary appends message id into sorted slice of message ids, if
// message id says that it needs to be acked
//
// input slice is always sorted as well as output slice. there is no duplicates
// in output slice
func ackIfNecessary(ids []payload.MsgID, id payload.MsgID, needToAck bool) []payload.MsgID {
	if ids == nil {
		ids = make([]payload.MsgID, 0)
	}

	if _, ok := slices.BinarySearch(ids, id); !ok && needToAck {
		ids = slices.AddSorted(ids, id)
	}

	return ids
}

func defaultUnknownMsgHandler(i []byte) error {
	var crcRaw [4]byte
	copy(crcRaw[:], i[:min(4, len(i))])
	crc := binary.LittleEndian.Uint32(crcRaw[:])

	return fmt.Errorf("got nonsystem message from server with crc 0x%08x", crc)
}
