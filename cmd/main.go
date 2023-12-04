package main

import (
	"diffy-helman/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(pkg.GetToken("config.json")) // инициализируем бота, передавая файл в котором хранится токен
	if err != nil {
		log.Panic(err)
	}
	telegramBot := pkg.NewBot(bot)              // создаем экземпляр бота
	if err := telegramBot.Start(); err != nil { // запускаем бота с помощью метода Start()
		log.Fatal("Ошибка")
	}
}
