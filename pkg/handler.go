package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) error {
	switch message.Command() {
	case "connect":
		data, err := b.connect(message, updates)
		if err != nil {
			return err
		}
		b.chat(data, updates)
	case "register":
		err := b.register(message)
		if err != nil {
			return err
		}
	}
	return nil
}
