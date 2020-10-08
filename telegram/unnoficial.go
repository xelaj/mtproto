// this is ALL helpful unoficial telegram api methods.

package telegram

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func (c *Client) GetChannelInfoByInviteLink(hashOrLink string) (*ChannelFull, error) {
	chat, err := c.GetChatInfoByHashLink(hashOrLink)
	if err != nil {
		return nil, err
	}

	channelSimpleData, ok := chat.(*Channel)
	if !ok {
		return nil, errors.New("not a channel")
	}
	id := channelSimpleData.Id
	hash := channelSimpleData.AccessHash

	data, err := c.ChannelsGetFullChannel(&ChannelsGetFullChannelParams{
		Channel: InputChannel(&InputChannelObj{
			ChannelId:  id,
			AccessHash: hash,
		}),
	})
	if err != nil {
		return nil, errors.Wrap(err, "retrieving full channel info")
	}
	fullChannel, ok := data.FullChat.(*ChannelFull)
	if !ok {
		return nil, errors.New("response not a ChannelFull, got '" + reflect.TypeOf(data.FullChat).String() + "'")
	}

	return fullChannel, nil
}

func (c *Client) GetChatInfoByHashLink(hashOrLink string) (Chat, error) {
	hash := hashOrLink
	hash = strings.TrimPrefix(hash, "http")
	hash = strings.TrimPrefix(hash, "s")
	hash = strings.TrimPrefix(hash, "://")
	hash = strings.TrimPrefix(hash, "t.me/")
	hash = strings.TrimPrefix(hash, "joinchat/")
	// checking now hash is HASH
	if !regexp.MustCompile(`^[a-zA-Z0-9+/=]+$`).MatchString(hash) {
		return nil, errors.New("'" + hash + "': not base64 hash")
	}

	resolved, err := c.MessagesCheckChatInvite(&MessagesCheckChatInviteParams{Hash: hash})
	if err != nil {
		return nil, errors.Wrap(err, "retrieving data by invite link")
	}

	switch res := resolved.(type) {
	case *ChatInviteAlready:
		return res.Chat, nil
	case *ChatInviteObj:
		return nil, errors.New("can't retrieve info due to user is not invited in chat  already")
	default:
		panic("impossible type: " + reflect.TypeOf(resolved).String() + ", can't process it")
	}
}
