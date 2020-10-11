package mtproto

import (
	"context"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/serialize"
	"github.com/xelaj/mtproto/utils"
)

func (m *MTProto) sendPacketNew(request serialize.TL, expectVector reflect.Type) (chan serialize.TL, error) {
	resp := make(chan serialize.TL)
	if m.serviceModeActivated {
		resp = m.serviceChannel
	}
	var data []byte
	var msgID = utils.GenerateMessageId()

	// может мы ожидаем вектор, см. erialize.RpcResult для понимания
	if expectVector != nil {
		m.msgsIdDecodeAsVector[msgID] = expectVector
	}

	if m.encrypted {
		requireToAck := false
		if MessageRequireToAck(request) {
			m.mutex.Lock()
			m.waitAck(msgID)
			m.mutex.Unlock()
			requireToAck = true
		}

		data = (&serialize.EncryptedMessage{
			Msg:         request.Encode(),
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
			Msg:   request.Encode(),
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
	//	dry.PanicIf(n != 3, fmt.Sprintf("expected read 3 bytes, got %d", n))
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
		return nil, fmt.Errorf("size is not length of int32, expected 4 bytes, got %d", n)
	}

	size := binary.LittleEndian.Uint32(sizeInBytes)
	// читаем сами данные
	data = make([]byte, int(size))
	n, err = reader.Read(data)
	dry.PanicIfErr(err)
	dry.PanicIf(n != int(size), fmt.Sprintf("expected read %d bytes, got %d", size, n))

	return data, nil
}
