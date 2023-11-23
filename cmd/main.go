package main

import (
	"diffy-helman/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(pkg.GetToken("config.json"))
	if err != nil {
		log.Panic(err)
	}
	telegramBot := pkg.NewBot(bot)
	if err := telegramBot.Start(); err != nil {
		log.Fatal("Ошибка")
	}
}
