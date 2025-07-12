package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type BotMsg struct {
	Text   string
	ChatId int64
	Doc    *tgbotapi.DocumentConfig
}
