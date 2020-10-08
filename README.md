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


## <p align="center">Features</p>

<div align="right">
<h3>Full native implementation</h3>
All code, from sending requests to encryption serialization is written on pure golang. You don't need to fetch any additional dependencies.
</br></br></br></br></br></br>
</div>

<div align="left">
<h3>Latest API version (117+)</h3>
Lib is supports all the API and MTProto features, including video calls and post comments. You can create additional pull request to push api updates! 
</br></br></br></br></br></br></br>
</div>

<div align="right">
<h3>Reactive API updates (generated from TL schema)</h3>
All changes in TDLib and Android client are monitoring to get the latest features and changes in TL schemas. New methods are creates by adding new lines into TL schema and updating generated code!
</br></br></br></br></br>
</div>

<div align="left">
<h3>Implements ONLY network tools</h3>
No more SQLite databases and caching unnecessary files, that **you** don't need. Also you can control how sessions are stored, auth process and literally everything that you need!
</br></br></br></br></br>
</div>

<div align="right">
<h3>Multiaccounting, Gateway mode</h3>
You can use more than 10 accounts at same time! <i>xelaj/MTProto</i> doesn't create huge overhead in memory or cpu consumption as TDLib. Thanks for that, you can create huge number of connection instances and don't worry about memory overload! 
</br></br></br></br></br>
</div>

## How to use

<!--
**СЮДА ИЗ asciinema ЗАПИХНУТЬ ДЕМОНСТРАЦИЮ**
![preview]({{ .PreviewUrl }})
-->

MTProto is really hard in implementation, but it's really easy to use. Basically, this lib sends serialized structures to Telegram servers (just like gRPC, but from Telegram LLC.). It looks like this:

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

Not so hard, huh? But there is even easier way to send request, which is included in TL API specification:

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

You do not need to think about encryption, key exchange, saving and restoring session, and more routine things. It is already implemented just for you.

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

That's it! You don't need to do anything more!

### What is InvokeWithLayer?

It's Telegram specific feature. I you want to create client instance and get information about the current servers configuration, you need to do something like this:

```go
    resp, err := client.InvokeWithLayer(apiVersion, &telegram.InitConnectionParams{
        ApiID:          124100,
        DeviceModel:    "Unknown",
        SystemVersion:  "linux/amd64",
        AppVersion:     "0.1.0",
	// just use "en", any other language codes will receive error. See telegram docs for more info.
        SystemLangCode: "en",
        LangCode:       "en", 
        // HelpGetConfig() is ACTUAL request, but wrapped in IvokeWithLayer
        Query:          &telegram.HelpGetConfigParams{},
    })
```

Why? We don't know! This method is described in Telegram API docs, any other starting requests will receive error.

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

### Docs are empty. Why?

There is a pretty huge chunk of documentation. We are ready to describe every method and object, but its requires a lot of work. Although **all** methods are **already** described [here](https://core.telegram.org/methods).

### Does this project support Windows?

Technically — yes. In practice — components don't require specific architecture, but we didn't test it yet. If you have any problems running it, just create an issue, we will try to help.

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
