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
	case "help":
		b.help(message)
	}
	return nil
}

func (b *Bot) help(message *tgbotapi.Message) {
	str := "/help - вывод всех команд\n/register - регистрация в Базе Данных\n/connect - создание чата с другими пользователями\nКоманды внутри соединения:\n/diffyHellmanCalculation [степень] [число] [модуль] - подсчет числа в степени по модулю для протокола Диффи-Хеллмана\n/encrypt [ключ] [сообщение] - шифрование Цезарем\n/decrypt [ключ] [сообщение] - дешифровка Цезаря\n/disconnect - отключение от соединения\n/generateRandomNumber [первое число] [степень первого числа] [второе число] [степень второго числа] - генерирует любое число в диапазоне от первого до второго"
	msg := tgbotapi.NewMessage(message.Chat.ID, str)
	b.bot.Send(msg)
}
