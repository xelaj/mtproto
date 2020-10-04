package mtproto

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/serialize"
	"github.com/xelaj/mtproto/utils"
)

func (m *MTProto) sendPacketNew(request serialize.TL) (chan serialize.TL, error) {
	resp := make(chan serialize.TL)
	if m.serviceModeActivated {
		resp = m.serviceChannel
	}
	var data []byte
	var msgID = utils.GenerateMessageId()
	if m.encrypted {
		requireToAck := false
		if MessageRequireToAck(request) {
			m.mutex.Lock()
			m.waitAck(msgID)
			m.mutex.Unlock()
			requireToAck = true
		}

		data = (&serialize.EncryptedMessage{
			Msg:         request,
			MsgID:       msgID,
			AuthKeyHash: m.authKeyHash,
		}).Serialize(m, requireToAck)

		if !isNullableResponse(request) {
			m.mutex.Lock()

			m.responseChannels[msgID] = resp
			m.mutex.Unlock()
		} else {
			// ответов на TL_Ack, TL_Pong и пр. не требуется
			go func() {
				// горутина, т.к. мы ПРЯМО СЕЙЧАС из resp не читаем
				resp <- &serialize.Null{}
			}()
		}
		// этот кусок не часть кодирования так что делаем при отправке
		m.lastSeqNo += 2
	} else {
		data = (&serialize.UnencryptedMessage{
			Msg:   request,
			MsgID: msgID,
		}).Serialize(m)
	}

	//? https://core.telegram.org/mtproto/mtproto-transports#intermediate
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	_, err := m.conn.Write(size)
	dry.PanicIfErr(err)

	//? https://core.telegram.org/mtproto/mtproto-transports#abridged
	// _, err := m.conn.Write(utils.PacketLengthMTProtoCompatible(data))
	// dry.PanicIfErr(err)
	println("writing message")
	_, err = m.conn.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "sending request")
	}

	return resp, nil
}

func (m *MTProto) writeRPCResponse(msgID int, data serialize.TL) error {
	m.mutex.Lock()
	v, ok := m.responseChannels[int64(msgID)]
	if !ok {
		return errs.NotFound("msgID", strconv.Itoa(msgID))
	}

	v <- data

	delete(m.responseChannels, int64(msgID))
	m.mutex.Unlock()
	return nil
}

func (m *MTProto) readFromConn(ctx context.Context) (data []byte, err error) {
	err = m.conn.SetReadDeadline(time.Now().Add(readTimeout)) // возможно поможет???
	dry.PanicIfErr(err)

	reader := dry.NewCancelableReader(ctx, m.conn)
	// https://core.telegram.org/mtproto/mtproto-transports#abridged
	// что делаем:
	// в conn есть определенный буффер, все что телега присылает, мы сохраняем в буффере, и потом через
	// Read читаем. т.к. маленькие пакеты (до 127 байт)  кодируют длину в 1 байт, а побольше в 4, то
	// мы читаем сначала 1 байт, смотрим, это 0xef или нет, если да, то читаем оставшиеся 3 байта и получаем длину
	//firstByte, err := reader.ReadByte()
	//dry.PanicIfErr(err)
	//
	//sizeInBytes, err := utils.GetPacketLengthMTProtoCompatible([]byte{firstByte})
	//if err == utils.ErrPacketSizeIsBigger {
	//	restOfSize := make([]byte, 3)
	//	n, err := reader.Read(restOfSize)
	//	dry.PanicIfErr(err)
	//	dry.PanicIf(n != 3, "expected read 3 bytes, got "+strconv.Itoa(n))
	//
	//	sizeInBytes, _ = utils.GetPacketLengthMTProtoCompatible(append([]byte{firstByte}, restOfSize...))
	//
	//	pp.Println(firstByte, restOfSize, sizeInBytes)
	//}

	// https://core.telegram.org/mtproto/mtproto-transports#intermediate
	sizeInBytes := make([]byte, 4)
	n, err := reader.Read(sizeInBytes)
	if err != nil {
		pp.Println(sizeInBytes, err)
		return nil, errors.Wrap(err, "reading length")
	}
	if n != 4 {
		return nil, errors.New("size is not length of int32, expected 4 bytes, got " + strconv.Itoa(n))
	}

	size := binary.LittleEndian.Uint32(sizeInBytes)
	// читаем сами данные
	data = make([]byte, int(size))
	n, err = reader.Read(data)
	dry.PanicIfErr(err)
	dry.PanicIf(n != int(size), "expected read "+strconv.Itoa(int(size))+" bytes, got "+strconv.Itoa(n))

	return data, nil
}

