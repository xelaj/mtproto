// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package session

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/xelaj/tl"
)

type genericFileSessionLoader struct {
	path string
}

var _ SessionLoader = (*genericFileSessionLoader)(nil)

func NewFromFile(path string) SessionLoader {
	return &genericFileSessionLoader{path: path}
}

func (l *genericFileSessionLoader) Load() (Session, error) {
	switch _, err := os.Stat(l.path); {
	case err == nil:
	case errors.Is(err, syscall.ENOENT):
		return Session{}, ErrSessionNotFound
	default:
		return Session{}, err
	}

	data, err := os.ReadFile(l.path)
	if err != nil {
		return Session{}, errors.Wrap(err, "reading file")
	}

	file := new(TokenStorageFormat)
	err = json.Unmarshal(data, file)
	if err != nil {
		return Session{}, errors.Wrap(err, "parsing file")
	}

	s, err := file.ReadSession()
	if err != nil {
		return Session{}, err
	}

	return *s, nil
}

func (l *genericFileSessionLoader) Store(s Session) error {
	dir, _ := filepath.Split(l.path)
	if info, err := os.Stat(dir); err != nil {
		return fmt.Errorf("%v: directory not found", dir)
	} else if !info.IsDir() {
		return fmt.Errorf("%v: not a directory", dir)
	}

	file := new(TokenStorageFormat)
	file.writeSession(&s)
	data, _ := json.Marshal(file)

	return os.WriteFile(l.path, data, 0600)
}

type TokenStorageFormat struct {
	Key  string `json:"key"`
	Salt string `json:"salt"`
}

func (t *TokenStorageFormat) writeSession(s *Session) {
	t.Key = base64.StdEncoding.EncodeToString(s.Key[:])
	t.Salt = encodeInt64ToBase64(s.Salt)
}

func (t *TokenStorageFormat) ReadSession() (*Session, error) {
	s := new(Session)
	var err error

	keyRaw, err := base64.StdEncoding.DecodeString(t.Key)
	if err != nil {
		return nil, errors.Wrap(err, "invalid binary data of 'key'")
	}
	copy(s.Key[:], keyRaw)

	s.Salt, err = decodeInt64ToBase64(t.Salt)
	if err != nil {
		return nil, errors.Wrap(err, "invalid binary data of 'salt'")
	}
	return s, nil
}

func encodeInt64ToBase64(i uint64) string {
	buf := make([]byte, tl.LongLen)
	binary.LittleEndian.PutUint64(buf, uint64(i))
	return base64.StdEncoding.EncodeToString(buf)
}

func decodeInt64ToBase64(i string) (uint64, error) {
	buf, err := base64.StdEncoding.DecodeString(i)
	if err != nil {
		return 0, err
	} else if len(buf) < 8 {
		return 0, errors.New("value is too short")
	}

	return binary.LittleEndian.Uint64(buf), nil
}
