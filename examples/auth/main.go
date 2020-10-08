package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto"
	"github.com/xelaj/mtproto/keys"
	"github.com/xelaj/mtproto/telegram"
)

var client *telegram.Client

func main() {
	keyfile := "/home/r0ck3t/go/src/github.com/xelaj/mtproto/keys/keys.pem"
	TelegramPublicKeys, err := keys.ReadFromFile(keyfile)
	dry.PanicIfErr(err)

	m, err := mtproto.NewMTProto(mtproto.Config{
		AuthKeyFile: "~/.local/var/lib/mtproto/session.json.lol",
		ServerHost:  "149.154.167.50:443",
		PublicKey:   TelegramPublicKeys[0],
		AppID:       94575,
		AppHash:     "a3406de8d171bb422bb6ddf3bbd800e2",
	})
	if err != nil {
		panic(fmt.Errorf("create failed: %w", err))
	}
	client = &telegram.Client{m}

	err = client.CreateConnection()
	if err != nil {
		panic(fmt.Errorf("connect failed: %w", err))
	}

	resp, err := client.InvokeWithLayer(117, &telegram.InitConnectionParams{
		ApiID:          94575,
		DeviceModel:    "Unknown",
		SystemVersion:  "linux/amd64",
		AppVersion:     "0.0.1",
		SystemLangCode: "en",
		LangCode:       "en",
		Proxy:          nil,
		Params:         nil,
		Query:          &telegram.HelpGetConfigParams{},
	})
	dry.PanicIfErr(err)
	pp.Println("resp:", resp)

	switch resp.(type) {
	case *telegram.Config:
	default:
		panic(fmt.Sprintf("Got: %T", resp))
	}

	phoneNumber := os.Args[1]

	setCode, err := client.AuthSendCode(&telegram.AuthSendCodeParams{
		phoneNumber, 94575, "a3406de8d171bb422bb6ddf3bbd800e2", &telegram.CodeSettings{},
	})
	dry.PanicIfErr(err)
	pp.Println(resp)

	fmt.Print("Код авторизации:")
	code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	code = strings.Replace(code, "\n", "", -1)

	pp.Println(os.Args[2], setCode.PhoneCodeHash, code)
	pp.Println(client.AuthSignIn(&telegram.AuthSignInParams{
		phoneNumber, setCode.PhoneCodeHash, code,
	}))
}
