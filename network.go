// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"github.com/xelaj/errs"

	"github.com/xelaj/mtproto/internal/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/messages"
	"github.com/xelaj/mtproto/internal/mtproto/objects"
	"github.com/xelaj/mtproto/internal/utils"
)

func (m *MTProto) sendPacket(request tl.Object, expectedTypes ...reflect.Type) (chan tl.Object, error) {
	msg, err := tl.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "encoding request message")
	}

	var (
		data  messages.Common
		msgID = utils.GenerateMessageId()
	)

	// adding types for parser if required
	if len(expectedTypes) > 0 {
		m.expectedTypes.Add(int(msgID), expectedTypes)
	}

	// dealing with response channel
	resp := m.getRespChannel()
	if isNullableResponse(request) {
		go func() { resp <- &objects.Null{} }() // goroutine cuz we don't read from it RIGHT NOW
	} else {
		m.responseChannels.Add(int(msgID), resp)
	}

	if m.encrypted {
		data = &messages.Encrypted{
			Msg:         msg,
			MsgID:       msgID,
			AuthKeyHash: m.authKeyHash,
		}
	} else {
		data = &messages.Unencrypted{ //nolint: errcheck нешифрованое не отправляет ошибки
			Msg:   msg,
			MsgID: msgID,
		}
	}

	// must write synchroniously, cuz seqno must be upper each request
	m.seqNoMutex.Lock()
	defer m.seqNoMutex.Unlock()

	err = m.transport.WriteMsg(data, MessageRequireToAck(request))
	if err != nil {
		return nil, errors.Wrap(err, "sending request")
	}

	if m.encrypted {
		// since we sending this message, we are incrementing the seqno BUT ONLY when we
		// are sending an encrypted message. why? I don’t know. But the fact remains:
		// we must to block seqno, cause messages with a bigger seqno can go faster than
		// messages with a smaller one.
		m.seqNo += 2
	}

	return resp, nil
}

func (m *MTProto) writeRPCResponse(msgID int, data tl.Object) error {
	v, ok := m.responseChannels.Get(msgID)
	if !ok {
		return errs.NotFound("msgID", strconv.Itoa(msgID))
	}

	v <- data

	m.responseChannels.Delete(msgID)
	m.expectedTypes.Delete(msgID)
	return nil
}

func (m *MTProto) getRespChannel() chan tl.Object {
	if m.serviceModeActivated {
		return m.serviceChannel
	}
	return make(chan tl.Object)
}

// проверяет, надо ли ждать от сервера пинга
func isNullableResponse(t tl.Object) bool {
	switch t.(type) {
	case /**objects.Ping,*/ *objects.Pong, *objects.MsgsAck:
		return true
	default:
		return false
	}
}
