# MTProto

![help wanted](https://img.shields.io/badge/-help%20wanted-success)
[![godoc reference](https://pkg.go.dev/badge/github.com/xelaj/mtproto?status.svg)](https://pkg.go.dev/github.com/xelaj/mtproto)
[![Go Report Card](https://goreportcard.com/badge/github.com/xelaj/mtproto)](https://goreportcard.com/report/github.com/xelaj/mtproto)
[![license MIT](https://img.shields.io/badge/license-MIT-green)](https://github.com/xelaj/mtproto/blob/master/README.md)
[![chat telegram](https://img.shields.io/badge/chat-telegram-0088cc)](https://bit.ly/2xlsVsQ)
![version v0.1.0](https://img.shields.io/badge/version-v0.1.0-red)
![unstable](https://img.shields.io/badge/stability-unstable-yellow)
<!--
code quality
golangci
contributors
go version
gitlab pipelines
-->

![FINALLY!](docs/assets/finally.jpg) Full-native implementation of MTProto protocol on Golang!

**english** [русский](https://github.com/xelaj/mtproto/blob/master/docs/ru_RU/README.md)

<p align="center">
<img src="docs/assets/MuffinMan-AgADRAADO2AkFA.gif"/>
</p>

## <p align="center">Features</p>

<div align="right">
<h3>Full native implementation</h3>
<img src="docs/assets/ezgif-3-a6bd45965060.gif" align="right"/>
All code  from sending requests to encryption serialisation is written on pure golang. You dont need to download any additional dependencies.
</br></br></br></br></br></br>
</div>

<div align="left">
<h3>Latest API version (117+)</h3>
<img src="docs/assets/ezgif-3-19ced73bc71f.gif" align="left"/>
It supports all the API and MTProto features, including video calls and post comments. You can create additional pull request to renew the data (???)! 
</br></br></br></br></br></br></br>
</div>

<div align="right">
<h3>Reactive API updates (generated from TL schema)</h3>
<img src="docs/assets/ezgif-3-5b6a808d2774.gif" align="right"/>
All the changes in TDLib and Android and being monitored to get the latests features and changes in TL schemas. New methods are created by adding new lines into the schema and updating generated code!
</br></br></br></br></br>
</div>

<div align="left">
<h3>Implements ONLY network tools</h3>
<img src="docs/assets/ezgif-3-3ac8a3ea5713.gif" align="left"/>
No SQLite databases and caching files are required. You can use only things you need. Also you can control how sessions are stored, authorisation process and literally everything you need!
</br></br></br></br></br>
</div>

<div align="right">
<h3>Multiaccounting, Gateway mode</h3>
<img src="docs/assets/ezgif-3-7bcf6dc78388.gif" align="right"/>
You can use more than 10 accounts at a time! _xelaj/MTProto_ does not create big overhead in memory and cpu consumption. Because of that you should not worry about having huge number of connection instances!  
</br></br></br></br></br>
</div>

## How to use

<!--
**СЮДА ИЗ asciinema ЗАПИХНУТЬ ДЕМОНСТРАЦИЮ**
![preview]({{ .PreviewUrl }})
-->

MTProto has a quiet hard implementation, but is quiet easy to use. In fact, you are sending serialized structures to Telegram servers (just like gRPC, but from Telegram LLC.). It looks like this:

```go
func main() {
    client := &Telegram.NewClient()
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

But there is even easier way to send request, which is included in TL API specification:

```go
func main() {
    client := &Telegram.NewClient()
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

You do not need to think about encryption, key exchange, saving and restoring a session. It is already implemented for you.

**Code examples are [here](https://github.com/xelaj/mtproto/blob/master/examples)**

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

That's it! Simple!

### What is InvokeWithLayer?

It is Telegram specific feature. I you want to create and get information about the current servers configuration, you need to do this:

```go
    resp, err := client.InvokeWithLayer(apiVersion, &telegram.InitConnectionParams{
        ApiID:          124100,
        DeviceModel:    "Unknown",
        SystemVersion:  "linux/amd64",
        AppVersion:     "0.1.0",
        SystemLangCode: "en",
        LangCode:       "en",
        Proxy:          nil,
        Params:         nil,
        // HelpGetConfig() is ACTUAL request, but wrapped in IvokeWithLayer
        Query:          &telegram.HelpGetConfigParams{},
    })
```

### How to use phone authorization?

**Example [here](https://github.com/xelaj/mtproto/blob/master/examples/auth)**


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

    // Можно выбрать любой удобный вам способ ввода,
    // базовые параметры сессии можно сохранить в любом месте
	fmt.Print("Auth code:")
	code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    code = strings.Replace(code, "\n", "", -1)

    // это весь процесс авторизации!
    fmt.Println(client.AuthSignIn(yourPhone, resp.PhoneCodeHash, code))
}
```

That's it! You don't need any cycles, the code is fully ready for asynchronous execution. You just need to follow the official Telegram API documentation.
 
### Docs are empty. Why?

It is a pretty huge chunk of documentation. We are ready to describe every method and object, but its requires a lot of work. Although **all** methods are **already** described [here](https://core.telegram.org/methods).

### Does this project support Windows?

Yes in theory. The components don't require specific architecture. Although we did not test it. If you have any problems running it, just create an issue, we will help.

## Who use it

## Contributing

Please read [contributing guide](https://github.com/xelaj/mtproto/blob/master/doc/en_US/CONTRIBUTING.md) if you want to help. And the help is very necessary!

## TODO

[ ]

## Authors

* **Richard Cooper** <[rcooper.xelaj@protonmail.com](mailto:rcooper.xelaj@protonmail.com)>

## License

<b style="color:red">WARNING!</b> This project is only maintained by Xelaj inc., however copyright of this source code **IS NOT** owned by Xelaj inc. at all. If you want to connect with code owners, write mail to <a href="mailto:up@khsfilms.ru">this email</a>. For all other questions like any issues, PRs, questions, etc. Use GitHub issues, or find email on official website.

This project is licensed under the MIT License - see the [LICENSE](https://github.com/xelaj/mtproto/blob/master/doc/en_US/LICENSE.md) file for details
