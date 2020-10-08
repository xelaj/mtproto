package serialize

// messages.go отвечает за имплементацию сериализации сообщений. зашифрованных и незашифрованных.

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"

	ige "github.com/xelaj/mtproto/aes_ige"
	"github.com/xelaj/mtproto/utils"
)

type EncryptedMessage struct {
	Msg         TL
	MsgID       int64
	AuthKeyHash []byte

	Salt      int64
	SessionID int64
	SeqNo     int32
	MsgKey    []byte
}

func (msg *EncryptedMessage) Serialize(client MessageInformator, requireToAck bool) []byte {
	obj := serializePacket(client, msg.Msg, msg.MsgID, requireToAck)
	encryptedData, err := ige.Encrypt(obj, client.GetAuthKey())
	dry.PanicIfErr(err)

	buf := NewEncoder()
	buf.PutRawBytes(utils.AuthKeyHash(client.GetAuthKey()))
	buf.PutRawBytes(ige.MessageKey(obj))
	buf.PutRawBytes(encryptedData)

	return buf.Result()
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
		return nil, fmt.Errorf("decrypting message: %w", err)
	}
	pp.Println("decoded", decrypted)
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
		return nil, fmt.Errorf("Wrong bits of message_id: %d", mod)
	}

	trimed := decrypted[0 : 32+messageLen] // суммарное сообщение, после расшифровки, с чет
	if !bytes.Equal(dry.Sha1Byte(trimed)[4:20], msg.MsgKey) {
		return nil, errors.New("Wrong message key, can't trust to sender")
	}
	msg.Msg = buf.PopObj()

	return msg, nil
	// TODO: мтпрото обновить msgID и seqNo
}

type UnencryptedMessage struct {
	Msg   TL
	MsgID int64
}

func (msg *UnencryptedMessage) Serialize(client MessageInformator) []byte {
	encodedMessage := msg.Msg.Encode()

	buf := NewEncoder()
	// authKeyHash, always 0 if unencrypted
	buf.PutLong(0)
	buf.PutLong(msg.MsgID)
	buf.PutInt(int32(len(encodedMessage)))
	buf.PutRawBytes(encodedMessage)
	return buf.Result()
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
		pp.Println(len(data), int(messageLen), int(messageLen+(LongLen+LongLen+WordLen)))
		return nil, fmt.Errorf("message not equal defined size: have %v, want %v", len(data), messageLen)
	}

	obj := buf.PopObj()

	msg.Msg = obj

	// TODO: в мтпрото объекте изменить msgID и задать seqNo 0
	return msg, nil
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

func serializePacket(client MessageInformator, msg TL, messageID int64, requireToAck bool) []byte {
	serializedMessage := msg.Encode()

	buf := NewEncoder()

	saltBytes := make([]byte, LongLen)
	binary.LittleEndian.PutUint64(saltBytes, uint64(client.GetServerSalt()))
	buf.PutRawBytes(saltBytes)
	pp.Println(saltBytes, fmt.Sprintf("%#v", uint64(client.GetServerSalt())))
	buf.PutLong(client.GetSessionID())
	buf.PutLong(messageID)
	if requireToAck { // не спрашивай, как это работает
		buf.PutInt(client.GetLastSeqNo() | 1) // почему тут добавляется бит не ебу
	} else {
		buf.PutInt(client.GetLastSeqNo())
	}
	buf.PutInt(int32(len(serializedMessage)))
	buf.PutRawBytes(serializedMessage)
	return buf.buf
}
