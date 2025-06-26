package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/types"
	"strings"
)

type TelegramService interface {
	HandleTelegramUpdates(*tgbotapi.BotAPI, *tgbotapi.UpdatesChannel) error
}

type telegramService struct{}

func newTelegramService() TelegramService {
	return telegramService{}
}

func (t telegramService) HandleTelegramUpdates(bot *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel) error {
	var userSessions = make(map[int64]*models.UserTelegramSession)
	for update := range *updates {
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID
		msgTxt := update.Message.Text
		session, exists := userSessions[chatId]

		if !exists {
			userSessions[chatId] = &models.UserTelegramSession{TelegramState: types.AWAIT_USER_NAME}
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter you name: ")); err != nil {
				return err
			}
			continue
		}

		switch session.TelegramState {
		case types.AWAIT_USER_NAME:
			session.Name = msgTxt
			session.TelegramState = types.AWAIT_JOB_ROLES
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter your interested job roles (seperated by commas)")); err != nil {
				return err
			}
		case types.AWAIT_JOB_ROLES:
			session.Keywords = strings.Split(msgTxt, ",")
			session.TelegramState = types.AWAIT_GEO_IDS
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter your interested job location geo-id (seperated by commas)")); err != nil {
				return err
			}
		case types.AWAIT_GEO_IDS:
			session.Locations = strings.Split(msgTxt, ",")
			session.TelegramState = types.AWAIT_EMAIL_NOTIFY
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Are you interested in daily email report for jobs (y/n) ?")); err != nil {
				return err
			}
		case types.AWAIT_EMAIL_NOTIFY:
			if msgTxt == "y" || msgTxt == "Y" {
				session.TelegramState = types.AWAIT_EMAIL
				if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Please enter your email: ")); err != nil {
					return err
				}
			} else {
				session.TelegramState = types.SEND_REPORT
			}
			//case types.AWAIT_EMAIL:
			//	session.Email = utils.ToPtr(msgTxt)
			//	svc := ser
			//
			//case types.AWAIT_USER_NAME:
			//	session.Name = msgTxt
			//	session.TelegramState = types.AWAIT_JOB_ROLES
			//	if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter your interested job roles (seperated by commas)")); err != nil {
			//		return err
			//	}
		}
	}
	return nil
}
