package mtproto

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/serialize"
)

func (m *MTProto) SaveSession() (err error) {
	m.encrypted = true
	s := new(Session)
	s.Key = m.authKey
	s.Hash = m.authKeyHash
	buf := make([]byte, serialize.LongLen)
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
		return nil, fmt.Errorf("reading file: %w", err)
	}

	file := new(tokenStorageFormat)
	err = json.Unmarshal(data, file)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	res := new(Session)

	res.Key, err = base64.StdEncoding.DecodeString(file.Key)
	if err != nil {
		return nil, fmt.Errorf("invalid binary data of 'key': %w", err)
	}
	res.Hash, err = base64.StdEncoding.DecodeString(file.Hash)
	if err != nil {
		return nil, fmt.Errorf("invalid binary data of 'hash': %w", err)
	}
	res.Salt, err = base64.StdEncoding.DecodeString(file.Salt)
	if err != nil {
		return nil, fmt.Errorf("invalid binary data of 'salt': %w", err)
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

	return ioutil.WriteFile(path, data, 0600)
}
