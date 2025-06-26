package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"job-scraper.go/internal/core/application"
)

type TelegramService interface {
	HandleTelegramUpdates(channel tgbotapi.UpdatesChannel) error
}

type telegramService struct {
	application application.Application
}

func newTelegramService(application application.Application) TelegramService {
	return telegramService{
		application: application,
	}
}

func (t telegramService) HandleTelegramUpdates(channel tgbotapi.UpdatesChannel) error {
	return nil
}
