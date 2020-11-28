// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto/encoding/tl"
)

func (m *MTProto) SaveSession() (err error) {
	m.encrypted = true
	s := new(Session)
	s.Key = m.authKey
	s.Hash = m.authKeyHash
	buf := make([]byte, tl.LongLen)
	binary.LittleEndian.PutUint64(buf, uint64(m.serverSalt))
	s.Salt = buf
	s.Hostname = m.addr
	err = SaveSession(s, m.tokensStorage)
	dry.PanicIfErr(err)

	return nil
}

func (m *MTProto) LoadSession() (err error) {
	s, err := LoadSession(m.tokensStorage)
	if errs.IsNotFound(err) {
		return err
	}
	dry.PanicIfErr(err)

	m.authKey = s.Key
	m.authKeyHash = s.Hash
	m.serverSalt = int64(binary.LittleEndian.Uint64(s.Salt)) // СОЛЬ ЭТО LONG
	m.addr = s.Hostname

	return nil
}

type tokenStorageFormat struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
	Hostname string `json:"hostname"`
}

type Session struct {
	Key      []byte
	Hash     []byte
	Salt     []byte
	Hostname string
}

func LoadSession(path string) (*Session, error) {
	if !dry.FileExists(path) {
		return nil, errs.NotFound("file", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}

	file := new(tokenStorageFormat)
	err = json.Unmarshal(data, file)
	if err != nil {
		return nil, errors.Wrap(err, "parsing file")
	}

	res := new(Session)

	res.Key, err = base64.StdEncoding.DecodeString(file.Key)
	if err != nil {
		return nil, errors.Wrap(err, "invalid binary data of 'key'")
	}
	res.Hash, err = base64.StdEncoding.DecodeString(file.Hash)
	if err != nil {
		return nil, errors.Wrap(err, "invalid binary data of 'hash'")
	}
	res.Salt, err = base64.StdEncoding.DecodeString(file.Salt)
	if err != nil {
		return nil, errors.Wrap(err, "invalid binary data of 'salt'")
	}
	res.Hostname = file.Hostname

	return res, nil
}

func SaveSession(s *Session, path string) error {
	file := new(tokenStorageFormat)
	file.Key = base64.StdEncoding.EncodeToString(s.Key)
	file.Hash = base64.StdEncoding.EncodeToString(s.Hash)
	file.Salt = base64.StdEncoding.EncodeToString(s.Salt)
	file.Hostname = s.Hostname

	data, _ := json.Marshal(file)

	dir, _ := filepath.Split(path)
	if !dry.FileExists(dir) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return errors.Wrap(err, "creating directory")
		}
	}
	if !dry.FileIsDir(dir) {
		return errors.New(path + ": not a directory")
	}

	return ioutil.WriteFile(path, data, 0600)
}
