# MTProto

[![godoc reference](https://pkg.go.dev/badge/github.com/xelaj/mtproto?status.svg)](https://pkg.go.dev/github.com/xelaj/mtproto)
[![Go Report Card](https://goreportcard.com/badge/github.com/xelaj/mtproto)](https://goreportcard.com/report/github.com/xelaj/mtproto)
[![codecov](https://codecov.io/gh/xelaj/mtproto/branch/master/graph/badge.svg)](https://codecov.io/gh/xelaj/mtproto)
[![license MIT](https://img.shields.io/badge/license-MIT-green)](https://github.com/xelaj/mtproto/blob/main/README.md)
[![chat telegram](https://img.shields.io/badge/chat-telegram-0088cc)](https://bit.ly/2xlsVsQ)
![version v1.0.0](https://img.shields.io/badge/version-v1.0.0-success)
![unstable](https://img.shields.io/badge/stability-stable-success)
<!--
contributors
go version
gitlab pipelines
-->

![FINALLY!](/docs/assets/finally.jpg) MTProto 协议的原生 Go 实现！

[english](https://github.com/xelaj/mtproto/blob/main/README.md) [русский](https://github.com/xelaj/mtproto/blob/main/docs/ru_RU/README.md) **简体中文**

<p align="center">
<img src="https://i.ibb.co/yYsPxhW/Muffin-Man-Ag-ADRAADO2-Ak-FA.gif"/>
</p>

## <p align="center">特性</p>

<div align="right">
<h3>原生实现</h3>
<img src="https://i.ibb.co/9Vfz6hj/ezgif-3-a6bd45965060.gif" align="right"/>
从发送请求到加密序列化的所有代码均是使用纯 Go 编写，你无需额外拉取任何依赖。
<br/><br/><br/><br/><br/><br/><br/><br/>
</div>

<div align="left">
<h3>最新的 API 版本（117+）</h3>
<img src="https://i.ibb.co/nw84W4h/ezgif-3-19ced73bc71f.gif" align="left"/>
支持所有的 API 和 MTProto 功能，包括视频通话和发表评论。如果你愿意的话，你也可以通过创建 PR 来更新这些 API！
<br/><br/><br/><br/><br/><br/><br/>
</div>

<div align="right">
<h3>响应式 API 更新（生成自 TL schema）</h3>
<img src="https://i.ibb.co/9WXrHq8/ezgif-3-5b6a808d2774.gif" align="right"/>
TDLib 和 Android 客户端的所有改动都在监控当中，确保了最新的功能和 TL schemas 中的变动都能够即时地同步进来。
<br/><br/><br/><br/><br/>
</div>

<div align="left">
<h3>仅是一个网络工具</h3>
<img src="https://i.ibb.co/bLj3PHx/ezgif-3-3ac8a3ea5713.gif" align="left"/>
不再需要 SQLite 数据库和其他多余的缓存文件。你还可以控制会话的存储方式和身份验证等一切内容！
<br/><br/><br/><br/><br/><br/><br/><br/>
</div>

<div align="right">
<h3>多账户，网关模式</h3>
<img src="https://i.ibb.co/8XbKRPG/ezgif-3-7bcf6dc78388.gif" align="right"/>
你可以同时使用多达 10 余个账户！<i>xelaj/MTProto</i> 不会像 TDLib 那样在内存或 CPU 消耗方面产生巨大的开销。因此，你可以创建大量的连接实例，而不必担心内存过载！
<br/><br/><br/><br/><br/>
</div>

## 如何使用

MTProto 真的很难实现，但却极易使用。简单地说，此库只是将序列化的结构发送至 Telegram 服务器（就像 gRPC 一样，但来自 Telegram LLC）。就像这样：

```go
func main() {
    client := &Telegram.NewClient()
    // 对于每一个方法，都有一个具体的结构体用来序列化（<method_name>Params{}）
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

不难吧？其实在 TL API 规范中还有更简单的发送请求的方法：

```go
func main() {
    client := &Telegram.NewClient()
    resp, err := client.GetSomeInfo(12345)
    if err != nil {
        panic(err)
    }

    // resp 已被按照 TLS API 规范中的描述来断言，可以直接使用
    // if _, ok := resp.(*SomeResponseObject); !ok {
    //     panic("No way, we found a bug! Create new issue!")
    // }

    println(resp.InfoAboutSomething)
}
```

你无需考虑加密、密钥交换、保存和还原会话以及其他常规事务，我们全都替你处理好了。

**示例代码在[这里](https://github.com/xelaj/mtproto/blob/main/examples)**

**完整的文档在[这里](https://pkg.go.dev/github.com/xelaj/mtproto)**

## 开始使用

### 安装

安装方式很简单，简单地执行 `go get`：

``` bash
go get github.com/xelaj/mtproto
```

之后，你可以通过执行 `go generate` 来根据需要生成方法和函数的源结构：

``` bash
go generate github.com/xelaj/mtproto
```

就这么简单！你不需要做其他的任何事情了！

### 什么是 InvokeWithLayer？

这是 Telegram 的一个功能，如果你想创建一个包含当前服务器配置信息的客户端，那么你需要这么做：

```go
resp, err := client.InvokeWithLayer(apiVersion, &telegram.InitConnectionParams{
    ApiID:          124100,
    DeviceModel:    "Unknown",
    SystemVersion:  "linux/amd64",
    AppVersion:     "0.1.0",
    // 请使用"en"，其他任何值都会导致错误产生
    SystemLangCode: "en",
    LangCode:       "en",
    // HelpGetConfig() 是一个确切的包含在 InvokeWithLayer 内的请求
    Query:          &telegram.HelpGetConfigParams{},
})
```

为什么？我们不知道！Telegram API 文档中介绍了此方法，其他任何启动请求都将收到错误。

### 如何使用电话验证？

**示例代码在[这里](https://github.com/xelaj/mtproto/blob/main/examples/auth)**

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

    // 你可以通过任意方式来获取收到的验证码，比如通过 HTTP 请求等。
    fmt.Print("Auth code:")
    code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    code = strings.Replace(code, "\n", "", -1)

    // 这就是电话验证的全部过程！
    fmt.Println(client.AuthSignIn(yourPhone, resp.PhoneCodeHash, code))
}
```

就这么简单！你不需要任何循环，异步执行代码开箱即用。你只需要遵循官方的 Telegram API 文档即可。

### Telegram Deeplinks

想处理那些奇葩的 `tg://` 链接吗？请查看 [`deeplinks`](https://github.com/xelaj/mtproto/blob/main/telegram/deeplinks) 包。如下是简单的示例：

``` go
package main

import (
    "fmt"

    "github.com/xelaj/mtproto/telegram/deeplinks"
)

func main() {
    link, _ := deeplinks.Resolve("t.me/xelaj_developers")
    // 顺便说一下，ResolveParameters 只是 tg://resolve 链接的结构体，并非所有链接都是 resolve 的
    resolve := link.(*deeplinks.ResolveParameters)
    fmt.Printf("Oh! Looks like @%v is the best developers channel in telegram!\n", resolve.Domain)
}
```

### 为什么文档是空的？

事实上是有相当庞大的文档的。我们准备描述每一个方法和对象，但是那将耗费巨大的工作量。尽管所有的方法全都已经描述在了[这里](https://core.telegram.org/methods)。

### 这个项目支持 Windows 吗？

从技术上讲，支持。在实践中，组件不需要具体的系统架构，但是我们目前并没有测试这一点。如果你遇到了任何问题，请随时创建 issue 来反馈，我们会尽力帮助。

### 为什么 Telegram 如此不稳定？

请使用 [ Google 翻译](https://translate.google.com/) 来查看 [这个 issue](https://github.com/ton-blockchain/ton/issues/31) 它可以回答你的所有问题。

## 谁在使用

## 贡献

如果你愿意提供帮助，请阅读我们的[贡献指南](https://github.com/xelaj/mtproto/blob/main/.github/CONTRIBUTING.md)。

不想写代码？请阅读[这个](https://github.com/xelaj/mtproto/blob/main/.github/SUPPORT.md)页面，我们同样欢迎 nocoders！

## 安全漏洞？

请千万不要创建 issue 来反馈安全漏洞，因为这样会影响到很多人。请通过[阅读](https://github.com/xelaj/mtproto/blob/main/.github/SECURITY.md)这个并遵循里面提到的步骤来通知我们。

## TODO

- [x] 基础 MTProto 实现
- [x] 实现最新 layer 中的所有方法
- [x] 实现 TL 编码器和解码器
- [x] 避免解析 TL 时产生 panics
- [ ] 支持 MTProxy
- [ ] 支持 socks5
- [ ] 添加测试
- [ ] 丰富文档

## 作者

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

## 许可证

**WARNING!** This project is only maintained by Xelaj inc., however copyright of this source code **IS NOT** owned by Xelaj inc. at all. If you want to connect with code owners, write mail to <a href="mailto:up@khsfilms.ru">this email</a>. For all other questions like any issues, PRs, questions, etc. Use GitHub issues, or find email on official website.

This project is licensed under the MIT License - see the [LICENSE](https://github.com/xelaj/mtproto/blob/main/docs/en_US/LICENSE.md) file for details

<!--

V2UndmUga25vd24gZWFjaCBvdGhlciBmb3Igc28gbG9uZwpZb3
VyIGhlYXJ0J3MgYmVlbiBhY2hpbmcgYnV0IHlvdSdyZSB0b28g
c2h5IHRvIHNheSBpdApJbnNpZGUgd2UgYm90aCBrbm93IHdoYX
QncyBiZWVuIGdvaW5nIG9uCldlIGtub3cgdGhlIGdhbWUgYW5k
IHdlJ3JlIGdvbm5hIHBsYXkgaXQKQW5kIGlmIHlvdSBhc2sgbW
UgaG93IEknbSBmZWVsaW5nCkRvbid0IHRlbGwgbWUgeW91J3Jl
IHRvbyBibGluZCB0byBzZWU=

-->
