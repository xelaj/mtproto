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

![FINALLY!](/docs/assets/finally.jpg) Полностью нативная имплементация MTProto на Golang!


[english](https://github.com/xelaj/mtproto/blob/main/README.md) **русский** [简体中文](https://github.com/xelaj/mtproto/blob/main/docs/zh_CN/README.md)

<p align="center">
<img src="https://i.ibb.co/yYsPxhW/Muffin-Man-Ag-ADRAADO2-Ak-FA.gif"/>
</p>

## <p align="center">Фичи</p>

<div align="right">
<h3>Полностью нативная реализация</h3>
<img src="https://i.ibb.co/9Vfz6hj/ezgif-3-a6bd45965060.gif" align="right"/>
Вся библиотека начиная с отправки запросов и шифрования и заканчивая сериализацией шифровния написаны исключительно на golang. Для работы с библиотекой не требуется никаких лишних зависимостей.
<br/><br/><br/><br/><br/>
</div>

<div align="left">
<h3>Самая свежая версия API (117+)</h3>
<img src="https://i.ibb.co/nw84W4h/ezgif-3-19ced73bc71f.gif" align="left"/>
Реализована поддержка всех возможностей API Telegram и MTProto включая функцию видеозвонков и комментариев к постам. Вы можете сделать дополнительный pull request на обновление данных!
<br/><br/><br/><br/><br/><br/><br/>
</div>

<div align="right">
<h3>Реактивные обновления (сгенерировано из TL спецификаций)</h3>
<img src="https://i.ibb.co/9WXrHq8/ezgif-3-5b6a808d2774.gif" align="right"/>
Все изменения в клиентах TDLib и Android мониторятся на предмет появления новых фич и изменений в TL схемах. Новые методы и объекты появляются просто по добавлению новых строк в схеме и обновления сгенерированного кода!
<br/><br/><br/><br/><br/>
</div>

<div align="left">
<h3>Implements ONLY network tools</h3>
<img src="https://i.ibb.co/bLj3PHx/ezgif-3-3ac8a3ea5713.gif" align="left"/>
Никаких SQLite баз данных и кеширования ненужных <b>вам</b> файлов. Вы можете использовать только тот функционал, который вам нужен. Вы так же можете управлять способом сохранения сессий, процессом авторизации, буквально всем, что вам необходимо!
<br/><br/><br/><br/><br/>
</div>

<div align="right">
<h3>Multiaccounting, Gateway mode</h3>
<img src="https://i.ibb.co/8XbKRPG/ezgif-3-7bcf6dc78388.gif" align="right"/>
Вы можете использовать больше 10 аккаунтов одновременно! xelaj/MTProto не создает большого оверхеда по вычислительным ресурсам, поэтому вы можете иметь огромное количество инстансов соединений и не переживать за перерасход памяти!
<br/><br/><br/><br/><br/>
</div>

## How to use

<!--
**СЮДА ИЗ asciinema ЗАПИХНУТЬ ДЕМОНСТРАЦИЮ**
![preview]({{ .PreviewUrl }})
-->

MTProto очень сложен в реализации, но при этом очень прост в использовании. По сути вы общаетесь с серверами Telegram посредством отправки сериализованых структур (аналог gRPC, разработанный Telegram llc.). Выглядит это примерно так:

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

Однако, есть более простой способ отправить запрос, который уже записан в TL спецификации API:

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

Вам не стоит задумываться о реализации шифрования, обмена ключами, сохранении и восстановлении сессии, все уже сделано за вас.

**Примеры кода [здесь](https://github.com/xelaj/mtproto/blob/main/examples)**

**Полная документация [здесь](https://pkg.go.dev/github.com/xelaj/mtproto)**

## Getting started

### Simple How-To

Все как обычно, вам необходимо загрузить пакет с помощью `go get`:

``` bash
go get github.com/xelaj/mtproto
```

Далее по желанию вы можете заново сгенерировать исходники структур методов и функций, для этого используйте команду `go generate`

``` bash
go generate github.com/xelaj/mtproto
```

Все! Больше ничего и не надо!

### Что за InvokeWithLayer?

Это специфическая особенность Telegram, для создания соединения и получения информации о текущей конфигурации серверов, нужно сделать что-то подобное:

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

Почему? А хер его знает! Этот метод описан в документации Telegram API, любые другие стартовые запросы получат ошибку.

### Как произвести авторизацию по телефону?

**Пример [здесь](https://github.com/xelaj/mtproto/blob/main/examples/auth)**

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

Все! вам не требуется никаких циклов или чего-то подобного, код уже готов к асинхронному выполнению, вам нужно только выполнить действия прописанные в документации к Telegram API

### Telegram Deeplinks

Нужно работать с этими стрёмными `tg://` ссылками? Загляните в [пакет `deeplinks`](https://github.com/xelaj/mtproto/blob/main/telegram/deeplinks), вот самый простейший пример:

``` go
package main

import (
    "fmt"

    "github.com/xelaj/mtproto/telegram/deeplinks"
)

func main() {
    link, _ := deeplinks.Resolve("t.me/xelaj_developers")
    // кстати говоря, ResolveParameters это просто структура для ссылок tg://resolve, не все ссылки это исключительно ResolveParameters
    resolve := link.(*deeplinks.ResolveParameters)
    fmt.Printf("Ого! похоже что @%v это лучший девелоперский чат в телеге!\n", resolve.Domain)
}
```

### Документация пустует! Почему?

Объем документации невероятно огромен. Мы бы готовы задокументировать каждый метод и объект, но это огромное количество работы. Несмотря на это, **все** методы **уже** описаны [здесь](https://core.telegram.org/methods), вы можете так же спокойно их

### Работает ли этот проект под Windows?

Технически — да. Компоненты не были заточены под определенную архитектуру. Однако, возможности протестировать у разработчиков не было. Если у вас возникли проблемы, напишите в issues, мы постараемся помочь


### Почему Telegram API НАСТОЛЬКО нестабильное?

Ну... Как сказать... А гляньте лучше [вот это ишью](https://github.com/ton-blockchain/ton/issues/31) про исходники TON (ныне закрытый проект telegram). Оно вам расскажет обо всех проблемах и самого телеграма, и его апи. И его разработчиках.

## Who use it

## Contributing

Please read [contributing guide](https://github.com/xelaj/mtproto/blob/main/docs/ru_RU/CONTRIBUTING.md) if you want to help. And the help is very necessary!

**Don't want code?** Read [this](https://github.com/xelaj/mtproto/blob/main/.github/SUPPORT.md) page! We love nocoders!

## Критические уязвимости?

Ага, мы стараемся подходить к ним серьезно. Пожалуйста, не создавайте ищью, описывающие ошибку безопасности, это может быть ОЧЕНЬ небезопасно! Вместо этого рекомендуем глянуть на [эту страничку](https://github.com/xelaj/mtproto/blob/main/.github/SECURITY.md) и следовать инструкции по уведомлению.

## TODO

- [x] Базовая реализация MTProto
- [x] Реализовать все методы для последнего слоя
- [x] Сделать TL Encoder/Decoder
- [x] Убрать все паники при парсинге TL
- [ ] Поддержка MTProxy
- [ ] Поддержка Socks5 из коробки (вообще вы и так можете гнать запросы через socks5)
- [ ] Капитальное тестирование (80%+)
- [ ] Написать **свою** суперскую документацию

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

This project is licensed under the MIT License - see the [LICENSE](https://github.com/xelaj/mtproto/blob/main/docs/ru_RU/LICENSE.md) file for details

<!--

V2UndmUga25vd24gZWFjaCBvdGhlciBmb3Igc28gbG9uZwpZb3
VyIGhlYXJ0J3MgYmVlbiBhY2hpbmcgYnV0IHlvdSdyZSB0b28g
c2h5IHRvIHNheSBpdApJbnNpZGUgd2UgYm90aCBrbm93IHdoYX
QncyBiZWVuIGdvaW5nIG9uCldlIGtub3cgdGhlIGdhbWUgYW5k
IHdlJ3JlIGdvbm5hIHBsYXkgaXQKQW5kIGlmIHlvdSBhc2sgbW
UgaG93IEknbSBmZWVsaW5nCkRvbid0IHRlbGwgbWUgeW91J3Jl
IHRvbyBibGluZCB0byBzZWU=

-->
