// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

// this is ALL helpful unoficial telegram api methods.

package telegram

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"

	"github.com/xelaj/mtproto/telegram/internal/calls"
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
		return nil, errors.New("can't retrieve info due to user is not invited in chat already")
	default:
		panic("impossible type: " + reflect.TypeOf(resolved).String() + ", can't process it")
	}
}

func (c *Client) GetPossibleAllUsersOfGroup(ch InputChannel) ([]User, error) {
	resp100, err := c.ChannelsGetParticipants(ch, ChannelParticipantsFilter(&ChannelParticipantsRecent{}), 100, 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "getting 0-100 recent users")
	}
	parts100 := resp100.(*ChannelsChannelParticipantsObj).Participants
	users100 := resp100.(*ChannelsChannelParticipantsObj).Users
	resp200, err := c.ChannelsGetParticipants(ch, ChannelParticipantsFilter(&ChannelParticipantsRecent{}), 100, 100, 0)
	if err != nil {
		return nil, errors.Wrap(err, "getting 100-200 recent users")
	}
	parts200 := resp200.(*ChannelsChannelParticipantsObj).Participants
	users200 := resp200.(*ChannelsChannelParticipantsObj).Users
	users := append(users100, users200...)

	idsStore := make(map[int]User)
	for _, participant := range append(parts100, parts200...) {
		uid := participant.GetUserID()
		var realUser User
		for _, user := range users {
			if _, ok := user.(*UserEmpty); ok {
				continue
			}
			v, ok := user.(*UserObj)
			if !ok {
				panic(reflect.TypeOf(user).String())
			}
			if int(v.ID) == uid {
				realUser = user
				break
			}
		}
		if realUser == nil {
			panic(fmt.Sprintf("user %v not found", uid))
		}

		idsStore[uid] = realUser
	}

	searchedUsers, err := getUsersOfChannelBySearching(c, ch, "")
	if err != nil {
		return nil, errors.Wrap(err, "searching")
	}

	for k, v := range searchedUsers {
		idsStore[k] = v
	}

	sortedIds := make([]int, 0, len(idsStore))
	for k := range idsStore {
		sortedIds = append(sortedIds, k)
	}
	sort.Ints(sortedIds)
	res := make([]User, len(sortedIds))
	for i, id := range sortedIds {
		res[i] = idsStore[id]
	}

	return res, nil
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

const symbols = "abcdefghijklmnopqrstuvwxyz0123456789"

func getParticipants(c *Client, ch InputChannel, lastQuery string) (map[int]struct{}, error) {
	idsStore := make(map[int]struct{})
	for _, symbol := range symbols {
		query := lastQuery + string([]rune{symbol})
		filter := ChannelParticipantsFilter(&ChannelParticipantsSearch{Q: query})

		// начинаем с 100-200, что бы проверить, может нам нужно дополнительный символ вставлять
		resp200, err := c.ChannelsGetParticipants(ch, filter, 100, 100, 0)
		if err != nil {
			return nil, errors.Wrap(err, "getting 100-200 users with query: '"+query+"'")
		}
		users200 := resp200.(*ChannelsChannelParticipantsObj).Participants
		if len(users200) >= 100 {
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

func getUsersOfChannelBySearching(c *Client, ch InputChannel, lastQuery string) (map[int]User, error) {
	idsStore := make(map[int]User)
	for _, symbol := range symbols {
		query := lastQuery + string([]rune{symbol})
		filter := ChannelParticipantsFilter(&ChannelParticipantsSearch{Q: query})

		// начинаем с 100-200, что бы проверить, может нам нужно дополнительный символ вставлять
		resp200, err := c.ChannelsGetParticipants(ch, filter, 100, 100, 0)
		if err != nil {
			return nil, errors.Wrap(err, "getting 100-200 users with query: '"+query+"'")
		}
		parts200 := resp200.(*ChannelsChannelParticipantsObj).Participants
		users200 := resp200.(*ChannelsChannelParticipantsObj).Users
		if len(parts200) >= 100 {
			deepParticipants, err := getUsersOfChannelBySearching(c, ch, query)
			if err != nil {
				return nil, errors.Wrapf(err, "query '%v'", query)
			}
			for k, v := range deepParticipants {
				idsStore[k] = v
			}
			continue
		}

		resp100, err := c.ChannelsGetParticipants(ch, filter, 0, 100, 0)
		if err != nil {
			return nil, errors.Wrap(err, "getting 0-100 users with query: '"+query+"'")
		}
		parts100 := resp100.(*ChannelsChannelParticipantsObj).Participants
		users100 := resp100.(*ChannelsChannelParticipantsObj).Users

		users := append(users100, users200...)
		for _, participant := range append(parts100, parts200...) {
			uid := participant.GetUserID()
			var realUser User
			for _, user := range users {
				if _, ok := user.(*UserEmpty); ok {
					continue
				}
				v, ok := user.(*UserObj)
				if !ok {
					panic(reflect.TypeOf(user).String())
				}
				if int(v.ID) == uid {
					realUser = user
					break
				}
			}
			if realUser == nil {
				panic(fmt.Sprintf("user %v not found", uid))
			}

			idsStore[uid] = realUser
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
		check(err)
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

func (c *Client) PhoneGetCallConfigFormatted() (*calls.CallConfig, error) {
	jsonString, err := c.PhoneGetCallConfig()
	if err != nil {
		return nil, errors.Wrap(err, "calling phone.getCallConfig method")
	}

	data := &calls.CallConfig{}
	err = json.Unmarshal([]byte(jsonString.Data), &data)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}
	return data, nil
}
