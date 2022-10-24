package mtproto

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"
)

const (
	minDiff = 300 * time.Second
	maxDiff = -30 * time.Second
)

type MsgTyp uint8

const (
	MsgClient     MsgTyp = 0b00
	MsgServerResp MsgTyp = 0b01
	MsgServerUpd  MsgTyp = 0b11
)

// https://core.telegram.org/mtproto/description#message-identifier-msg-id
// 1-32:  current unix time
// 33-62: random bits
// 63:    for server-side only: 0 if msg is answer to request, 1 if server update
// 64:    1 if msg server-side, 0 if client-side
type MsgID uint64

func (m MsgID) At() time.Time { return time.Unix(int64(m)>>32, 0) }

func (m MsgID) IsValid(at time.Time) bool {
	t := m.At().Sub(at)
	return minDiff > t && t > maxDiff
}

func NewMsgID(at time.Time, typ MsgTyp) MsgID {
	if typ >= 4 {
		panic("invalid type")
	}

	timePart := uint64(at.Unix()) << 32
	randPart := randUint32() & (math.MaxUint32 - 3)
	return MsgID(timePart) | MsgID(randPart) | MsgID(typ)
}

func randUint32() uint32 {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		panic("unreachable") // to be sure that we will never get any problems
	}
	return binary.LittleEndian.Uint32(b)
}
