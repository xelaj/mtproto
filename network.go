package mtproto

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/serialize"
)

const (
	// если длина пакета больше или равн 127 слов, то кодируем 4 байтами, 1 это магическое число, оставшиеся 3 — дилна
	magicValueSizeMoreThanSingleByte = 0x7f
)

func isNullableResponse(t serialize.TL) bool {
	switch t.(type) {
	case /**serialize.Ping,*/ *serialize.Pong, *serialize.MsgsAck:
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
	buf := serialize.NewDecoder(data)
	authKeyHash := buf.PopRawBytes(serialize.DoubleLen)
	return binary.LittleEndian.Uint64(authKeyHash) != 0
}

func (m *MTProto) decodeRecievedData(data []byte) (serialize.CommonMessage, error) {
	// проверим, что это не код ошибки
	err := CatchResponseErrorCode(data)
	if err != nil {
		return nil, errors.Wrap(err, "Server response error")
	}

	var msg serialize.CommonMessage

	if IsPacketEncrypted(data) {
		msg, err = serialize.DeserializeEncryptedMessage(data, m.GetAuthKey())
	} else {
		msg, err = serialize.DeserializeUnencryptedMessage(data)
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
