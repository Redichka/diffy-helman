package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) error {
	switch message.Command() { // конструкция, которая выберет из списка ниже введенную команду
	case "connect":
		data, err := b.connect(message, updates) // функция для образования соединения между пользователями
		if err != nil {
			return err
		}
		err = b.chat(data, updates) // функция для работы чата между пользователями
		if err != nil {
			return err
		}
	case "register": // команда регистрации пользователя
		err := b.register(message) // функция регистрации пользователя в системе
		if err != nil {
			return err
		}
	}
	return nil
}
