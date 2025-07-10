package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/repository"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
	"log/slog"
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
	*slog.Logger
}

func newTelegramService(ctx context.Context, q *repository.Queries, config config.Config, clients client.Client, logger *slog.Logger) TelegramService {
	return telegramService{
		Context: ctx,
		Client:  clients,
		Config:  config,
		Queries: q,
		Logger:  logger,
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
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter your name: ")); err != nil {
				return err
			}
			continue
		}

		createdUser := &models.User{}
		userService := newUserService(t.Context, t.Queries, t.EncryptionKey())

		switch session.TelegramState {
		case types.AWAIT_USER_NAME:
			session.Name = msgTxt
			session.TelegramState = types.AWAIT_JOB_ROLES
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter your interested job roles (seperated by commas)")); err != nil {
				return err
			}
			t.Info("User name received", slog.String("username", session.Name), slog.Int64("chat_id", chatId))
		case types.AWAIT_JOB_ROLES:
			session.Keywords = strings.Split(msgTxt, ",")
			session.TelegramState = types.AWAIT_GEO_IDS
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Enter your interested job location geo-id (seperated by commas)")); err != nil {
				return err
			}
			t.Info("User job roles received", slog.String("keywords", msgTxt), slog.Int64("chat_id", chatId))
		case types.AWAIT_GEO_IDS:
			session.Locations = strings.Split(msgTxt, ",")
			session.TelegramState = types.AWAIT_EMAIL_NOTIFY
			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Are you interested in daily email report for jobs (y/n) ?")); err != nil {
				return err
			}
			t.Info("User geo-ids received", slog.String("geo_ids", msgTxt), slog.Int64("chat_id", chatId))
		case types.AWAIT_EMAIL_NOTIFY:
			if msgTxt == "y" || msgTxt == "Y" {
				session.TelegramState = types.AWAIT_EMAIL
				if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Please enter your email: ")); err != nil {
					return err
				}
				t.Info("User email notification preference received", slog.String("preference", msgTxt), slog.Int64("chat_id", chatId))
			} else {
				user, err := userService.CreateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, nil, session.Keywords, session.Locations))
				if err != nil {
					return err
				}
				createdUser = user

				session.TelegramState = types.SEND_REPORT
				t.Info("User opted out of email notifications", slog.String("preference", msgTxt), slog.Int64("chat_id", chatId))
			}
		case types.AWAIT_EMAIL:
			session.Email = utils.ToPtr(msgTxt)
			user, err := userService.GetUserByEmail(*session.Email)
			if err != nil && !errors.Is(err, types.ErrRecordNotFound) {
				return err
			}
			if user == nil && errors.Is(err, types.ErrRecordNotFound) {
				user, err := userService.CreateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, session.Email, session.Keywords, session.Locations))
				if err != nil {
					return err
				}
				createdUser = user

				if _, err := bot.Send(tgbotapi.NewMessage(chatId, "You are registered successfully to our service! Sending report to you and your email...")); err != nil {
					return err
				}
				session.TelegramState = types.SEND_REPORT
				t.Info("New user registered", slog.String("email", *session.Email), slog.Int64("chat_id", chatId))
			} else {
				if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Your account already exists! Updating your details....")); err != nil {
					return err
				}
				user, err := userService.GetUserByEmail(*session.Email)
				createdUser = user
				if err != nil {
					return err
				}
				if _, err := userService.UpdateUser(models.NewUserInput(user.Name, user.Cookie, user.CsrfToken, user.Email, user.Keywords, user.Locations)); err != nil {
					return err
				}
				if _, err := bot.Send(tgbotapi.NewMessage(chatId, "You are updated successfully to our service! Sending report to you and your email...")); err != nil {
					return err
				}
				session.TelegramState = types.FINISHED
				t.Info("Existing user updated", slog.String("email", *session.Email), slog.Int64("chat_id", chatId))
			}
		}

		if session.TelegramState == types.FINISHED {
			jobs, err := newAccumulatorService(t.Client, t.Config).FetchJobs(createdUser)
			if err != nil {
				return err
			}
			t.Info("Jobs fetched", slog.Int("count", len(jobs)), slog.Int64("chat_id", chatId))

			file, err := newReportService().GenerateReport(jobs, createdUser.Name)
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			if err := file.Write(&buf); err != nil {
				return err
			}

			if _, err := bot.Send(tgbotapi.NewMessage(chatId, "Report generated successfully! Sending to you and your email...")); err != nil {
				return err
			}
			fileName := fmt.Sprintf("%s_report.xlsx", createdUser.Name)
			t.Info("Report generated", slog.String("file", fileName), slog.Int64("chat_id", chatId))

			doc := tgbotapi.NewDocumentUpload(chatId, tgbotapi.FileBytes{
				Name:  fileName,
				Bytes: buf.Bytes(),
			})
			if _, err := bot.Send(doc); err != nil {
				return err
			}
			t.Info("Report sent to user", slog.String("file", fileName), slog.Int64("chat_id", chatId))
			t.Info("Sending email to user", slog.String("email", *createdUser.Email), slog.Int64("chat_id", chatId))
			if err := t.Client.GoMailClient().SendEmail(createdUser, file, len(jobs)); err != nil {
				return err
			}
			t.Info("Email sent to user", slog.String("email", *createdUser.Email), slog.Int64("chat_id", chatId))
		}
	}
	return nil
}
