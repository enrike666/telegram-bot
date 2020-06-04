package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type UpdatesGetter interface {
	//что-то должен получать на вход наверное
	GetUpdates() ([]tgbotapi.Update, error)
}

type TelegramUpdatesGetter struct {
	telegramBot *tgbotapi.BotAPI
	udateConfig tgbotapi.UpdateConfig
}

func NewTelegramUpdatesGetter(telegramBot *tgbotapi.BotAPI, udateConfig tgbotapi.UpdateConfig) *TelegramUpdatesGetter {
	return &TelegramUpdatesGetter{telegramBot: telegramBot, udateConfig: udateConfig}
}

func (t *TelegramUpdatesGetter) GetUpdates() ([]tgbotapi.Update, error) {

	updates, err := t.telegramBot.GetUpdates(t.udateConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var newUpdates []tgbotapi.Update

	for _, update := range updates {
		if update.UpdateID >= t.udateConfig.Offset {
			t.udateConfig.Offset = update.UpdateID + 1
			newUpdates = append(newUpdates, update)
		}
	}

	return newUpdates, nil
}
