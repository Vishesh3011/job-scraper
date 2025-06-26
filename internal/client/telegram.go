package client

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type TelegramClient struct {
	token string
}

func NewTelegramClient(token string) *TelegramClient {
	return &TelegramClient{
		token: token,
	}
}

func (c *TelegramClient) GetTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return bot, nil
}