//! DEPRECATED
func (m *MTProto) sendPacket(msg serialize.TL, resp chan serialize.TL) error {
	var data []byte
	var msgID = utils.GenerateMessageId()
	if m.encrypted {
		requireToAck := false
		if MessageRequireToAck(msg) {
			m.mutex.Lock()
			//m.msgsIdToAck[msgID] = packetToSend{msg, resp}
			m.mutex.Unlock()
			requireToAck = true

		}

		data = (&serialize.EncryptedMessage{
			Msg:         msg,
			MsgID:       msgID,
			AuthKeyHash: m.authKeyHash,
		}).Serialize(m, requireToAck)

		if resp != nil {
			m.mutex.Lock()
			m.msgsIdToResp[msgID] = resp
			m.mutex.Unlock()
		}
		// этот кусок не часть кодирования так что делаем при отправке
		m.lastSeqNo += 2
	} else {
		data = (&serialize.UnencryptedMessage{
			Msg:   msg.(serialize.TL),
			MsgID: msgID,
		}).Serialize(m)
	}
	_, err := m.conn.Write(data)
	dry.PanicIfErr(err)
	return nil
}

//! DEPRECATED
func (m *MTProto) read(stop <-chan struct{}) (serialize.TL, error) {
	var err error
	var obj serialize.TL

	err = m.conn.SetReadDeadline(time.Now().Add(readTimeout))
	dry.PanicIfErr(err)

	// что делаем:
	// в conn есть определенный буффер, все что телега присылает, мы сохраняем в буффере, и потом через
	// Read читаем. т.к. маленькие пакеты (до 127 байт)  кодируют длину в 1 байт, а побольше в 4, то
	// мы читаем сначала 1 байт, смотрим, это 0xef или нет, если да, то читаем оставшиеся 3 байта и получаем длину
	firstByte := make([]byte, 1)
	println("start reading data")
	_, err = m.conn.Read(firstByte)
	println("data was read")
	if stop != nil {
		select {
		case <-stop:
			return nil, nil
		default:
		}
	}
	if err != nil {
		pp.Println(err)
		panic(err)
	}

	sizeInBytes, err := utils.GetPacketLengthMTProtoCompatible(firstByte)
	if err == utils.ErrPacketSizeIsBigger {
		restOfSize := make([]byte, 3)
		_, err = m.conn.Read(restOfSize)
		dry.PanicIfErr(err)
		sizeInBytes, _ = utils.GetPacketLengthMTProtoCompatible(append(firstByte, restOfSize...))
	}

	// читаем сами данные
	data := make([]byte, sizeInBytes)
	_, err = m.conn.Read(data)
	dry.PanicIfErr(err)

	// проверим, что это не код ошибки
	err = CatchResponseErrorCode(data)
	if err != nil {
		return nil, errors.Wrap(err, "Server response error")
	}

	if IsPacketEncrypted(data) {
		msg, err := serialize.DeserializeEncryptedMessage(data, m.GetAuthKey())
		dry.PanicIfErr(err)
		obj = msg.Msg
		m.seqNo = msg.SeqNo
		m.msgId = msg.MsgID
	} else {
		msg, err := serialize.DeserializeUnencryptedMessage(data)
		dry.PanicIfErr(err)
		obj = msg.Msg
		m.seqNo = 0
		m.msgId = msg.MsgID
	}

	mod := m.msgId & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("Wrong bits of message_id: %d", mod)
	}

	return obj, nil
}
