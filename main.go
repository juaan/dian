package main

import (
	bot "github.com/dian/bot"
	config "github.com/dian/etc"
)

func main() {
	config.InitConfig()

	bot.BotToken = config.Get().DiscordToken
	bot.Run() // call the run function of bot/bot.go
}
