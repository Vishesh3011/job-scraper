package service

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/repository"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
	"strings"
)

type TelegramService interface {
	HandleTelegramUpdates(*tgbotapi.BotAPI, *tgbotapi.UpdatesChannel) error
}

type telegramService struct {
	context.Context
	client.Client
	config.Config
	*repository.Queries
}

func newTelegramService(ctx context.Context, q *repository.Queries, config config.Config, clients client.Client) TelegramService {
	return telegramService{
		Context: ctx,
		Client:  clients,
		Config:  config,
		Queries: q,
	}
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

		ui := &models.UserInput{}
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
		case types.AWAIT_EMAIL:
			session.Email = utils.ToPtr(msgTxt)
			userService := newUserService(t.Context, t.Queries)
			user, err := userService.GetUserByEmail(*session.Email)
			if err != nil && !errors.Is(err, types.ErrRecordNotFound) {
				return err
			}
			if user == nil {
				ui = models.NewUserInput(session.Name, session.Email, session.Cookie, session.CsrfToken, session.Keywords, session.Locations)
				if _, err := userService.CreateUser(ui); err != nil {
					return err
				}
				session.TelegramState = types.SEND_REPORT
			}
		case types.SEND_REPORT:
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "You are registered successfully to our service! Sending report to you and your email...")); err != nil {
				return err
			}

			jobs, err := newAccumulatorService(t.Client, t.Config).FetchJobs(ui)
			if err != nil {
				return err
			}

			f, err := newReportService().GenerateReport(jobs, ui.Name)
			if err != nil {
				return err
			}

			doc := tgbotapi.NewDocumentUpload(chatId, f)
			if _, err := bot.Send(doc); err != nil {
				return err
			}

			if err := t.Client.GoMailClient().SendEmail(ui, f, len(jobs)); err != nil {
				return err
			}

			session.TelegramState = types.FINISHED
		case types.FINISHED:
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "You are done with all the steps. Please enter /start to begin again.")); err != nil {
				return err
			}
		}
	}
	return nil
}
