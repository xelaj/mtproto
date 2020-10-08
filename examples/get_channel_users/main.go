package main

import (
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
		// host address of mtproto server. actualy, it can'be mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is patrh to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile: "/home/me/go/src/github.com/xelaj/mtproto/keys/keys.pem",
		AppID:          94575,                              // app id, could be find at https://my.telegram.org
		AppHash:        "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
	})
	dry.PanicIfErr(err)

	// get this hash from channel invite link (after t.me/join/<HASH>)
	hash := "AAAAAEkCCtUBS84eqWdEeA"

	// syntax sugared method, more easy to read than default ways to solve some troubles
	//pp.Println(client.GetChannelInfoByInviteLink(hash))

	// get channel participants
	chat, err := client.GetChatInfoByHashLink(hash)
	dry.PanicIfErr(err)
	channelSimpleData, ok := chat.(*telegram.Channel)
	if !ok {
		panic("not a channel")
	}

	inChannel := telegram.InputChannel(&telegram.InputChannelObj{
		ChannelId:  channelSimpleData.Id,
		AccessHash: channelSimpleData.AccessHash,
	})

	resp, err := client.ChannelsGetParticipants(&telegram.ChannelsGetParticipantsParams{
		Channel: inChannel,
		Filter:  telegram.ChannelParticipantsFilter(&telegram.ChannelParticipantsRecent{}),
		Limit:   100,
	})
	dry.PanicIfErr(err)
	pp.Println(resp)
	//users := resp.(*telegram.ChannelsChannelParticipantsObj)
	//for i, participant := range users.Participants {
	//	user := participant.(*telegram.ChannelParticipantObj)
	//	pp.Println(i, user.UserId)
	//
	//}

}
