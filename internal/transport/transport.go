package transport

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/xelaj/tl"

	"github.com/xelaj/mtproto/internal/mode"
	"github.com/xelaj/mtproto/internal/payload"
	"github.com/xelaj/mtproto/internal/utils"
)

type Transport interface {
	Close() error
	WriteMsg(ctx context.Context, body []byte, initiator payload.Initiator) (payload.MsgID, error)
	Ack(ctx context.Context, ids []payload.MsgID) error
	ReadMsg(context.Context) (id payload.MsgID, body []byte, ack bool, err error)
}

type transport struct {
	mode      mode.Mode
	cipher    payload.Cipher
	key       [256]byte
	salt      uint64
	sessionID uint64
	clock     Clock
	role      payload.Side // client or server

	// mutex for sending messages only, reading is always sequential
	m sync.Mutex

	//  for detecting that reading call is only one
	isReading atomic.Bool

	// client side sequence number, server side is useless for tcp
	seqNo uint32
}

var _ Transport = (*transport)(nil)

func NewTransport(mode mode.Mode, key [256]byte, salt uint64) Transport {
	return &transport{
		mode:      mode,
		cipher:    payload.NewClientCipher(rand.Reader, key),
		key:       key,
		salt:      salt,
		sessionID: utils.GenerateSessionID(),
		clock:     ClockFunc(time.Now),
		role:      payload.SideClient,
	}
}

func (t *transport) Close() error { return t.mode.Close() }

func (t *transport) getSeqNo(ack bool) uint32 {
	defer func() { t.seqNo += 2 }()
	if ack {
		return t.seqNo
	} else {
		return t.seqNo | 1
		// return taken
	}
}

func (t *transport) WriteMsg(ctx context.Context, body []byte, initiator payload.Initiator) (payload.MsgID, error) {
	return t.writeMsg(ctx, body, initiator, false)
}

// https://core.telegram.org/mtproto/service_messages_about_messages#acknowledgment-of-receipt
func (t *transport) Ack(ctx context.Context, ids []payload.MsgID) error {
	body := encodeAck(ids)

	_, err := t.writeMsg(ctx, body, payload.Initiator(t.role), true)

	return err
}

func (t *transport) writeMsg(_ context.Context, body []byte, initiator payload.Initiator, isAck bool) (payload.MsgID, error) {
	t.m.Lock()
	defer t.m.Unlock()

	id := payload.GenerateMessageID(t.clock.Now(), initiator, t.role)

	seqno := t.getSeqNo(isAck)

	envelope := payload.BuildEnvelope(t.salt, t.sessionID, seqno, id, body, rand.Reader)

	data, err := t.cipher.Encrypt([8]byte{}, envelope)
	if err != nil {
		panic(err)
	}

	if err := t.mode.WriteMsg(data); err != nil {
		return 0, err
	}

	return id, nil
}

func (t *transport) ReadMsg(ctx context.Context) (id payload.MsgID, body []byte, ack bool, err error) {
	// locking...
	if !t.isReading.CompareAndSwap(false, true) {
		return 0, nil, false, errors.New("ReadMsg called from multiple goroutines")
	}
	defer func() { t.isReading.Store(false) }()

	// and now reading

	data, err := t.mode.ReadMsg(ctx)
	if err != nil {
		switch err {
		case io.EOF, context.Canceled:
			return 0, nil, false, err
		default:
			return 0, nil, false, errors.Wrap(err, "reading message")
		}
	}

	// checking that response is not error code
	if len(data) == 4 {
		code := int(binary.LittleEndian.Uint32(data))
		return 0, nil, false, ErrCode(code)
	}

	if !isPacketEncrypted(data) {
		return 0, nil, false, errors.New("received unencrypted message on encrypted transport")
	}

	decrypted, err := t.cipher.Decrypt(data)
	if err != nil {
		return 0, nil, false, err
	}

	e, err := payload.DeserializeEnvelope(decrypted)
	if err != nil {
		return 0, nil, false, err
	}

	// обработка особенных сервисных сообщений (обновление соли, проблемы с seqno и прочие)
	if binary.LittleEndian.Uint32(e.Msg) == crcBadServerSalt {
		var bsdSaltMsg badServerSalt
		if err := tl.Unmarshal(e.Msg, &bsdSaltMsg); err == nil {
			t.salt = uint64(bsdSaltMsg.NewSalt)
			return 0, nil, false, ErrCancelledAllRequests
		} else {
			// ok, failed to parse, just throwing to reader
			panic(err)
		}
	}

	return e.MsgID, e.Msg, e.SeqNo|1 > 0, nil
}

