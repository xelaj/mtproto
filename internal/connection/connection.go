// Copyright (c) 2020-2022 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package connection

import (
	"github.com/xelaj/mtproto/internal/mode"
	"github.com/xelaj/mtproto/internal/mtproto/messages"
	"github.com/xelaj/mtproto/internal/transport"
)

type Connection interface {
	Close() error
	WriteMsg(msg messages.Common) error
	ReadMsg() (messages.Common, error)
}

type conn struct {
	t transport.Transport
	m mode.Variant
}

func New(t transport.Transport, m mode.Variant) Connection {
	return &conn{t, m}
}
