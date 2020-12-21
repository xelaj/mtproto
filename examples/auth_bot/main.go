// Example of using Mtproto for Telegram bot
package main

import (
	"fmt"
	"os"

	"github.com/xelaj/mtproto/telegram"
)

const (
	// from https://my.telegram.org/apps
	TgAppID       = XXXXX                // integer value from "App api_id" field
	TgAppHash     = "XXXXXXXXXXXX"       // string value from "App api_hash" field
	TgTestServer  = "149.154.167.40:443" // string value from "Test configuration" field
	TgProdServer  = "149.154.167.50:443" // string value from "Production configuration" field

	// from https://t.me/BotFather
	TgBotToken    = "XXXXX"  // bot token from BotFather
	TgBotUserName = "YourBotUserName" // username of the bot
)

func main() {
	client, err := telegram.NewClient(telegram.ClientConfig{
		// current dir must be writable
		// file 'session.json' will be created here
		SessionFile: "./session.json",
		// file 'keys.pem' must contain text from "Public keys" field
		// from https://my.telegram.org/apps
		PublicKeysFile: "./keys.pem",
		// we need to use production Telegram API server
		// because test server don't know about our bot
		ServerHost: TgProdServer,
		AppID: TgAppID,
		AppHash: TgAppHash,
	})
	if err != nil {
		fmt.Println("NewClient error:", err.Error())
		os.Exit(1)
	}

	// Trying to auth as bot with our bot token
	_, err = client.AuthImportBotAuthorization(&telegram.AuthImportBotAuthorizationParams{
		Flags: 1, // reserved, must be set (not 0)
		ApiId: TgAppID,
		ApiHash: TgAppHash,
		BotAuthToken: TgBotToken,
	})
	if err != nil {
		fmt.Println("ImportBotAuthorization error:", err.Error())
		os.Exit(1)
	}

	// Request info about username of our bot
	uname, err := client.ContactsResolveUsername(&telegram.ContactsResolveUsernameParams{Username: TgBotUserName})
	if err != nil {
		fmt.Println("ResolveUsername error:", err.Error())
		os.Exit(1)
	}

	chatsCount := len(uname.Chats)
	// No chats for requested username of our bot
	if chatsCount > 0 {
		fmt.Println("Chats number:", chatsCount)
		os.Exit(1)
	}

	usersCount := len(uname.Users)
	// Users vector must contain single item with information about our bot
	if usersCount != 1 {
		fmt.Println("Users number:", usersCount)
		os.Exit(1)
	}

	user := uname.Users[0].(*telegram.UserObj)

	// dump our bot info
	fmt.Println("\nSelf ->", user.Self)
	fmt.Println("Username ->", user.Username)
	fmt.Println("FirstName ->", user.FirstName)
	fmt.Println("LastName ->", user.LastName)
	fmt.Println("Id ->", user.Id)
	fmt.Println("Bot ->", user.Bot)
	fmt.Println("Verified ->", user.Verified)
	fmt.Println("Restricted ->", user.Restricted)
	fmt.Println("Support ->", user.Support)
	fmt.Println("Scam ->", user.Scam)
	fmt.Println("BotInfoVersion ->", user.BotInfoVersion)

	// fmt.Println("Contact ->", user.Contact)
	// fmt.Println("MutualContact ->", user.MutualContact)
	// fmt.Println("Deleted ->", user.Deleted)
	// fmt.Println("BotChatHistory ->", user.BotChatHistory)
	// fmt.Println("BotNochats ->", user.BotNochats)
	// fmt.Println("Min ->", user.Min)
	// fmt.Println("BotInlineGeo ->", user.BotInlineGeo)
	// fmt.Println("ApplyMinPhoto ->", user.ApplyMinPhoto)
	// fmt.Println("AccessHash ->", user.AccessHash)
	// fmt.Println("Phone ->", user.Phone)
	// fmt.Println("Photo ->", user.Photo)
	// fmt.Println("Status ->", user.Status)
	// fmt.Println("RestrictionReason ->", user.RestrictionReason)
	// fmt.Println("BotInlinePlaceholder ->", user.BotInlinePlaceholder)
	// fmt.Println("LangCode ->", user.LangCode)
}
