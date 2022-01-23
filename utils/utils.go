package utils

import (
	"encoding/json"
	"log"
	"os"
)

var logFn = log.Panic

func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

type EnvConfig struct {
	TelegramToken  string `json:"telegramToken"`
	TelegramChatID int    `json:"telegramChatId"`
}

func LoadConfig() (EnvConfig, error) {
	var config EnvConfig
	file, err := os.Open("env.json")
	HandleErr(err)
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	HandleErr(err)

	return config, err
}
