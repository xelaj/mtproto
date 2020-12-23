// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

// this is ALL helpful unoficial telegram api methods.

package telegram

import (
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/xelaj/go-dry"
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
	id := channelSimpleData.ID
	hash := channelSimpleData.AccessHash

	data, err := c.ChannelsGetFullChannel(&InputChannelObj{
		ChannelID:  id,
		AccessHash: hash,
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

	resolved, err := c.MessagesCheckChatInvite(hash)
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

func (c *Client) GetPossibleAllParticipantsOfGroup(ch InputChannel) ([]int, error) {
	resp100, err := c.ChannelsGetParticipants(ch, ChannelParticipantsFilter(&ChannelParticipantsRecent{}), 100, 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "getting 0-100 recent users")
	}
	users100 := resp100.(*ChannelsChannelParticipantsObj).Participants
	resp200, err := c.ChannelsGetParticipants(ch, ChannelParticipantsFilter(&ChannelParticipantsRecent{}), 100, 100, 0)
	if err != nil {
		return nil, errors.Wrap(err, "getting 100-200 recent users")
	}
	users200 := resp200.(*ChannelsChannelParticipantsObj).Participants

	idsStore := make(map[int]struct{})
	for _, participant := range append(users100, users200...) {
		switch user := participant.(type) {
		case *ChannelParticipantObj:
			idsStore[int(user.UserID)] = struct{}{}
		case *ChannelParticipantAdmin:
			idsStore[int(user.UserID)] = struct{}{}
		case *ChannelParticipantCreator:
			idsStore[int(user.UserID)] = struct{}{}
		default:
			pp.Println(user)
			panic("что?")
		}
	}

	searchedUsers, err := getParticipants(c, ch, "")
	if err != nil {
		return nil, errors.Wrap(err, "searching")
	}

	for k, v := range searchedUsers {
		idsStore[k] = v
	}

	res := make([]int, 0, len(idsStore))
	for k := range idsStore {
		res = append(res, k)
	}

	sort.Ints(res)

	return res, nil
}

var symbols = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
	"l", "m", "n", "o", "p", "q", "r", "s", "t", "u",
	"v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func getParticipants(c *Client, ch InputChannel, lastQuery string) (map[int]struct{}, error) {
	idsStore := make(map[int]struct{})
	for _, symbol := range symbols {
		query := lastQuery + symbol
		filter := ChannelParticipantsFilter(&ChannelParticipantsSearch{Q: query})

		// начинаем с 100-200, что бы проверить, может нам нужно дополнительный символ вставлять
		resp200, err := c.ChannelsGetParticipants(ch, filter, 100, 100, 0)
		if err != nil {
			return nil, errors.Wrap(err, "getting 100-200 users with query: '"+query+"'")
		}
		users200 := resp200.(*ChannelsChannelParticipantsObj).Participants
		if len(users200) >= 200 {
			deepParticipants, err := getParticipants(c, ch, query)
			if err != nil {
				return nil, err
			}
			for k := range deepParticipants {
				idsStore[k] = struct{}{}
			}
			continue
		}

		resp100, err := c.ChannelsGetParticipants(ch, filter, 0, 100, 0)
		if err != nil {
			return nil, errors.Wrap(err, "getting 0-100 users with query: '"+query+"'")
		}
		users100 := resp100.(*ChannelsChannelParticipantsObj).Participants

		for _, participant := range append(users100, users200...) {
			switch user := participant.(type) {
			case *ChannelParticipantObj:
				idsStore[int(user.UserID)] = struct{}{}
			case *ChannelParticipantAdmin:
				idsStore[int(user.UserID)] = struct{}{}
			case *ChannelParticipantCreator:
				idsStore[int(user.UserID)] = struct{}{}
			default:
				pp.Println(user)
				panic("что?")
			}
		}
	}

	return idsStore, nil
}

// GetChatByID is searching in all user chats specific chat with input id
// TODO: need to test
func (c *Client) GetChatByID(chatID int) (Chat, error) {
	resp, err := c.MessagesGetAllChats([]int32{})
	if err != nil {
		return nil, errors.Wrap(err, "getting all chats")
	}
	chats := resp.(*MessagesChatsObj)
	for _, chat := range chats.Chats {
		switch c := chat.(type) {
		case *ChatObj:
			if int(c.ID) == chatID {
				return c, nil
			}
		case *Channel:
			if -1*(int(c.ID)+(1000000000000)) == chatID { // -100<channelID, specific for bots>
				return c, nil
			}
		default:
			pp.Println(c)
			panic("???")
		}
	}

	return nil, errs.NotFound("chatID", strconv.Itoa(chatID))
}

// returning all user ids in specific SUPERGROUP. Note that, SUPERGROUP IS NOT CHANNEL! Major difference in how
// users list returning: in supergroup you aren't limited in offset of fetching users. But channel is
// different: telegram forcely limit you in up to 200 users per single request (you can sort it by recently
// joined, search query, etc.)
func (c *Client) AllUsersInChat(chatID int) ([]int, error) {
	chat, err := c.GetChatByID(chatID)
	if err != nil {
		return nil, errors.Wrap(err, "getting chat by id: "+strconv.Itoa(chatID))
	}

	channel, ok := chat.(*Channel)
	if !ok {
		return nil, errors.New("Not a channel")
	}

	inCh := InputChannel(&InputChannelObj{
		ChannelID:  channel.ID,
		AccessHash: channel.AccessHash,
	})

	res := make(map[int]struct{})
	totalCount := 100 // at least 100
	offset := 0
	for offset < totalCount {
		resp, err := c.ChannelsGetParticipants(
			inCh,
			ChannelParticipantsFilter(&ChannelParticipantsRecent{}),
			100,
			int32(offset),
			0,
		)
		dry.PanicIfErr(err)
		data := resp.(*ChannelsChannelParticipantsObj)
		totalCount = int(data.Count)
		for _, participant := range data.Participants {
			switch user := participant.(type) {
			// здесь хоть и параметр userId одинаковый, да вот объекты разные...
			case *ChannelParticipantSelf:
				res[int(user.UserID)] = struct{}{}
			case *ChannelParticipantObj:
				res[int(user.UserID)] = struct{}{}
			case *ChannelParticipantAdmin:
				res[int(user.UserID)] = struct{}{}
			case *ChannelParticipantCreator:
				res[int(user.UserID)] = struct{}{}
			default:
				pp.Println(user)
				return nil, errors.New("found too specific object")
			}
		}

		offset += 100
		pp.Println(offset, totalCount)
	}

	total := make([]int, 0, len(res))
	for k := range res {
		total = append(total, k)
	}

	sort.Ints(total)

	return total, nil
}

// returning all user ids in specific SUPERGROUP. Note that, SUPERGROUP IS NOT CHANNEL! Major difference in how
// users list returning: in supergroup you aren't limited in offset of fetching users. But channel is
// different: telegram forcely limit you in up to 200 users per single request (you can sort it by recently
// joined, search query, etc.)
//
// This method is running too long for simple call, if channel is big, so call it inside goroutine with
// callback.
func (c *Client) AllUsersInChannel(channelID int) ([]int, error) {
	chat, err := c.GetChatByID(channelID)
	if err != nil {
		return nil, errors.Wrap(err, "getting chat by id: "+strconv.Itoa(channelID))
	}

	channel, ok := chat.(*Channel)
	if !ok {
		return nil, errors.New("Not a channel")
	}

	inCh := InputChannel(&InputChannelObj{
		ChannelID:  channel.ID,
		AccessHash: channel.AccessHash,
	})

	ids, err := c.GetPossibleAllParticipantsOfGroup(inCh)
	if err != nil {
		return nil, errors.Wrap(err, "getting clients of chat")
	}
	return ids, nil
}
