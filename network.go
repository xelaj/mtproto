// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/xelaj/tl"

	"github.com/xelaj/mtproto/v2/internal/objects"
	"github.com/xelaj/mtproto/v2/internal/payload"
)

func (m *MTProto) makeRequest(ctx context.Context, msg []byte) (resp []byte, err error) {
	for resp == nil {
		resp, err = m.sendPacket(ctx, msg, true)
		if errors.Is(err, errRetryRequest) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("sending message: %w", err)
		}
	}

	return resp, nil
}

func (m *MTProto) sendPacket(ctx context.Context, msg []byte, expectAnswer bool) ([]byte, error) {
	id, err := m.transport.WriteMsg(ctx, msg, payload.InitiatorClient)
	if err != nil {
		return nil, err
	}

	if expectAnswer {
		ch := m.expectAnswer(id)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case r := <-ch:
			return r.data, r.err
		}
	}

	return nil, nil
}

func (m *MTProto) expectAnswer(msgID payload.MsgID) <-chan responseChanMsg {
	resp := make(chan responseChanMsg, 1)
	m.chanMux.Lock()
	defer m.chanMux.Unlock()
	m.responseChannels[msgID] = resp

	return resp
}

func (m *MTProto) writeRPCResponse(msgID payload.MsgID, data []byte) bool {
	m.chanMux.Lock()
	v, ok := m.responseChannels[msgID]
	if ok {
		delete(m.responseChannels, msgID)
	}
	m.chanMux.Unlock()

	if ok {
		if len(data) > tl.WordLen && binary.LittleEndian.Uint32(data) == objects.CrcRpcError {
			var e objects.RpcError
			if err := objects.Unmarshal(data, &e); err != nil {
				e = objects.RpcError{
					ErrorCode:    -1,
					ErrorMessage: fmt.Sprintf("can't unmarshal error: %v, data: %v", err, data),
				}
			}

			v <- responseChanMsg{err: rpcErrorToNative(e)}
		} else {
			v <- responseChanMsg{data: data}
		}

		close(v)
	}

	return ok
}

func (m *MTProto) rejectAllRequests(err error) {
	m.chanMux.Lock()
	defer m.chanMux.Unlock()

	for k, v := range m.responseChannels {
		delete(m.responseChannels, k)
		v <- responseChanMsg{err: err}
		close(v)
	}
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
