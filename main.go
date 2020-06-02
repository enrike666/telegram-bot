package main

import (
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
)

type Config struct {
	TelegramBotToken string
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		// тут можем шаманить с сообщениями и ответами
		if update.Message == nil {
			continue
		}

		if update.Message.ReplyToMessage != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			replMessage := update.Message.ReplyToMessage
			switch replMessage.Text {
			case "введите информацию ответом на это сообщение":
				info := ParseCheckInfo(update.Message.Text)
				msg.Text = info.Sum
			}
			bot.Send(msg)
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "check":
				msg.Text = "введите информацию ответом на это сообщение"
			case "status":
				msg.Text = "I'm ok."
			default:
				msg.Text = "не понял!"
			}
			bot.Send(msg)
		}

	}
}

func ParseCheckInfo(str string) CheckInfo {
	paramsMap := map[string]string{}

	params := strings.Split(str, "&")
	for _, p := range params {
		param := strings.Split(p, "=")
		paramsMap[param[0]] = param[1]
	}

	info := CheckInfo{
		DateTime: paramsMap["t"],
		Sum: paramsMap["s"],
		FD: paramsMap["i"],
		FN: paramsMap["fn"],
		FP: paramsMap["fp"],
		N: paramsMap["n"],
	}

	return info
}



func getToken()  string {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	return configuration.TelegramBotToken
}
