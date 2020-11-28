package main

import (
	"sort"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"
)

func main() {
	println("firstly, you need to authorize. after exapmle 'auth', uo will signin")
	// edit these params for you!
	client, err := telegram.NewClient(telegram.ClientConfig{
		// where to store session configuration. must be set
		SessionFile: "/home/me/.local/var/lib/mtproto/session1.json",
		// host address of mtproto server. Actually, it can'be mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is patrh to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile: "/home/me/.local/var/lib/mtproto/tg_public_keys.pem",
		AppID:          94575,                              // app id, could be find at https://my.telegram.org
		AppHash:        "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
	})
	dry.PanicIfErr(err)

	// get this hash from channel invite link (after t.me/join/<HASH>)
	hash := "AAAAAEkCCtoerhjfii34iiii" // add here any link that you are ADMINISTRATING cause participants can be viewed only by admins

	// syntax sugared method, more easy to read than default ways to solve some troubles
	peer, err := client.GetChatInfoByHashLink(hash)
	dry.PanicIfErr(err)

	total, err := client.GetPossibleAllParticipantsOfGroup(&telegram.InputChannelObj{
		ChannelID:  peer.(*telegram.Channel).ID,
		AccessHash: peer.(*telegram.Channel).AccessHash,
	})

	dry.PanicIfErr(err)
	pp.Println(total, len(total))

	println("this is partial users in CHANNEL. In supergroup you can use more easy way to find, see below")

	resolved, err := client.ContactsResolveUsername("gogolang")
	dry.PanicIfErr(err)

	channel := resolved.Chats[0].(*telegram.Channel)
	inCh := telegram.InputChannel(&telegram.InputChannelObj{
		ChannelID:  channel.ID,
		AccessHash: channel.AccessHash,
	})

	res := make(map[int]struct{})
	totalCount := 100 // at least 100
	offset := 0
	for offset < totalCount {
		resp, err := client.ChannelsGetParticipants(inCh, telegram.ChannelParticipantsFilter(&telegram.ChannelParticipantsRecent{}), 100, int32(offset), 0)
		dry.PanicIfErr(err)
		data := resp.(*telegram.ChannelsChannelParticipantsObj)
		totalCount = int(data.Count)
		for _, participant := range data.Participants {
			switch user := participant.(type) {
			case *telegram.ChannelParticipantSelf:
				res[int(user.UserID)] = struct{}{}
			case *telegram.ChannelParticipantObj:
				res[int(user.UserID)] = struct{}{}
			case *telegram.ChannelParticipantAdmin:
				res[int(user.UserID)] = struct{}{}
			case *telegram.ChannelParticipantCreator:
				res[int(user.UserID)] = struct{}{}
			default:
				pp.Println(user)
				panic("что?")
			}
		}

		offset += 100
		pp.Println(offset, totalCount)
	}

	total = make([]int, 0, len(res))
	for k := range res {
		total = append(total, k)
	}

	sort.Ints(total)

	pp.Println(total, len(total))
}
