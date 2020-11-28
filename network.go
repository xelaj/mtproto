// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/messages"
	"github.com/xelaj/mtproto/internal/mtproto/objects"
)

const (
	// если длина пакета больше или равн 127 слов, то кодируем 4 байтами, 1 это магическое число, оставшиеся 3 — дилна
	magicValueSizeMoreThanSingleByte = 0x7f
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

const (
	readTimeout = 300 * time.Second
)

func CatchResponseErrorCode(data []byte) error {
	if len(data) == 4 {
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

	m.msgId = int64(msg.GetMsgID())
	m.seqNo = int32(msg.GetSeqNo())
	mod := m.msgId & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("Wrong bits of message_id: %d", mod)
	}

	return msg, nil
}
