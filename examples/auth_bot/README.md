# Example of using Mtproto for Telegram bot

Mtproto lib can be used with Telegram bot, that can access to full Telegram API, not just the simplified Telegram Bot API.

First, you need to create your telegram bot, like you did this in any another case. Don't know how? Read [this official guide](https://core.telegram.org/bots)

Next we must [register our app](https://my.telegram.org/apps) as usual and obtain all data that we need to run any other examples in this package.

Then fill this code in:
```go
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
```

Okay, you've got bot token, appID, appHash, saved your app public keys, added test and prod ip addresses. What next?

just run in this folder:

```
go run main.go
```

And done! That's full example of how to run bot api under mtproto! Easy-peasy!