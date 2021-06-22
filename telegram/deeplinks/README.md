# Simple FAQ about this package

## Purpose of this package

#### Q: Why does this package even exist?

Deeplinks are big part of telegram infrastructure: with them you can login via
qr code, start any bot, get post link, etc. With deep links you can even join
hidden chats! But dealing with tg:// links pretty hard, so this is why
`deeplinks` package exist. Are you required to use it? Hell no, in most of cases
you don't need resolving, parsing, or do anything with these links, so, have
some problem with them? Use this package.

#### Q: Where are all references/specs of deeplinks? How you know that
`tg://resolve` is working?

This is really PainInTheAss problem: telegram docs don't describe this feature
**AT ALL.** The first one (and the last one, unfortunately) doc page which is
describing deeplinking system is
[tiny part of bots api](https://core.telegram.org/bots#deep-linking). But did
you know, that you can set your language to english
[with `tg://setlanguage` href](tg://resolve?domain=DeepLink&post=42)? No? That's
the problem.

So, here is a few sources, where we get all these references:
[Deep Link telegram channel](tg://resolve?domain=DeepLink),
[telegram desktop](https://git.io/JtYos) source code,
[telegram android](https://git.io/JtcJf) (creepiest spaghetti code i've ever
seen, honestly), [this](https://git.io/JtcJs) and [this](https://git.io/JtcJn)
files from telegram ios client. Maybe there is more good implementations? If
yes, please add issue with feature request.

## How to use it

#### Q: Which scheme is preferrable?

It's way more better to use `tg://` instead `http(s)://` for two main reasons:
first is that MOST OF web browsers, android/ios apps, windows/macos and linux
apps. Don't believe me? Okay, just try to paste this command in your terminal:
`xdg-open 'tg://resolve?domain=xelaj_developers'` (note that `xdg-open` works
only on linux). See? If telegram desktop installed on your computer, you can
contemplate our beautiful amazing xelaj developers chat!

So, Why you don't want to use only `tg://` links? Sometimes, maybe a few editors
doesn't support it, so in that situations it's better to use https. Note that
only a few deeplinks can be converted to https.

#### Why this package is separated of mtproto/telegram?

`deeplinks` using `github.com/gorilla/schema` package for decoding queries in
links. But define this dependency to all `mtproto` package is really
overheading. We think that most of time you don't need parsing deeplinks, so if
you don't want to, you can just don't depend on `schema` package.
