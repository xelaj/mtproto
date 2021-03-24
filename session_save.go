// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"github.com/xelaj/mtproto/internal/session"
)

func (m *MTProto) SaveSession() (err error) {
	s := new(session.Session)
	s.Key = m.authKey
	s.Hash = m.authKeyHash
	s.Salt = m.serverSalt
	s.Hostname = m.addr
	return session.SaveSession(s, m.tokensStorage)
}

func (m *MTProto) LoadSession(s *session.Session) {
	m.authKey = s.Key
	m.authKeyHash = s.Hash
	m.serverSalt = s.Salt
	m.addr = s.Hostname
}
