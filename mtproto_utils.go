// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto/internal/encoding/tl"
	"github.com/xelaj/mtproto/internal/utils"
)

// мелкие методы, которые сделаны для понимания алгоритмов и кода вцелом

// waitAck добавляет в список id сообщения, которому нужно подтверждение
// возвращает true, если ранее этого id не было
func (m *MTProto) waitAck(msgID int64) bool {
	_, ok := m.idsToAck[msgID]
	m.idsToAck[msgID] = null{}
	return !ok
}

// gotAck удаляет элемент из списка id сообщений, на который ожидается ack.
// возвращается true, если id был найден
func (m *MTProto) gotAck(msgID int64) bool {
	m.idsToAckMutex.Lock()
	_, ok := m.idsToAck[msgID]
	delete(m.idsToAck, msgID)
	m.idsToAckMutex.Unlock()
	return ok
}

// resetAck сбрасывает целиком список сообщений, которым нужен ack
func (m *MTProto) resetAck() {
	m.idsToAck = make(map[int64]null)
}

// получает текущий идентификатор сессии
func (m *MTProto) GetSessionID() int64 {
	return m.sessionId
}

// Получает lastSeqNo
func (m *MTProto) GetLastSeqNo() int32 {
	return m.lastSeqNo
}

// получает соль
func (m *MTProto) GetServerSalt() int64 {
	return m.serverSalt
}

// получает ключ авторизации
func (m *MTProto) GetAuthKey() []byte {
	return m.authKey
}

func (m *MTProto) SetAuthKey(key []byte) {
	m.authKey = key
	m.authKeyHash = utils.AuthKeyHash(m.authKey)
}

func (m *MTProto) MakeRequest(msg tl.Object) (any, error) {
	return m.makeRequest(msg)
}

func (m *MTProto) MakeRequestWithHintToDecoder(msg tl.Object, expectedTypes ...reflect.Type) (any, error) {
	if len(expectedTypes) == 0 {
		return nil, errors.New("expected a few hints. If you don't need it, use m.MakeRequest")
	}
	return m.makeRequest(msg, expectedTypes...)
}

func (m *MTProto) recoverGoroutine() {
	if r := recover(); r != nil {
		if m.RecoverFunc != nil {
			fmt.Println(dry.StackTrace(0))
			m.RecoverFunc(r)
		} else {
			panic(r)
		}
	}
}

func (m *MTProto) AddCustomServerRequestHandler(handler customHandlerFunc) {
	m.serverRequestHandlers = append(m.serverRequestHandlers, handler)
}

func (m *MTProto) warnError(err error) {
	if err == nil {
		return
	}
	if m.Warnings != nil {
		m.Warnings <- err
	}
}
