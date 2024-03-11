// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

// utility functions that do not depend much on declared structures, but are
// used in many places.

package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"time"
)

// GenerateMessageId essentially gives unix timestamp, but in a horribly
// specific way
//
// Scheme of message_id (showed):
//
//	|1-31|32-61|62|63|
//	|A   |B    |C |D |
//
// Where:
//
//   - A: Approximately equal curent unix time
//   - B: Any random unique 30-bit number
//   - C: Indicates the message initiator: 0 — for client-initiated (request or
//     server response), 1 — for server-initiated (notification or client
//     response to server request).
//   - D: message side (0 means client sent message, 1 means server did)
//
// More info:
// https://core.telegram.org/mtproto/description#message-identifier-msg-id
func GenerateMessageId() int64 {
	const billion = 1000 * 1000 * 1000
	unixnano := time.Now().UnixNano()
	seconds := unixnano / billion
	nanoseconds := unixnano % billion
	return (seconds << 32) | (nanoseconds & -4)
}

func GenerateSessionID() uint64 {
	var raw [8]byte
	rand.Read(raw[:])

	return binary.LittleEndian.Uint64(raw[:])
}

func AuthKeyHash(key []byte) []byte { return Sha1Byte(key)[12:20] }
func Sha1Byte(input []byte) []byte  { r := sha1.Sum(input); return r[:] }
