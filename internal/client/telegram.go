package client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramClient struct {
	token string
}

func newTelegramClient(token string) *telegramClient {
	return &telegramClient{
		token: token,
	}
}

func (c *telegramClient) GetTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false
	return bot, nil
}
