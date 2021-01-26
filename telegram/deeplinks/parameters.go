package deeplinks

import (
	"net/url"
	"strconv"
)

// most of concepts stolen from tdesktop app https://git.io/JtYos
// also, some more info gathered from @deeplink channel at https://t.me/DeepLink

const (
	ReservedSchema = "tg"
)

func ReservedHosts() []string {
	return []string{
		"telegram.me",
		"telegram.dog",
		"t.me",
		"tx.me", // a few official durov's projects redirects to tx.me. looks like simple link mirror
		"telesco.pe",
	}
}

// https://t.me/DeepLink/16 converting to -> tg://resolve?domain=deeplink&post=16
type ResolveParameters struct {
	Domain string `schema:"domain"`
	Start  string `schema:"start"` // looks like not working
	//StartGroup string `schema:"startgroup"` // looks like not working
	//Game       string `schema:"game"`       // looks like not working
	Post    int `schema:"post"`    // if you copy some post link, use this param with domain, example in desc
	Thread  int `schema:"thread"`  // we don't know what does it mean
	Comment int `schema:"comment"` // we don't know what does it mean
}

func (p *ResolveParameters) String() string {
	values := url.Values{}
	if p.Domain != "" {
		values.Add("domain", p.Domain)
	}
	if p.Start != "" {
		values.Add("start", p.Start)
	}
	if p.Post != 0 {
		values.Add("post", strconv.Itoa(p.Post))
	}
	if p.Thread != 0 {
		values.Add("thread", strconv.Itoa(p.Thread))
	}
	if p.Comment != 0 {
		values.Add("comment", strconv.Itoa(p.Comment))
	}
	return (&url.URL{
		Scheme:   ReservedSchema,
		Path:     "resolve",
		RawQuery: values.Encode(),
	}).String()
}

// this parameters works ONLY if resolve?domain parameter is 'telegrampassport'
// tg://resolve?domain=telegrampassport&...
type ResolvePassportParameters struct {
	ResolveParameters
	// next parameters wasn't tested yet, need to ask telegram developers of full guide how to use these params
	// BotID       string `schema:"bot_id"`
	// Scope       string `schema:"scope"`
	// PublicKey   string `schema:"public_key"`
	// CallbackURL string `schema:"callback_url"`
	// Nonce       string `schema:"nonce"`
	// Payload     string `schema:"payload"`
	// Scope       string `schema:"scope"`
}



// https://t.me/joinchat/abcdefg
// tg://join?invite=abcdefg
type JoinParameters struct {
	Invite string `schema:"invite,required"`
}

func (p *JoinParameters) String() string {
	return (&url.URL{
		Scheme: ReservedSchema,
		Path:   "join",
		RawQuery: url.Values{
			"invite": []string{p.Invite},
		}.Encode(),
	}).String()
}

// tg://addstickers?set=abcd
type AddstickersParameters struct {
	Set string `schema:"set,required"`
}

// all below pseudohosts used to share specific link, all of them have identical purpose
// tg://msg tg://share tg://msg_url
type MsgParameters struct {
	URL  string `schema:"url,required"` // link to some resource
	Text string `schema:"text"`         // text in share message
}

// tg://confirmphone?phone=88005553535&hash=hash+used+by+telegram+api
type ConfirmPhoneParameters struct {
	Phone string `schema:"phone"` // phone with + sign like +79161234567
	Hash  string `schema:"hash"`  // confirm hash which is using by account.confirmPhone method
}

// tg://passport // TODO: figure out what the hell is this
type PassportParameters struct {
	// Scope       string `schema:"scope`
	// Nonce       string `schema:"nonce`
	// Payload     string `schema:"payload`
	// BotID       string `schema:"bot_id`
	// PublicKey   string `schema:"public_key`
	// CallbackURL string `schema:"callback_url`
}

// tg://proxy and tg://socks
// server (address)
// port (port)
// user (user)
// pass (password)
// secret (secret)

// tg://user?id=1234 or MAYBE it's not an id, but an username. or not. i don't know
type UserParameters struct {
	ID string `schema:"id"` // not sure how does it works
}

// tg:// filename (?)
// filename = dc_id + _ + document_id (?)
// filename = volume_id + _ + local_id + . + jpg (?)
// filename = md5(url) + . + extension (?)
// filename = "" (?)
// filename = dc_id + _ + document_id + _ + document_version + extension (?)
//
// id (document id)
// hash (access hash)
// dc (dc id)
// size (size)
// mime (mime type)
// name (document file name)

// tg:bg
// tg://bg
//
// slug (wallpaper)
// mode (blur+motion)
// color
// bg_color
// rotation
// intensity

// tg://search_hashtag
//
// hashtag
//
// (used internally by Telegram Web/Telegram React, you can use it by editing a href)

// tg://bot_command
//
// command
// bot
//
// (used internally by Telegram Web/Telegram React, you can use it by editing a href)

// tg://unsafe_url
//
// url
//
// (used internally by Telegram Web, you can use it by editing a href)

// tg://setlanguage
//
// lang

// tg://statsrefresh
//
// (something related to getStatsURL, probably not implemented yet)

// tg://openmessage
//
// user_id
// chat_id
// message_id
//
// (used internally by Android Stock (and fork), do not use, use tg://privatepost)

// tg://privatepost
//
// channel (channelId)
// post (messageId)
// thread (messageId)
// comment (messageId)

// links to theme i hope
//tg://addtheme
//
//slug

// ton stuff. leaved here just for know, we'll never implement it here
//ton://test/test?test=test&test=test
//
//ton://<domain>/<method>?<field1>=<value1>&<field2>=. . .
//
//ton://transfer/WALLET?amount=123&text=test

// tg://login
//
// token
// code

// da fuck is this???
// tg://settings
//
// themes
// devices
// folders
// language
// change_number

// tg://call
//
// format
// name
// phone

// REALLY specific link, used for login in telegram desktop, link works only on mobile apps
// tg://scanqr

// add contact, u no
// tg://addcontact
//
// name
// phone

// tg://search
// qyery

// Next stuff working in http(s) links but generally all of them covered by tg://, leaved here just to know
//
//
// this is really specific link, works on ios
// https://t.me/@id1234
//
// joinchat/
//
// addstickers/
//
// addtheme/
//
// iv/
//   url
//   rhash
//
// msg/
// share/
// share/url
//   url
//   text
// (Only android)
//
// confirmphone
//   phone
//   hash
//
// start
//
// startgroup
//
// game
//
// socks
// proxy
//   server (address)
//   port (port)
//   user (user)
//   pass (password)
//   secret (secret)
//
// setlanguage/
//   (12char max)
//
// bg
//   slug
//   mode
//   intensity
//   bg_color
//
// c/
//  (/chatid/messageid/ t.me/tgbeta/3539)
//   threadId
//   comment
//
// s/
//  (channel username/messageid)
//  q (search query)
//
// ?comment=
