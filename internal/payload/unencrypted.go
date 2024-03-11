// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Unencrypted struct {
	ID  MsgID
	Msg []byte
}

// compatible with MTProto 1.0 and 2.0
//
// More info: https://core.telegram.org/mtproto/description#unencrypted-message
func (msg *Unencrypted) Serialize() (res []byte) {
	res = make([]byte, 0, uMsgLenOffset+uMsgLenLen+len(msg.Msg))

	res = binary.LittleEndian.AppendUint64(res, 0)              // authKeyHash, always 0 if unencrypted
	res = binary.LittleEndian.AppendUint64(res, uint64(msg.ID)) // msg_id
	res = binary.LittleEndian.AppendUint32(res, uint32(len(msg.Msg)))

	return append(res, msg.Msg...)
}

const (
	uAuthKeyHashOffset = 0
	uAuthKeyHashLen    = 8
	uMsgIDOffset       = uAuthKeyHashOffset + uAuthKeyHashLen
	uMsgIDLen          = 8
	uMsgLenOffset      = uMsgIDOffset + uMsgIDLen
	uMsgLenLen         = 4
)

// More info: https://core.telegram.org/mtproto/description#unencrypted-message
func DeserializeUnencrypted(data []byte) (*Unencrypted, error) {
	if len(data) < uMsgLenOffset+uMsgLenLen {
		return nil, errors.New("message is too short")
	}

	var (
		le = binary.LittleEndian

		authKeyHash = le.Uint64(data[uAuthKeyHashOffset : uAuthKeyHashOffset+uAuthKeyHashLen])
		msgID       = le.Uint64(data[uMsgIDOffset : uMsgIDOffset+uMsgIDLen])
		msgLen      = le.Uint32(data[uMsgLenOffset : uMsgLenOffset+uMsgLenLen])
		body        = data[uMsgLenOffset+uMsgLenLen:]
	)

	if authKeyHash != 0 {
		return nil, errors.New("received unencrypted message with auth_key_hash != 0")
	}

	if len(body) != int(msgLen) {
		return nil, fmt.Errorf("message not equal defined size: have %v, want %v", len(body), int(msgLen))
	}

	return &Unencrypted{
		Msg: body,
		ID:  MsgID(msgID),
	}, nil
}
