# MTProto

[![godoc reference](https://pkg.go.dev/badge/github.com/xelaj/mtproto?status.svg)](https://pkg.go.dev/github.com/xelaj/mtproto)
[![Go Report Card](https://goreportcard.com/badge/github.com/xelaj/mtproto)](https://goreportcard.com/report/github.com/xelaj/mtproto)
[![codecov](https://codecov.io/gh/xelaj/mtproto/branch/master/graph/badge.svg)](https://codecov.io/gh/xelaj/mtproto)
[![license MIT](https://img.shields.io/badge/license-MIT-green)](https://github.com/xelaj/mtproto/blob/main/README.md)
[![chat telegram](https://img.shields.io/badge/chat-telegram-0088cc)](https://bit.ly/2xlsVsQ)
![version v1.0.0](https://img.shields.io/badge/version-v1.0.0-success)
![unstable](https://img.shields.io/badge/stability-stable-success)
<!--
code quality
golangci
contributors
go version
gitlab pipelines
-->

![FINALLY!](/docs/assets/finally.jpg) Full-native implementation of MTProto protocol on Golang!

**english** [русский](https://github.com/xelaj/mtproto/blob/main/docs/ru_RU/README.md) [简体中文](https://github.com/xelaj/mtproto/blob/main/docs/zh_CN/README.md)

<p align="center">
<img src="https://i.ibb.co/yYsPxhW/Muffin-Man-Ag-ADRAADO2-Ak-FA.gif"/>
</p>

## <p align="center">Features</p>

<div align="right">
<h3>Full native implementation</h3>
<img src="https://i.ibb.co/9Vfz6hj/ezgif-3-a6bd45965060.gif" align="right"/>
All code, from sending requests to encryption serialization is written on pure golang. You don't need to fetch any additional dependencies.
<br/><br/><br/><br/><br/>
</div>

<div align="left">
<h3>Latest API version (117+)</h3>
<img src="https://i.ibb.co/nw84W4h/ezgif-3-19ced73bc71f.gif" align="left"/>
Lib is supports all the API and MTProto features, including video calls and post comments. You can create additional pull request to push api updates!
<br/><br/><br/><br/><br/><br/><br/>
</div>

<div align="right">
<h3>Reactive API updates (generated from TL schema)</h3>
<img src="https://i.ibb.co/9WXrHq8/ezgif-3-5b6a808d2774.gif" align="right"/>
All changes in TDLib and Android client are monitoring to get the latest features and changes in TL schemas. New methods are creates by adding new lines into TL schema and updating generated code!
<br/><br/><br/><br/><br/>
</div>

<div align="left">
<h3>Implements ONLY network tools</h3>
<img src="https://i.ibb.co/bLj3PHx/ezgif-3-3ac8a3ea5713.gif" align="left"/>
No more SQLite databases and caching unnecessary files, that <b>you</b> don't need. Also you can control how sessions are stored, auth process and literally everything that you want to!
<br/><br/><br/><br/><br/>
</div>

<div align="right">
<h3>Multiaccounting, Gateway mode</h3>
<img src="https://i.ibb.co/8XbKRPG/ezgif-3-7bcf6dc78388.gif" align="right"/>
You can use more than 10 accounts at same time! <i>xelaj/MTProto</i> doesn't create huge overhead in memory or cpu consumption as TDLib. Thanks for that, you can create huge number of connection instances and don't worry about memory overload!
<br/><br/><br/><br/><br/>
</div>

## How to use

<!--
TODO: **HERE GOES asciinema DEMO**
![preview]({{ .PreviewUrl }})
-->

MTProto is really hard in implementation, but it's really easy to use. Basically, this lib sends serialized structures to Telegram servers (just like gRPC, but from Telegram LLC.). It looks like this:

```go
func main() {
    client := &telegram.NewClient()
    // for each method there is specific struct for serialization (<method_name>Params{})
    result, err := client.MakeRequest(&telegram.GetSomeInfoParams{FromChatId: 12345})
    if err != nil {
        panic(err)
    }

    resp, ok := result.(*SomeResponseObject)
    if !ok {
        panic("Oh no! Wrong type!")
    }
}
```

Not so hard, huh? But there is even easier way to send request, which is included in TL API specification:

```go
func main() {
    client := &telegram.NewClient()
    resp, err := client.GetSomeInfo(12345)
    if err != nil {
        panic(err)
    }

    // resp will be already asserted as described in TL specs of API
    // if _, ok := resp.(*SomeResponseObject); !ok {
    //     panic("No way, we found a bug! Create new issue!")
    // }

    println(resp.InfoAboutSomething)
}
```

You do not need to think about encryption, key exchange, saving and restoring session, and more routine things. It is already implemented just for you.

**Code examples are [here](https://github.com/xelaj/mtproto/blob/main/examples)**

**Full docs are [here](https://pkg.go.dev/github.com/xelaj/mtproto)**

## Getting started

### Simple How-To

Installation is simple. Just do `go get`:

``` bash
go get github.com/xelaj/mtproto
```

After that you can generate source structures of methods and functions if you wish to. To do it, use `go generate`

``` bash
go generate github.com/xelaj/mtproto
```

That's it! You don't need to do anything more!

### What is InvokeWithLayer?

It's Telegram specific feature. If you want to create client instance and get information about the current server's configuration, you need to do something like this:

```go
resp, err := client.InvokeWithLayer(apiVersion, &telegram.InitConnectionParams{
    ApiID:          124100,
    DeviceModel:    "Unknown",
    SystemVersion:  "linux/amd64",
    AppVersion:     "0.1.0",
    // just use "en", any other language codes will receive error. See telegram docs for more info.
    SystemLangCode: "en",
    LangCode:       "en",
    // HelpGetConfig() is ACTUAL request, but wrapped in InvokeWithLayer
    Query:          &telegram.HelpGetConfigParams{},
})
```

Why? We don't know! This method is described in Telegram API docs, any other starting requests will receive error.

### How to use phone authorization?

**Example [here](https://github.com/xelaj/mtproto/blob/main/examples/auth)**

```go
func AuthByPhone() {
    resp, err := client.AuthSendCode(
        yourPhone,
        appID,
        appHash,
        &telegram.CodeSettings{},
    )
    if err != nil {
        panic(err)
    }

    // You can make any way to enter verification code, like in
    // http requests, or what you like. You just need to call two
    // requests, that's main method.
    fmt.Print("Auth code:")
    code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    code = strings.Replace(code, "\n", "", -1)

    // this is ALL process of authorization! :)
    fmt.Println(client.AuthSignIn(yourPhone, resp.PhoneCodeHash, code))
}
```

That's it! You don't need any cycles, code is ready-to-go for async execution. You just need to follow the official Telegram API documentation.

### Telegram Deeplinks

Want to deal those freaky `tg://` links? See [`deeplinks` package](https://github.com/xelaj/mtproto/blob/main/telegram/deeplinks), here is the simplest how-to:

``` go
package main

import (
    "fmt"

    "github.com/xelaj/mtproto/telegram/deeplinks"
)

func main() {
    link, _ := deeplinks.Resolve("t.me/xelaj_developers")
    // btw, ResolveParameters is just struct for tg://resolve links, not all links are resolve
    resolve := link.(*deeplinks.ResolveParameters)
    fmt.Printf("Oh! Looks like @%v is the best developers channel in telegram!\n", resolve.Domain)
}
```

### Docs are empty. Why?

There is a pretty huge chunk of documentation. We are ready to describe every method and object, but it requires a lot of work. Although **all** methods are **already** described [here](https://core.telegram.org/methods).

### Does this project support Windows?

Technically — yes. In practice — components don't require specific architecture, but we didn't test it yet. If you have any problems running it, just create an issue, we will try to help.

### Why Telegram API soooo unusable?

Well... Read [this issue](https://github.com/ton-blockchain/ton/issues/31) about TON source code. Use google translate, this issue will answer to all your questions.

## Who use it

## Contributing

Please read [contributing guide](https://github.com/xelaj/mtproto/blob/main/.github/CONTRIBUTING.md) if you want to help. And the help is very necessary!

**Don't want code?** Read [this](https://github.com/xelaj/mtproto/blob/main/.github/SUPPORT.md) page! We love nocoders!

## Security bugs?

Please, don't create issue which describes security bug, this can be too offensive! Instead, please read [this notification](https://github.com/xelaj/mtproto/blob/main/.github/SECURITY.md) and follow that steps to notify us about problem.

## TODO

- [x] Basic MTProto implementation
- [x] Implement all Methods for latest layer
- [x] Make TL Encoder/Decoder
- [x] Get away from panics in parsing TL
- [ ] Support MTProxy
- [ ] Support socks5 as well
- [ ] Multiple tests
- [ ] Write amazing docs

## Authors

* **Richard Cooper** <[rcooper.xelaj@protonmail.com](mailto:rcooper.xelaj@protonmail.com)>
* **Anton Larionov** <[Anton.Larionov@infobip.com](mailto:Anton.Larionov@infobip.com)>
* **Arthur Petukhovsky** <[petuhovskiy@yandex.ru](mailto:petuhovskiy@yandex.ru)>
* **Roman Timofeev** <[timofeev@uteka.ru](mailto:timofeev@uteka.ru)>
* **Artem** <[webgutar@gmail.com](mailto:webgutar@gmail.com)>
* **Bo-Yi Wu** <[appleboy.tw@gmail.com](mailto:appleboy.tw@gmail.com)>
* **0xflotus** <[0xflotus@gmail.com](mailto:0xflotus@gmail.com)>
* **Luclu7** <[me@luclu7.fr](mailto:me@luclu7.fr)>
* **Vladimir Stolyarov** <[xakep6666@gmail.com](mailto:xakep6666@gmail.com)>
* **grinrill** [@grinrill](https://github.com/grinrill)
* **kulallador** <[ilyastalk@bk.ru](ilyastalk@bk.ru)>
* **rs** <[yuiop1955@mail.ru](mailto:yuiop1955@mail.ru)>

## License

<b style="color:red">WARNING!</b> This project is only maintained by Xelaj inc., however copyright of this source code **IS NOT** owned by Xelaj inc. at all. If you want to connect with code owners, write mail to <a href="mailto:up@khsfilms.ru">this email</a>. For all other questions like any issues, PRs, questions, etc. Use GitHub issues, or find email on official website.

This project is licensed under the MIT License - see the [LICENSE](https://github.com/xelaj/mtproto/blob/main/LICENSE.md) file for details

<!--

V2UndmUga25vd24gZWFjaCBvdGhlciBmb3Igc28gbG9uZwpZb3
VyIGhlYXJ0J3MgYmVlbiBhY2hpbmcgYnV0IHlvdSdyZSB0b28g
c2h5IHRvIHNheSBpdApJbnNpZGUgd2UgYm90aCBrbm93IHdoYX
QncyBiZWVuIGdvaW5nIG9uCldlIGtub3cgdGhlIGdhbWUgYW5k
IHdlJ3JlIGdvbm5hIHBsYXkgaXQKQW5kIGlmIHlvdSBhc2sgbW
UgaG93IEknbSBmZWVsaW5nCkRvbid0IHRlbGwgbWUgeW91J3Jl
IHRvbyBibGluZCB0byBzZWU=

-->