type transportUnencrypted struct {
	role  payload.Side
	clock Clock
	mode  mode.Mode

	//  for detecting that reading call is only one
	isReading atomic.Bool
}

var _ Transport = (*transportUnencrypted)(nil)

func NewUnencrypted(mode mode.Mode) Transport {
	return &transportUnencrypted{
		role:  payload.SideClient,
		clock: ClockFunc(time.Now),
		mode:  mode,
	}
}

func (t *transportUnencrypted) Close() error { return t.mode.Close() }

func (t *transportUnencrypted) WriteMsg(_ context.Context, body []byte, initiator payload.Initiator) (payload.MsgID, error) {
	msgID := payload.GenerateMessageID(t.clock.Now(), initiator, t.role)

	encoded := payload.Unencrypted{
		ID:  msgID,
		Msg: body,
	}

	if err := t.mode.WriteMsg(encoded.Serialize()); err != nil {
		return 0, err
	}

	return msgID, nil
}

func (t *transportUnencrypted) Ack(ctx context.Context, ids []payload.MsgID) error {
	// except for Telegram servers, for mtproto as a protocol, acknowledgments
	// in unecrypted messages might also be really usable
	panic("forbidden")
}

func (t *transportUnencrypted) ReadMsg(ctx context.Context) (id payload.MsgID, body []byte, ack bool, err error) {
	// locking...
	if !t.isReading.CompareAndSwap(false, true) {
		return 0, nil, false, errors.New("ReadMsg called from multiple goroutines")
	}
	defer func() { t.isReading.Store(false) }()

	// and now reading

	data, err := t.mode.ReadMsg(ctx)
	if err != nil {
		switch err {
		case io.EOF, context.Canceled:
			return 0, nil, false, err
		default:
			return 0, nil, false, errors.Wrap(err, "reading message")
		}
	}

	// checking that response is not error code
	if len(data) == 4 {
		code := int(binary.LittleEndian.Uint32(data))
		return 0, nil, false, ErrCode(code)
	}

	if isPacketEncrypted(data) {
		return 0, nil, false, errors.New("received encrypted message on unencrypted transport")
	}

	msg, err := payload.DeserializeUnencrypted(data)
	if err != nil {
		return 0, nil, false, errors.Wrap(err, "parsing message")
	}

	// https://core.telegram.org/mtproto/description#message-identifier-msg-id
	if side := msg.ID.Side(); side == t.role {
		return 0, nil, false, fmt.Errorf("received message from %v, expected from %v", side, t.role)
	}

	return msg.ID, msg.Msg, false, nil
}

func isPacketEncrypted(data []byte) bool {
	if len(data) < 8 {
		return false
	}
	authKeyHash := data[:8]
	return binary.LittleEndian.Uint64(authKeyHash) != 0
}

// See more: https://core.telegram.org/mtproto/mtproto-transports#transport-errors
type ErrCode uint32

func (e ErrCode) Error() string {
	switch e {
	case ErrAuthKeyNotFound:
		return "auth key not found"
	case ErrFlood:
		return "transport flood"
	case ErrInvalidDC:
		return "invalid DC"
	default:
		return fmt.Sprintf("code %v", int(e))
	}
}

const (
	ErrAuthKeyNotFound ErrCode = ^ErrCode(404) + 1
	ErrFlood           ErrCode = ^ErrCode(429) + 1
	ErrInvalidDC       ErrCode = ^ErrCode(444) + 1
)

var ErrCancelledAllRequests = errors.New("session config changed, need to retry all requests")

type Clock interface {
	Now() time.Time
}

type ClockFunc func() time.Time

func (c ClockFunc) Now() time.Time { return c() }
