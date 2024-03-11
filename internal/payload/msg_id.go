// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math"
	"time"
)

// MsgID is a unique identifier of a message in MTProto protocol.
//
// Scheme of message_id (showed):
//
//	|1-31|32-61|62|63|
//	|A   |B    |C |D |
//
// Where:
//
//   - A: Approximately equal current unix time
//   - B: Any random unique 30-bit number
//   - C: Indicates the message initiator: 0 — for client-initiated (request or
//     server response), 1 — for server-initiated (notification or client
//     response to server request).
//   - D: message side (0 means client sent message, 1 means server did)
//
// More info:
// https://core.telegram.org/mtproto/description#message-identifier-msg-id
type MsgID uint64

const msgIDTimeOffset = 32

func (m MsgID) Time() time.Time      { return time.Unix(int64(m>>msgIDTimeOffset), 0) }
func (m MsgID) UniqueID() uint32     { return uint32(m & uniqueIDMask) }
func (m MsgID) Initiator() Initiator { return Initiator((m >> 1) & 1) }
func (m MsgID) Side() Side           { return Side(m & 1) }

const uniqueIDMask MsgID = math.MaxUint32 & ^(0b11)

// GenerateMessageId essentially gives unix timestamp, but in a horribly
// specific way.
//
// See [MsgID] for more info how it works
func GenerateMessageID(now time.Time, initiator Initiator, msgSide Side) MsgID {
	return generateMessageID(now, rand.Reader, initiator, msgSide)
}

const (
	randomnessInMsgID = msgIDTimeDependent

	msgIDRandom = iota
	msgIDTimeDependent
	msgIDIterative
)

// might panic if there is some error in rand reader
func generateMessageID(now time.Time, random io.Reader, initiator Initiator, msgSide Side) (id MsgID) {
	// ! WARN: we can have a collision, if two message ids will be generated in
	// same time, that's why we are using random lower bits. However, telegram
	// servers expecting to have INCREMENTAL message ids with incremental seqno,
	// (no idea how it works on server side, but it is what it is)
	var uniqID uint32
	switch randomnessInMsgID {
	case msgIDRandom:
		const uint32ByteSize = 4
		rawUniqueID := make([]byte, uint32ByteSize)
		if _, err := random.Read(rawUniqueID); err != nil {
			panic(err)
		}
		uniqID = binary.LittleEndian.Uint32(rawUniqueID)

	case msgIDTimeDependent:
		// 1 billion nanoseconds in 1 second, used to not include seconds in
		// fractional part
		const billion = 1e9
		uniqID = uint32(now.UnixNano() % billion)

	case msgIDIterative:
		panic("unimplemented")

	default:
		panic("unreachable")
	}

	id |= MsgID(now.Unix()) << 32
	id |= MsgID(uniqID) & uniqueIDMask
	id |= MsgID(initiator&1) << 1
	id |= MsgID(msgSide & 1)

	return id
}

type Initiator int8

const (
	InitiatorClient Initiator = iota
	InitiatorServer
)

type Side int8

const (
	SideClient Side = iota
	SideServer
)

func (s Side) String() string {
	if s&1 == 1 {
		return "server"
	}
	return "client"
}

// Initial vector extension for AES key. See more here:
// https://core.telegram.org/mtproto/description#defining-aes-key-and-initialization-vector
//
// if side is odd (server side) — returns 8
func (s Side) X() int { return 8 * (int(s) & 1) } //nolint:gomnd // come on

// Invert returns Side for decryption.
func (s Side) Invert() Side { return s ^ 1 } // flips bit, so 0 becomes 1, 1 becomes 0
