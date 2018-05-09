package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(update Update) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "start message")
	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func Help(update Update) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "help message")
	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func Default(update Update) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "use /help")
	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
