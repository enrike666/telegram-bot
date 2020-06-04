package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UpdateHandler interface {
	HandleUpdates([]tgbotapi.Update) error
}

type TelegramUpdateHandler struct {
	telegramBot    *tgbotapi.BotAPI
	getFullChecker GetFullChecker
}

func NewTelegramUpdateHandler(telegramBot *tgbotapi.BotAPI, getFullChecker GetFullChecker) *TelegramUpdateHandler {
	return &TelegramUpdateHandler{telegramBot: telegramBot, getFullChecker: getFullChecker}
}

func (t *TelegramUpdateHandler) HandleUpdates(updates []tgbotapi.Update) error {
	for _, update := range updates {
		err := t.handleUpdate(update)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TelegramUpdateHandler) handleUpdate(update tgbotapi.Update) error {
	if update.Message.ReplyToMessage != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		replMessage := update.Message.ReplyToMessage
		switch replMessage.Text {
		case "введите параметры ответом на это сообщение":
			checkInfo, err := ParseCheckInfo(update.Message.Text)
			if err != nil {
				return err
			}
			check, err := t.getFullChecker.GetFullCheck(*checkInfo)
			if err != nil {
				return err
			}
			msg.Text = check.String()
		}

		_, err := t.telegramBot.Send(msg)
		if err != nil {
			return err
		}
	}

	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "check":
			msg.Text = "введите параметры ответом на это сообщение"
		default:
			msg.Text = "не понял!"
		}

		_, err := t.telegramBot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
