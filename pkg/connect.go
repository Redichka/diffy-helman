package pkg

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

const dataBase string = "root:root@tcp(127.0.0.1:3306)/helly-hafman"

type Connect struct {
	firstPerson  int64
	secondPerson int64
}

func (b *Bot) register(message *tgbotapi.Message) error {
	db, err := sql.Open("mysql", dataBase)
	defer db.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	find, _ := db.Prepare("SELECT `username` FROM `user` WHERE `username` = ?")
	defer find.Close()
	var usname string
	err = find.QueryRow(message.Chat.UserName).Scan(&usname)
	if message.Chat.UserName == usname {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы уже зарегестрированы")
		b.bot.Send(msg)
		return nil
	}
	add, _ := db.Prepare("INSERT INTO `user` (`username`, `telegramID`) VALUES (?, ?)")
	defer add.Close()
	_, err = add.Exec(message.Chat.UserName, message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно зарегестрировались")
	b.bot.Send(msg)
	return nil
}

func (b *Bot) connect(message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) (*Connect, error) {
	connect := &Connect{firstPerson: message.Chat.ID}
	db, _ := sql.Open("mysql", dataBase)
	defer db.Close()
	stmt, _ := db.Prepare("SELECT `id`, `username` FROM `user` WHERE `username` != ?")
	defer stmt.Close()
	data, err := stmt.Query(message.Chat.UserName)
	str := "Введите номер пользователя, с которым вы хотите установить соединенние\n"
	for data.Next() {
		var id int
		var username string
		err = data.Scan(&id, &username)
		if err != nil {
			return nil, err
		}
		str += fmt.Sprintf("%d. @%s\n", id, username)
	}
	if err != nil {
		return nil, err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, str)
	b.bot.Send(msg)
	find, _ := db.Prepare("SELECT `telegramID` FROM `user` WHERE `id` = ?")
	var secondPerson int64
	for update := range updates {
		_ = find.QueryRow(update.Message.Text).Scan(&secondPerson)
		msg = tgbotapi.NewMessage(secondPerson, fmt.Sprintf("С вами хочет установить соединение @%s, вы согласны?\n1 - Да, 2 - Нет", message.Chat.UserName))
		b.bot.Send(msg)
		if update.Message.Chat.ID == secondPerson {
			if update.Message.Text == "1" {
				connect.secondPerson = secondPerson
				break
			} else {
				msg = tgbotapi.NewMessage(message.Chat.ID, "В соединении отказано")
				b.bot.Send(msg)
				break
			}
		}
	}
	msg = tgbotapi.NewMessage(connect.firstPerson, "Соединение установлено")
	b.bot.Send(msg)
	msg = tgbotapi.NewMessage(connect.secondPerson, "Соединение установлено")
	b.bot.Send(msg)
	return connect, nil
}

func (b *Bot) chat(connect *Connect, updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		var first, second int64
		if update.Message.Chat.ID == connect.firstPerson {
			first, second = connect.firstPerson, connect.secondPerson
		} else if update.Message.Chat.ID == connect.secondPerson {
			second, first = connect.firstPerson, connect.secondPerson
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "disconnect":
				{
					msg := tgbotapi.NewMessage(first, "Соединение прервано")
					b.bot.Send(msg)
					msg = tgbotapi.NewMessage(second, "Соединение прервано")
					b.bot.Send(msg)
					return nil
				}
			case "send":
				{
					msg := tgbotapi.NewMessage(second, update.Message.CommandArguments())
					b.bot.Send(msg)
					log.Println(update.Message.Chat.UserName, ":", update.Message.CommandArguments())
				}
			case "generateNumber":
				{
					numbers := strings.Split(update.Message.CommandArguments(), " ")
					first1, _ := strconv.ParseInt(numbers[0], 10, 64)
					second2, _ := strconv.ParseInt(numbers[1], 10, 64)
					third3, _ := strconv.ParseInt(numbers[2], 10, 64)
					code := generateYouNumber(first1, second2, third3)
					msg := tgbotapi.NewMessage(first, strconv.FormatInt(code, 10))
					b.bot.Send(msg)
				}
			case "getSecretNumber":
				{
					numbers := strings.Split(update.Message.CommandArguments(), " ")
					first1, _ := strconv.ParseInt(numbers[0], 10, 64)
					second2, _ := strconv.ParseInt(numbers[1], 10, 64)
					third3, _ := strconv.ParseInt(numbers[2], 10, 64)
					code := getSecretNumber(first1, second2, third3)
					msg := tgbotapi.NewMessage(first, strconv.FormatInt(code, 10))
					b.bot.Send(msg)
				}
			case "encrypt":
				{
					arguments := strings.Split(update.Message.CommandArguments(), " ")
					second2, _ := strconv.ParseInt(arguments[1], 10, 64)
					code := encrypt(arguments[0], second2)
					msg := tgbotapi.NewMessage(first, code)
					b.bot.Send(msg)
				}
			case "decrypt":
				{
					arguments := strings.Split(update.Message.CommandArguments(), " ")
					second2, _ := strconv.ParseInt(arguments[1], 10, 64)
					code := decrypt(arguments[0], second2)
					msg := tgbotapi.NewMessage(first, code)
					b.bot.Send(msg)
				}
			}
		}
	}
	return nil
}
