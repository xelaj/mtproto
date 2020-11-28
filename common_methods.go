// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

// методы (ну или функции, you name it), которые взяты отсюда https://core.telegram.org/schema/mtproto
// по сути это такие "алиасы" методов описанных объектов (которые в internal/mtproto/objects), идея взята
// из github.com/xelaj/vk

import (
	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/internal/mtproto/objects"
)

func (m *MTProto) reqPQ(nonce *tl.Int128) (*objects.ResPQ, error) {
	return objects.ReqPQ(m, nonce)
}

func (m *MTProto) reqDHParams(nonce, serverNonce *tl.Int128, p, q []byte, publicKeyFingerprint int64, encryptedData []byte) (objects.ServerDHParams, error) {
	return objects.ReqDHParams(m, nonce, serverNonce, p, q, publicKeyFingerprint, encryptedData)
}

func (m *MTProto) setClientDHParams(nonce, serverNonce *tl.Int128, encryptedData []byte) (objects.SetClientDHParamsAnswer, error) {
	return objects.SetClientDHParams(m, nonce, serverNonce, encryptedData)
}

func (m *MTProto) ping(pingID int64) (*objects.Pong, error) {
	return objects.Ping(m, pingID)
}
