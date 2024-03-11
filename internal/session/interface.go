package session

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
)

// SessionLoader is the interface which allows you to access sessions from
// different storages (like filesystem, database, s3 storage, etc.)
type SessionLoader interface {
	// must return ErrSessionNotFound, if session is not found
	Load() (Session, error)
	Store(Session) error
}

var ErrSessionNotFound = errors.New("session not found")

// Sesion is a basic data of specific session. Typically, session stores default
// hostname of mtproto server (cause all accounts ties to specific server after
// sign in), session key, server hash and salt.
type Session struct {
	Key  [256]byte // Authorization key
	Salt uint64    // соль сессии
}

func (s Session) KeyID() (keyID uint64) { return KeyID(s.Key) }

func KeyID(key [256]byte) (keyID uint64) {
	hash := sha1.Sum(key[:])

	return binary.LittleEndian.Uint64(hash[12:20])
}

// Defined here:
// https://core.telegram.org/mtproto/auth_key#9-server-responds-in-one-of-three-ways
//
// According to documentation, auth_key_aux_hash is a first 64 bits of key hash
// (or 8 bytes), and needs for retry request of diffie-hellman key exchange. But
// in real life, there is no use case for this part of hash.
//
// DEPRECATED: Even that it's mentioned in docs, it has no use case.
func (s Session) AuxHash() (keyID uint64) {
	hash := sha1.Sum(s.Key[:])

	return binary.LittleEndian.Uint64(hash[0:8])
}

type loader[T any] interface {
	Load() (T, error)
	Store(T) error
}

func IsSet(l SessionLoader) (bool, error) {
	if _, err := l.Load(); err == nil || errors.Is(err, ErrSessionNotFound) {
		return err == nil, nil
	} else {
		return false, err
	}
}

func GetKey(l SessionLoader) ([]byte, error) {
	return get(l, func(s Session) []byte { return s.Key[:] })
}

func GetSalt(l SessionLoader) (uint64, error) {
	return get(l, func(s Session) uint64 { return s.Salt })
}

func get[T, K any](l loader[T], f func(T) K) (res K, err error) {
	if v, err := l.Load(); err != nil {
		return res, err
	} else {
		return f(v), nil
	}
}
