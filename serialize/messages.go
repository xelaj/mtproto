package serialize

// messages.go отвечает за имплементацию сериализации сообщений. зашифрованных и незашифрованных.

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"

	ige "github.com/xelaj/mtproto/aes_ige"
	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/utils"
)

// CommonMessage это сообщение (зашифрованое либо открытое) которыми общаются между собой клиент и сервер
type CommonMessage interface {
	GetMsg() []byte
	GetMsgID() int64
	GetSeqNo() int32
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

func (msg *EncryptedMessage) GetMsg() []byte {
	return msg.Msg
}

func (msg *EncryptedMessage) GetMsgID() int64 {
	return msg.MsgID
}

func (msg *EncryptedMessage) GetSeqNo() int32 {
	return msg.SeqNo
}

func (msg *EncryptedMessage) Serialize(client MessageInformator, requireToAck bool) ([]byte, error) {
	obj, err := serializePacket(client, msg.Msg, msg.MsgID, requireToAck)
	if err != nil {
		return nil, err
	}

	encryptedData, err := ige.Encrypt(obj, client.GetAuthKey())
	if err != nil {
		return nil, errors.Wrap(err, "encrypting")
	}

	buf := bytes.NewBuffer(nil)
	cw := tl.NewWriteCursor(buf)
	err = cw.PutRawBytes(utils.AuthKeyHash(client.GetAuthKey()))
	if err != nil {
		return nil, err
	}

	err = cw.PutRawBytes(ige.MessageKey(obj))
	if err != nil {
		return nil, err
	}

	err = cw.PutRawBytes(encryptedData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeEncryptedMessage(data, authKey []byte) (*EncryptedMessage, error) {
	msg := new(EncryptedMessage)

	cr := tl.NewReadCursor(bytes.NewBuffer(data))
	keyHash, err := cr.PopRawBytes(tl.LongLen)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(keyHash, utils.AuthKeyHash(authKey)) {
		return nil, errors.New("wrong encryption key")
	}

	msg.MsgKey, err = cr.PopRawBytes(int128Len) // msgKey это хэш от расшифрованного набора байт, последние 16 символов
	if err != nil {
		return nil, err
	}

	encryptedData, err := cr.PopRawBytes(len(data) - (tl.LongLen + int128Len))
	if err != nil {
		return nil, err
	}

	decrypted, err := ige.Decrypt(encryptedData, authKey, msg.MsgKey)
	if err != nil {
		return nil, errors.Wrap(err, "decrypting message")
	}

	cr = tl.NewReadCursor(bytes.NewBuffer(decrypted))
	msg.Salt, err = cr.PopLong()
	if err != nil {
		return nil, err
	}

	msg.SessionID, err = cr.PopLong()
	if err != nil {
		return nil, err
	}

	msg.MsgID, err = cr.PopLong()
	if err != nil {
		return nil, err
	}

	seqNo, err := cr.PopUint()
	if err != nil {
		return nil, err
	}

	msg.SeqNo = int32(seqNo)

	messageLen, err := cr.PopUint()
	if err != nil {
		return nil, err
	}

	if len(decrypted) < int(messageLen)-(tl.LongLen+tl.LongLen+tl.LongLen+tl.WordLen+tl.WordLen) {
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

	msg.Msg, err = cr.PopRawBytes(int(messageLen))
	if err != nil {
		return nil, err
	}

	return msg, nil
}

type UnencryptedMessage struct {
	Msg   []byte
	MsgID int64
}

func (msg *UnencryptedMessage) GetMsg() []byte {
	return msg.Msg
}

func (msg *UnencryptedMessage) GetMsgID() int64 {
	return msg.MsgID
}

func (msg *UnencryptedMessage) GetSeqNo() int32 {
	return 0
}

func (msg *UnencryptedMessage) Serialize(client MessageInformator) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	cw := tl.NewWriteCursor(buf)

	// authKeyHash, always 0 if unencrypted
	if err := cw.PutLong(0); err != nil {
		return nil, err
	}

	if err := cw.PutLong(msg.MsgID); err != nil {
		return nil, err
	}

	if err := cw.PutUint(uint32(len(msg.Msg))); err != nil {
		return nil, err
	}

	if err := cw.PutRawBytes(msg.Msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeUnencryptedMessage(data []byte) (*UnencryptedMessage, error) {
	msg := new(UnencryptedMessage)

	buf := bytes.NewBuffer(data)
	cr := tl.NewReadCursor(buf)

	_, err := cr.PopRawBytes(tl.LongLen) // authKeyHash, always 0 if unencrypted
	if err != nil {
		return nil, err
	}

	msg.MsgID, err = cr.PopLong()
	if err != nil {
		return nil, err
	}

	mod := msg.MsgID & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("Wrong bits of message_id: %#v", uint64(mod))
	}

	messageLen, err := cr.PopUint()
	if err != nil {
		return nil, err
	}

	if len(data)-(tl.LongLen+tl.LongLen+tl.WordLen) != int(messageLen) {
		fmt.Println(len(data), int(messageLen), int(messageLen+(tl.LongLen+tl.LongLen+tl.WordLen)))
		return nil, fmt.Errorf("message not equal defined size: have %v, want %v", len(data), messageLen)
	}

	msg.Msg = buf.Bytes()

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
	MakeRequest(req tl.Object, resp interface{}) error
}

func serializePacket(client MessageInformator, msg []byte, messageID int64, requireToAck bool) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	cw := tl.NewWriteCursor(buf)

	saltBytes := make([]byte, tl.LongLen)
	binary.LittleEndian.PutUint64(saltBytes, uint64(client.GetServerSalt()))
	if err := cw.PutRawBytes(saltBytes); err != nil {
		return nil, err
	}

	if err := cw.PutLong(client.GetSessionID()); err != nil {
		return nil, err
	}

	if err := cw.PutLong(messageID); err != nil {
		return nil, err
	}

	if requireToAck { // не спрашивай, как это работает
		// почему тут добавляется бит не ебу
		if err := cw.PutUint(uint32(client.GetLastSeqNo() | 1)); err != nil {
			return nil, err
		}
	} else {
		if err := cw.PutUint(uint32(client.GetLastSeqNo())); err != nil {
			return nil, err
		}
	}

	if err := cw.PutUint(uint32(len(msg))); err != nil {
		return nil, err
	}

	if err := cw.PutRawBytes(msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
