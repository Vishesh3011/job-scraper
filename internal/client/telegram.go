package client

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func GetTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return bot, nil
}
