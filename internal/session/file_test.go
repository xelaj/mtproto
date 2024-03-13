// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package session_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/xelaj/mtproto/v2/internal/session"
)

const keyStr = "qgelGXIGnuXOWfx8ZfTV8RhcAhxInOxGiBbYA6H7edSVc49AYi1Y/rc+/dLA5" +
	"IvV4lVh1SUpaduCiXEiGT3MOSsR1DyiqVutzqgKcBXJOAPnLIu59q5SJA5VM2Z0NeLY/F3qo" +
	"+AQq+11olPaMZYeN44lHXIARxuZdVeAwXIXgCaRmlx3KAtVCjlGRhIlQLaOL+HDq1EfVMM7H" +
	"/esyxjZ6mu+Lqbuii7W1sf2ojYF7a1STrHjH4Hd+FLv8nQFAm/Xzl+znAFv3HtYpd9ETZqPJ" +
	"7OtATi5tiK70dd8qRY1Qwu4WndFyNL8HpLuQQ8GtnIS9PsVlN+Xzlk/YWyg/bh9aQ=="

var key = must(keyFromBase64(keyStr))

func TestMTProto_SaveSession(t *testing.T) {
	storePath := filepath.Join(os.TempDir(), "session.json")
	defer os.Remove(storePath)

	os.Remove(storePath)

	storage := NewFromFile(storePath)
	err := storage.Store(Session{
		Key:  key,
		Salt: 0,
	})
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(storePath)
	check(err)

	assert.Equal(t, `{"key":"`+keyStr+`","salt":"AAAAAAAAAAA="}`, string(data))
}

func TestMTProto_LoadSession(t *testing.T) {
	storePath := filepath.Join(os.TempDir(), "session.json")
	tmpData := `{"key":"` + keyStr + `","salt":"AAAAAAAAAAA="}`
	ioutil.WriteFile(storePath, []byte(tmpData), 0666)
	defer os.Remove(storePath)

	storage := NewFromFile(storePath)

	sess, err := storage.Load()
	require.NoError(t, err)

	assert.Equal(t, Session{
		Key:  key,
		Salt: 0,
	}, sess)
}

func keyFromBase64(s string) (key [256]byte, err error) {
	keyRaw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return key, fmt.Errorf("invalid binary data of 'key': %w", err)
	}
	copy(key[:], keyRaw)

	return key, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](t T, err error) T { check(err); return t }
