package tl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/telegram"
)

func TestEquality(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
		fill interface{}
	}{
		{
			name: "MessagesChatsObj",
			obj: &telegram.MessagesChatsObj{
				Chats: []telegram.Chat{
					&telegram.ChatObj{
						Creator:           false,
						Kicked:            false,
						Left:              false,
						Deactivated:       true,
						Id:                123,
						Title:             "abcdef",
						Photo:             &telegram.ChatPhotoEmpty{},
						ParticipantsCount: 123,
						Date:              1,
						Version:           1,
						MigratedTo:        &telegram.InputChannelEmpty{},
						AdminRights: &telegram.ChatAdminRights{
							ChangeInfo: true,
							BanUsers:   true,
						},
						DefaultBannedRights: &telegram.ChatBannedRights{
							SendGames: true,
						},
					},
				},
			},
			fill: &telegram.MessagesChatsObj{},
		},
	}

	for _, tt := range tests {
		encoded, err := tl.Encode(tt.obj)
		assert.NoError(t, err)

		err = tl.Decode(encoded, tt.fill)
		assert.NoError(t, err)

		assert.Equal(t, tt.obj, tt.fill)
	}
}
