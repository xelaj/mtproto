// Copyright (c) 2020-2022 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package session

// SessionLoader is the interface which allows you to access sessions from
// different storages (like filesystem, database, s3 storage, etc.)
type SessionLoader interface {
	Load() (*Session, error)
	Store(*Session) error
}

// Session is a basic data of specific session. Typically, session stores
// default hostname of mtproto server (cause all accounts ties to specific
// server after sign in), session key, server hash and salt.
type Session struct {
	Key      []byte
	Hash     []byte
	Salt     int64
	Hostname string
}
