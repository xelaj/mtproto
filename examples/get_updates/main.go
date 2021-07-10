package main

import (
	"path/filepath"
	"sync"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"

	utils "github.com/xelaj/mtproto/examples/example_utils"
)

func main() {
	println("firstly, you need to authorize. after example 'auth', you will sign in")

	// helper variables
	appStorage := utils.PrepareAppStorageForExamples()
	sessionFile := filepath.Join(appStorage, "session.json")
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")

	client, err := telegram.NewClient(telegram.ClientConfig{
		// where to store session configuration. must be set
		SessionFile: sessionFile,
		// host address of mtproto server. Actually, it can be any mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is path to file with public keys, which you must get from https://my.telegram.org
		PublicKeysFile:  publicKeys,
		AppID:           94575,                              // app id, could be find at https://my.telegram.org
		AppHash:         "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
		InitWarnChannel: true,                               // if we want to get errors, otherwise, client.Warnings will be set nil
	})
	utils.ReadWarningsToStdErr(client.Warnings)
	dry.PanicIfErr(err)

	client.AddCustomServerRequestHandler(func(i interface{}) bool {
		pp.Println(i)
		return false
	})

	// we need to call updates.getState, after that telegram server will send you updates
	state, err := client.UpdatesGetState()
	dry.PanicIfErr(err)
	// this state could be useful, if you want to get old unread updates
	pp.Println(state)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
