package serialize

// messages.go отвечает за имплементацию сериализации сообщений. зашифрованных и незашифрованных.

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"

	ige "github.com/xelaj/mtproto/aes_ige"
	"github.com/xelaj/mtproto/utils"
)

// CommonMessage это сообщение (зашифрованое либо открытое) которыми общаются между собой клиент и сервер
type CommonMessage interface {
	GetMsg() []byte
	GetMsgID() int
	GetSeqNo() int
}

type EncryptedMessage struct {
	Msg         []byte
	MsgID       int64
	AuthKeyHash []byte

	Salt      int64
	SessionID int64
	SeqNo     int32
	MsgKey    []byte
}

func (msg *EncryptedMessage) Serialize(client MessageInformator, requireToAck bool) ([]byte, error) {
	obj := serializePacket(client, msg.Msg, msg.MsgID, requireToAck)
	encryptedData, err := ige.Encrypt(obj, client.GetAuthKey())
	if err != nil {
		return nil, errors.Wrap(err, "encrypting")
	}

	buf := NewEncoder()
	buf.PutRawBytes(utils.AuthKeyHash(client.GetAuthKey()))
	buf.PutRawBytes(ige.MessageKey(obj))
	buf.PutRawBytes(encryptedData)

	return buf.Result(), nil
}

func DeserializeEncryptedMessage(data, authKey []byte) (*EncryptedMessage, error) {
	msg := new(EncryptedMessage)

	buf := NewDecoder(data)
	keyHash := buf.PopRawBytes(LongLen)
	if !bytes.Equal(keyHash, utils.AuthKeyHash(authKey)) {
		return nil, errors.New("wrong encryption key")
	}
	msg.MsgKey = buf.PopRawBytes(Int128Len) // msgKey это хэш от расшифрованного набора байт, последние 16 символов
	encryptedData := buf.PopRawBytes(len(data) - (LongLen + Int128Len))

	decrypted, err := ige.Decrypt(encryptedData, authKey, msg.MsgKey)
	if err != nil {
		return nil, errors.Wrap(err, "decrypting message")
	}
	buf = NewDecoder(decrypted)
	msg.Salt = buf.PopLong()
	msg.SessionID = buf.PopLong()
	msg.MsgID = buf.PopLong()
	msg.SeqNo = buf.PopInt()
	messageLen := buf.PopInt()

	if len(decrypted) < int(messageLen)-(LongLen+LongLen+LongLen+WordLen+WordLen) {
		return nil, fmt.Errorf("message is smaller than it's defining: have %v, but messageLen is %v", len(decrypted), messageLen)
	}

	mod := msg.MsgID & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("wrong bits of message_id: %d", mod)
	}

	// этот кусок проверяет валидность данных по ключу
	trimed := decrypted[0 : 32+messageLen] // суммарное сообщение, после расшифровки
	if !bytes.Equal(dry.Sha1Byte(trimed)[4:20], msg.MsgKey) {
		return nil, errors.New("wrong message key, can't trust to sender")
	}
	msg.Msg = buf.PopRawBytes(int(messageLen))

	return msg, nil
}

func (msg *EncryptedMessage) GetMsg() []byte {
	return msg.Msg
}

func (msg *EncryptedMessage) GetMsgID() int {
	return int(msg.MsgID)
}

func (msg *EncryptedMessage) GetSeqNo() int {
	return int(msg.SeqNo)
}

type UnencryptedMessage struct {
	Msg   []byte
	MsgID int64
}

func (msg *UnencryptedMessage) Serialize(client MessageInformator) ([]byte, error) {
	buf := NewEncoder()
	// authKeyHash, always 0 if unencrypted
	buf.PutLong(0)
	buf.PutLong(msg.MsgID)
	buf.PutInt(int32(len(msg.Msg)))
	buf.PutRawBytes(msg.Msg)
	return buf.Result(), nil
}

func DeserializeUnencryptedMessage(data []byte) (*UnencryptedMessage, error) {
	msg := new(UnencryptedMessage)
	buf := NewDecoder(data)
	_ = buf.PopRawBytes(LongLen) // authKeyHash, always 0 if unencrypted

	msg.MsgID = buf.PopLong()

	mod := msg.MsgID & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("Wrong bits of message_id: %#v", uint64(mod))
	}

	messageLen := buf.PopUint()
	if len(data)-(LongLen+LongLen+WordLen) != int(messageLen) {
		fmt.Println(len(data), int(messageLen), int(messageLen+(LongLen+LongLen+WordLen)))
		return nil, fmt.Errorf("message not equal defined size: have %v, want %v", len(data), messageLen)
	}

	msg.Msg = buf.GetRestOfMessage()

	// TODO: в мтпрото объекте изменить msgID и задать seqNo 0
	return msg, nil
}

func (msg *UnencryptedMessage) GetMsg() []byte {
	return msg.Msg
}

func (msg *UnencryptedMessage) GetMsgID() int {
	return int(msg.MsgID)
}

func (msg *UnencryptedMessage) GetSeqNo() int {
	return 0
}

//------------------------------------------------------------------------------------------

// MessageInformator нужен что бы отдавать информацию о текущей сессии для сериализации сообщения
// по факту это *MTProto структура
type MessageInformator interface {
	GetSessionID() int64
	GetLastSeqNo() int32
	GetServerSalt() int64
	GetAuthKey() []byte
	MakeRequest(msg TL) (TL, error)
}

func serializePacket(client MessageInformator, msg []byte, messageID int64, requireToAck bool) []byte {
	buf := NewEncoder()

	saltBytes := make([]byte, LongLen)
	binary.LittleEndian.PutUint64(saltBytes, uint64(client.GetServerSalt()))
	buf.PutRawBytes(saltBytes)
	buf.PutLong(client.GetSessionID())
	buf.PutLong(messageID)
	if requireToAck { // не спрашивай, как это работает
		buf.PutInt(client.GetLastSeqNo() | 1) // почему тут добавляется бит не ебу
	} else {
		buf.PutInt(client.GetLastSeqNo())
	}
	buf.PutInt(int32(len(msg)))
	buf.PutRawBytes(msg)
	return buf.buf
}
