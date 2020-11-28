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
		// host address of mtproto server. Actually, it can'be mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is patrh to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile: "/home/me/.local/var/lib/mtproto/tg_public_keys.pem",
		AppID:          94575,                              // app id, could be find at https://my.telegram.org
		AppHash:        "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
	})
	dry.PanicIfErr(err)

	resp, err := client.AccountInitTakeoutSession(&telegram.AccountInitTakeoutSessionParams{
		Contacts: true,
	})
	dry.PanicIfErr(err)

	res, err := client.MakeRequest(
		&telegram.InvokeWithTakeoutParams{
			TakeoutID: resp.ID,
			Query:     &telegram.ContactsGetSavedParams{},
		},
	)

	pp.Sprintln(res, err)

	contacts, err := client.ContactsGetContacts(0)
	dry.PanicIfErr(err)
	c := contacts.(*telegram.ContactsContactsObj)
	pp.Println(int(c.SavedCount), len(c.Users), len(c.Contacts))
}
