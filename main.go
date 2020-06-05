package main

import (
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	contectTypeHeader              = "Content-Type"
	contectTypeXWWWWFormUrlencoded = "application/x-www-form-urlencoded"
	contentLengthHeader            = "Content-Length"
	proverkachekaApiURL            = "https://proverkacheka.com/check/get"
)

type Config struct {
	TelegramBotToken string
}

func main() {
	botToken, err := getToken()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	getter := NewTelegramUpdatesGetter(bot, updateConfig)
	getFullChecker := NewTelegramGetFullChecker(proverkachekaApiURL, &http.Client{})
	handler := NewTelegramUpdateHandler(bot, getFullChecker)

	telegramWorker := NewTelegramWorker(getter, handler)

	for {
		time.Sleep(time.Second)

		err = telegramWorker.Work()
		if err != nil {
			break
		}
	}

	log.Println(err)
}

func ParseCheckInfo(checkQueryString string) (*CheckInfoFromBot, error) {
	var checkInfo CheckInfoFromBot

	values, err := url.ParseQuery(checkQueryString)
	if err != nil {
		return nil, err
	}

	var decoder = schema.NewDecoder()

	err = decoder.Decode(&checkInfo, values)
	if err != nil {
		return nil, err
	}

	return &checkInfo, nil
}

func getToken() (string, error) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return configuration.TelegramBotToken, nil
}
