package tl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/mtproto/encoding/tl"
)

type MultipleChats struct {
	Chats []any
}

func (*MultipleChats) CRC() uint32 {
	return uint32(0xff1144cc)
}

type Chat struct {
	Creator           bool `tl:"flag:0,encoded_in_bitflags"`
	Kicked            bool `tl:"flag:1,encoded_in_bitflags"`
	Left              bool `tl:"flag:2,encoded_in_bitflags"`
	Deactivated       bool `tl:"flag:5,encoded_in_bitflags"`
	ID                int32
	Title             string
	Photo             string
	ParticipantsCount int32
	Date              int32
	Version           int32
	AdminRights       *Rights `tl:"flag:14"`
	BannedRights      *Rights `tl:"flag:18"`
}

func (*Chat) CRC() uint32 {
	return uint32(0x3bda1bde)
}

func (*Chat) FlagIndex() int {
	return 0
}

type Rights struct {
	DeleteMessages bool `tl:"flag:3,encoded_in_bitflags"`
	BanUsers       bool `tl:"flag:4,encoded_in_bitflags"`
}

func (*Rights) CRC() uint32 {
	return uint32(0x5fb224d5)
}

func (*Rights) FlagIndex() int {
	return 0
}

func TestEquality(t *testing.T) {
	tests := []struct {
		name string
		obj  any
		fill any
	}{
		{
			name: "MessagesChatsObj",
			obj: &MultipleChats{
				Chats: []any{
					&Chat{
						Creator:           false,
						Kicked:            false,
						Left:              false,
						Deactivated:       true,
						ID:                123,
						Title:             "abcdef",
						Photo:             "pikcha.png",
						ParticipantsCount: 123,
						Date:              1,
						Version:           1,
						AdminRights: &Rights{
							DeleteMessages: true,
							BanUsers:       true,
						},
						BannedRights: &Rights{
							DeleteMessages: false,
							BanUsers:       false,
						},
					},
				},
			},
			fill: &MultipleChats{},
		},
	}

	for _, tt := range tests {
		encoded, err := tl.Marshal(tt.obj)
		assert.NoError(t, err)

		err = tl.Decode(encoded, tt.fill)
		assert.NoError(t, err)

		assert.Equal(t, tt.obj, tt.fill)
	}
}
