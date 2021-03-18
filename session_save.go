// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"github.com/xelaj/errs"

	"github.com/xelaj/mtproto/internal/session"
)

func (m *MTProto) SaveSession() (err error) {
	m.encrypted = true
	s := new(session.Session)
	s.Key = m.authKey
	s.Hash = m.authKeyHash
	s.Salt = m.serverSalt
	s.Hostname = m.addr
	err = session.SaveSession(s, m.tokensStorage)
	check(err)

	return nil
}

func (m *MTProto) LoadSession() (err error) {
	s, err := session.LoadSession(m.tokensStorage)
	if errs.IsNotFound(err) {
		return err
	}
	check(err)

	m.authKey = s.Key
	m.authKeyHash = s.Hash
	m.serverSalt = s.Salt
	m.addr = s.Hostname

	return nil
}
