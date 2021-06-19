package main

// original could be found at https://github.com/umputun/feed-master/pull/37

import (
	"bytes"
	"fmt"
	"hash/maphash"
	"html/template"
	"log"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"
	"golang.org/x/net/html"

	utils "github.com/xelaj/mtproto/examples/example_utils"
)

const kb = 1024
const mb = kb * 1024

// Item for rss
type Item struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`
}

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
		// public keys file is path to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile:  publicKeys,
		AppID:           94575,                              // app id, could be find at https://my.telegram.org
		AppHash:         "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
		InitWarnChannel: true,                               // if we want to get errors, otherwise, client.Warnings will be set nil
	})
	utils.ReadWarningsToStdErr(client.Warnings)
	dry.PanicIfErr(err)

	// authorize the bot
	_, err = client.AuthImportBotAuthorization(0, 94575, "a3406de8d171bb422bb6ddf3bbd800e2", "here_goes_the_token")
	dry.PanicIfErr(errors.Wrapf(err, "error authorizing with telegram bot"))

	// get reference for the channel we'll be posting to
	chanRef, err := getChannelReference(client, "here_goes_your_public_channel_name")
	dry.PanicIfErr(errors.Wrapf(err, "error retrieving channel metadata"))

	// here you should prepare an object to send
	item := Item{
		Title:       "This is a test message to send",
		Description: "<a href=\"https://example.org/test\">This is a test link</a>",
		Link:        "https://example.org",
	}

	htmlMessage := getMessageHTML(item)
	plainMessage := getPlainMessage(htmlMessage)

	entities := getMessageFormatting(htmlMessage, plainMessage)

	// send the formatted message
	_, err = client.MessagesSendMessage(&telegram.MessagesSendMessageParams{
		NoWebpage: true,
		Peer:      chanRef,
		Message:   plainMessage,
		Entities:  entities,
		RandomID:  getInt64Hash(plainMessage),
	})
	dry.PanicIfErr(err)
}

// getChannelReference returns telegram channel metadata reference which
// is enough to send messages to that channel using the telegram API
func getChannelReference(client *telegram.Client, channelID string) (telegram.InputPeer, error) {
	channel, err := client.ContactsResolveUsername(channelID)
	if err != nil {
		return nil, err
	}
	return &telegram.InputPeerChannel{
		ChannelID:  channel.Chats[0].(*telegram.Channel).ID,
		AccessHash: channel.Chats[0].(*telegram.Channel).AccessHash,
	}, nil
}

// getInt64Hash generates int64 hash from the provided string, returns 0 in case of error
func getInt64Hash(s string) int64 {
	hash := maphash.Hash{}
	_, _ = hash.Write([]byte(s))
	return int64(hash.Sum64())
}

// https://core.telegram.org/api/entities
// currently only links are supported, but it's possible to parse all listed entities
func tagLinkOnlySupport(htmlText string) string {
	p := bluemonday.NewPolicy()
	p.AllowAttrs("href").OnElements("a")
	return p.Sanitize(htmlText)
}

// getPlainMessage strips provided HTML to the bare text
func getPlainMessage(htmlText string) string {
	p := bluemonday.NewPolicy()
	return p.Sanitize(htmlText)
}

// getMessageHTML generates HTML message from provided feed.Item
func getMessageHTML(item Item) string {
	title := strings.TrimSpace(item.Title)

	description := tagLinkOnlySupport(string(item.Description))
	description = strings.TrimSpace(description)

	messageHTML := fmt.Sprintf("<a href=\"%s\">%s</a>\n\n%s", item.Link, title, description)

	return messageHTML
}

// getMessageFormatting gets links from HTML text and maps them to same text in plain format using MessageEntity
func getMessageFormatting(htmlMessage, plainMessage string) []telegram.MessageEntity {
	doc, err := html.Parse(bytes.NewBufferString(htmlMessage))
	if err != nil {
		log.Printf("[WARN] can't parse HTML message: %v", err)
		return nil
	}

	b, err := getBody(doc)
	if err != nil {
		log.Printf("[WARN] problem finding HTML message body: %v", err)
		return nil
	}

	n := b.FirstChild
	var entities []telegram.MessageEntity
	for n != nil {
		if n.Data == "a" {
			url := ""
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url = attr.Val
				}
			}
			if n.FirstChild == nil || n.FirstChild != n.LastChild {
				log.Printf("[WARN] problem parsing a href=%s, can't retrieve link text", url)
				n = n.NextSibling
				continue
			}
			aText := strings.TrimSpace(n.FirstChild.Data)
			offsetIndexUTF8 := strings.Index(plainMessage, aText)
			offsetIndexUTF16 := len(utf16.Encode([]rune(plainMessage[:offsetIndexUTF8])))
			lengthUTF16 := len(utf16.Encode([]rune(aText)))
			entities = append(entities, &telegram.MessageEntityTextURL{
				Offset: int32(offsetIndexUTF16),
				Length: int32(lengthUTF16),
				URL:    url,
			})
		}
		n = n.NextSibling
	}

	return entities
}

// getBody returns provided document <body> node if found
func getBody(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("missing <body> in the node tree")
}
