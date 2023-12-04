package pkg

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TelegramBotToken string // структура для хранения телеграм бота
}

func GetToken(fileName string) string {
	file, _ := os.Open(fileName)          // открываем файл
	defer file.Close()                    // закрываем файл при выходе из функции
	decoder := json.NewDecoder(file)      // создаем декодировщик для файла
	configuration := Config{}             // создаем экземпляр структуры для хранения токена бота
	err := decoder.Decode(&configuration) // сканируем данные из файла в нашу структуру
	if err != nil {
		log.Panic(err)
	}
	return configuration.TelegramBotToken // возвращаем токен бота
}
