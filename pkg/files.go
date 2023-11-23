package pkg

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TelegramBotToken string
}

func GetToken(fileName string) string {
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	return configuration.TelegramBotToken
}
