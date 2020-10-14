package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("second argument must be phone number!")
		os.Exit(1)
	}
	phoneNumber := os.Args[1]

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

	setCode, err := client.AuthSendCode(&telegram.AuthSendCodeParams{
		phoneNumber, 94575, "a3406de8d171bb422bb6ddf3bbd800e2", &telegram.CodeSettings{},
	})
	dry.PanicIfErr(err)
	pp.Println(setCode)

	fmt.Print("Код авторизации:")
	code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	code = strings.Replace(code, "\n", "", -1)

	pp.Println(client.AuthSignIn(&telegram.AuthSignInParams{
		phoneNumber, setCode.PhoneCodeHash, code,
	}))
}
