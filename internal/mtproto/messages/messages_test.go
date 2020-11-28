// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package messages_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/go-dry"

	. "github.com/xelaj/mtproto/internal/mtproto/messages"
)

type DummyClient struct {
	sessionID  int64
	lastSeqNo  int32
	serverSalt int64
	authKey    []byte
}

func (d *DummyClient) GetSessionID() int64 {
	return d.sessionID
}

func (d *DummyClient) GetLastSeqNo() int32 {
	return d.lastSeqNo
}

func (d *DummyClient) GetServerSalt() int64 {
	return d.serverSalt
}

func (d *DummyClient) GetAuthKey() []byte {
	return d.authKey
}

var client = &DummyClient{
	authKey: Hexed("28F43A9E1F5B15C093445BDBA697C78DCE12B53C8F05AE86F1E25338DC8EF962" +
		"E9B89C8E560955FFA0E1A45C8D121A9AEFDB89C88BB1493374959C6D6E5C46D1" +
		"42775B2A56C889D184ECB0B570E5BC763AAC504A35F6A5B259D1C20A01671C47" +
		"24185463D7DBC8F7376743F32AFEEC4C272C814DC10612AE8A12C861C4BFA04B" +
		"2CFC96880C9AFCDDD3584465F93C6D2597E433D39777BDF9C4C613D7F43D7F65" +
		"F59E3462137BD4C009049D154A73048679C09D832A41A12A1F646455B5BD6263" +
		"02AAF9798BC8A97A219CF9FF22FB3362943FD67E460258295D0984BD3FBA15A0" +
		"D6BDF1F48F51CA65BD6C1CDD9C0509A73EB320379118BC586F7564F391DA1490"),
}

func TestSerializeUnencryptedMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *Unencrypted
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "usualMessage",
			msg: &Unencrypted{
				Msg:   []byte("hello mtproto messages!"),
				MsgID: 123,
			}, //         |   authKeyHash||    msgID     ||mlen ||rest of data                 >>
			want: Hexed("00000000000000007b000000000000001700000068656c6c6f206d7470726f746f206d65" +
				"73736167657321"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := tt.wantErr
			if wantErr == nil {
				wantErr = assert.NoError
			}

			got, err := tt.msg.Serialize(client)
			if !wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSerializeEncryptedMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *Encrypted
		ack     bool
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "usualMessage",
			msg: &Encrypted{
				Msg:   []byte("hello mtproto messages!"),
				MsgID: 123,
			},
			want: Hexed("26C877F943462A4247DC1ACF8232053834D146BE164547066924AB8509629E8C" +
				"2C2B353A77C8A37EAB2D8982723DD7027941408F91F84BF1FE8FD7CDE3E4D29F" +
				"AFAFFD26489DFFC18DFC09C9C4A53973B9C943910C28B687"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := tt.wantErr
			if wantErr == nil {
				wantErr = assert.NoError
			}

			got, err := tt.msg.Serialize(client, tt.ack)

			if !wantErr(t, err) {
				return
			}
			if !assert.Equal(t, tt.want, got) {
				return
			}

			// TODO: тоже не работает
			//pp.Println(DeserializeEncryptedMessage(got, client.GetAuthKey()))
		})
	}
}

func Hexed(in string) []byte {
	res, err := hex.DecodeString(in)
	dry.PanicIfErr(err)
	return res
}
