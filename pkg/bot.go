package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI // структура бота
}

func NewBot(bot *tgbotapi.BotAPI) *Bot { // Функция, которая создает бота
	return &Bot{bot: bot}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account Diffy-Helman") // Лог, который показывает что бот запущен
	updates, err := b.initUpdatesChannel()           // Создание канала, в который будут посылаться все обновления с Телеграм API
	if err != nil {
		return err
	}
	err = b.handleUpdates(updates) // отправляем канал с обновлениями в обработчик этих самых обновлений
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates { // запускаем цикл, который будет работать когда происходит запрос в телеграме
		if update.CallbackQuery != nil { // проверяем не является ли сообщение нажатием на кнопку
			continue // так как у нас не предусмотрена логика на такой случай, то просто пропускаем
		}
		if update.Message.IsCommand() { // проверка на то является ли обновление командой
			err := b.handleCommand(update.Message, updates) // передаем обновление в обработчик команд
			if err != nil {
				return err
			}
		}
		continue // в случае если это просто сообщение без команды то пропускаем
	}
	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)          // Создание канала
	u.Timeout = 60                      // выставляем кд проверки поступления запросов
	return b.bot.GetUpdatesChan(u), nil // возвращаем готовый канал
}
