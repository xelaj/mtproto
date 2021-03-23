// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/internal/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/messages"
	"github.com/xelaj/mtproto/internal/mtproto/objects"
)

// проверяет, надо ли ждать от сервера пинга
func isNullableResponse(t tl.Object) bool {
	switch t.(type) {
	case /**objects.Ping,*/ *objects.Pong, *objects.MsgsAck:
		return true
	default:
		return false
	}
}

func CatchResponseErrorCode(data []byte) error {
	if len(data) == tl.WordLen {
		code := int(binary.LittleEndian.Uint32(data))
		return &ErrResponseCode{Code: code}
	}
	return nil
}

func IsPacketEncrypted(data []byte) bool {
	if len(data) < tl.DoubleLen {
		return false
	}
	authKeyHash := data[:tl.DoubleLen]
	return binary.LittleEndian.Uint64(authKeyHash) != 0
}

func (m *MTProto) decodeRecievedData(data []byte) (messages.Common, error) {
	// проверим, что это не код ошибки
	err := CatchResponseErrorCode(data)
	if err != nil {
		return nil, errors.Wrap(err, "Server response error")
	}

	var msg messages.Common

	if IsPacketEncrypted(data) {
		msg, err = messages.DeserializeEncrypted(data, m.GetAuthKey())
	} else {
		msg, err = messages.DeserializeUnencrypted(data)
	}
	if err != nil {
		return nil, errors.Wrap(err, "parsing message")
	}

	mod := msg.GetMsgID() & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("wrong bits of message_id: %d", mod)
	}

	return msg, nil
}
