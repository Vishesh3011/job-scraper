package workers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
	"log/slog"
	"strings"
	"sync"
)

type telegramReceiverWorker struct {
	bot          *tgbotapi.BotAPI
	appCtx       context.Context
	svc          *service.Service
	mailbox      actor.MailboxSender[models.BotMsg]
	logger       *slog.Logger
	updates      tgbotapi.UpdatesChannel
	userSessions map[int64]*models.UserTelegramSession
	mu           sync.RWMutex
	client.Client
}

func (w *telegramReceiverWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	case update := <-w.updates:
		if err := w.handleSingleUpdate(update); err != nil {
			w.logger.Error(utils.PrepareLogMsg(fmt.Sprintf("Error handling update: %v", err)))
			if err := w.mailbox.Send(ctx, models.BotMsg{
				ChatId: update.Message.Chat.ID,
				Text:   string(types.PromptErrorProcessingRequest),
			}); err != nil {
				w.logger.Error(utils.PrepareLogMsg(fmt.Sprintf("Error sending update: %v", err)))
			}
		}
		return actor.WorkerContinue
	}
}

func (w *telegramReceiverWorker) handleSingleUpdate(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	chatId := update.Message.Chat.ID
	msgTxt := update.Message.Text

	if msgTxt == "/start" {
		w.mu.Lock()
		w.userSessions[chatId] = &models.UserTelegramSession{TelegramState: types.AWAIT_USER_NAME}
		w.mu.Unlock()

		return w.mailbox.Send(w.appCtx, models.BotMsg{
			ChatId: chatId,
			Text:   string(types.PromptWelcome),
		})
	}

	w.mu.Lock()
	session, exists := w.userSessions[chatId]
	if !exists {
		w.userSessions[chatId] = &models.UserTelegramSession{TelegramState: types.AWAIT_USER_NAME}
		session = w.userSessions[chatId]
	}
	w.mu.Unlock()

	if !exists {
		return w.mailbox.Send(w.appCtx, models.BotMsg{
			ChatId: chatId,
			Text:   string(types.PromptEnterName),
		})
	}

	switch session.TelegramState {
	case types.AWAIT_USER_NAME:
		return w.handleAwaitUserName(session, chatId, msgTxt)
	case types.AWAIT_JOB_ROLES:
		return w.handleAwaitJobRoles(session, chatId, msgTxt)
	case types.AWAIT_LOCATION:
		return w.handleAwaitLocations(session, chatId, msgTxt)
	case types.AWAIT_COOKIE:
		return w.handleAwaitCookie(session, chatId, msgTxt)
	case types.AWAIT_CSRF_TOKEN:
		return w.handleAwaitCsrfToken(session, chatId, msgTxt)
	case types.AWAIT_EMAIL_NOTIFY:
		return w.handleAwaitEmailNotify(session, chatId, msgTxt)
	case types.AWAIT_EMAIL:
		return w.handleAwaitEmail(session, chatId, msgTxt)
	case types.AWAIT_UPDATE_DETAILS:
		return w.handleAwaitUpdateDetails(session, chatId, msgTxt)
	}
	return nil
}

func (w *telegramReceiverWorker) handleAwaitUserName(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Name = msgTxt
	w.logger.Info("User name received", slog.String("username", session.Name), slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_JOB_ROLES
	return w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Text:   string(types.PromptEnterJobRoles),
	})
}

func (w *telegramReceiverWorker) handleAwaitJobRoles(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Keywords = strings.Split(msgTxt, ",")
	w.logger.Info("User job roles received", slog.String("keywords", msgTxt), slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_LOCATION
	return w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Text:   string(types.PromptEnterJobLocations),
	})
}

func (w *telegramReceiverWorker) handleAwaitLocations(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	locations := strings.Split(msgTxt, ",")
	w.logger.Info("User locations received", slog.String("locations", msgTxt), slog.Int64("chat_id", chatId))

	var geoIds []string
	for _, l := range locations {
		location := w.svc.Location().FetchGeoIdBasedOnLocation(l)
		if location == "" {
			w.logger.Info("User input for location is invalid", slog.String("location", l), slog.Int64("chat_id", chatId))
			session.TelegramState = types.AWAIT_LOCATION
			return w.mailbox.Send(w.appCtx, models.BotMsg{
				ChatId: chatId,
				Text:   string(types.PromptEnterJobLocationsAgain),
			})
		}
		geoIds = append(geoIds, location)
	}

	session.GeoIds = geoIds
	session.TelegramState = types.AWAIT_COOKIE
	return w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Text:   string(types.PromptEnterLinkedInCookie),
	})
}

func (w *telegramReceiverWorker) handleAwaitCookie(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Cookie = msgTxt
	w.logger.Info("User cookie received", slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_CSRF_TOKEN
	return w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Text:   string(types.PromptEnterLinkedInCSRFToken),
	})
}

func (w *telegramReceiverWorker) handleAwaitCsrfToken(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.CsrfToken = msgTxt
	w.logger.Info("User csrf token received", slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_EMAIL_NOTIFY
	return w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Text:   string(types.PromptAskEmailReportPreference),
	})
}

func (w *telegramReceiverWorker) handleAwaitEmailNotify(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	w.logger.Info("User email notification preference received", slog.String("preference", msgTxt), slog.Int64("chat_id", chatId))

	if msgTxt == "y" || msgTxt == "Y" {
		session.TelegramState = types.AWAIT_EMAIL
		return w.mailbox.Send(w.appCtx, models.BotMsg{
			ChatId: chatId,
			Text:   string(types.PromptEnterEmail),
		})
	} else {
		w.logger.Info("User opted out of email notifications", slog.String("preference", msgTxt), slog.Int64("chat_id", chatId))

		userService := w.svc.User()
		user, err := userService.CreateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, nil, session.Keywords, session.GeoIds))
		if err != nil {
			return err
		}

		session.User = user
		session.TelegramState = types.SEND_REPORT
		return w.handleSendReport(session, chatId)
	}
}

func (w *telegramReceiverWorker) handleAwaitEmail(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Email = &msgTxt

	userService := w.svc.User()
	user, err := userService.GetUserByEmail(*session.Email)
	if err != nil && !errors.Is(err, types.ErrRecordNotFound) {
		return err
	}

	if user == nil && errors.Is(err, types.ErrRecordNotFound) {
		user, err := userService.CreateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, session.Email, session.Keywords, session.GeoIds))
		if err != nil {
			return err
		}
		session.User = user
		session.TelegramState = types.SEND_REPORT
		w.logger.Info("New user registered", slog.String("email", *session.Email), slog.Int64("chat_id", chatId))

		if err := w.mailbox.Send(w.appCtx, models.BotMsg{
			ChatId: chatId,
			Text:   string(types.PromptRegistrationSuccess),
		}); err != nil {
			return err
		}

		return w.handleSendReport(session, chatId)
	} else {
		if err := w.mailbox.Send(w.appCtx, models.BotMsg{
			ChatId: chatId,
			Text:   string(types.PromptAccountExistsUpdateRequest),
		}); err != nil {
			return err
		}
		session.TelegramState = types.AWAIT_UPDATE_DETAILS
		session.User = user
	}
	return nil
}

func (w *telegramReceiverWorker) handleAwaitUpdateDetails(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	if msgTxt == "y" || msgTxt == "Y" {
		updatedUser, err := w.svc.User().UpdateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, session.Email, session.Keywords, session.GeoIds))
		if err != nil {
			return err
		}
		session.User = updatedUser

		if err := w.mailbox.Send(w.appCtx, models.BotMsg{
			ChatId: chatId,
			Text:   string(types.PromptPreferencesUpdated),
		}); err != nil {
			return err
		}

		session.TelegramState = types.SEND_REPORT
		w.logger.Info("Existing user updated", slog.String("email", *session.Email), slog.Int64("chat_id", chatId))
	} else {
		session.TelegramState = types.SEND_REPORT
	}
	return w.handleSendReport(session, chatId)
}

func (w *telegramReceiverWorker) handleSendReport(session *models.UserTelegramSession, chatId int64) error {
	accumulatorService := w.svc.Accumulator()
	jobs, err := accumulatorService.FetchJobs(session.User)
	if err != nil {
		return err
	}
	w.logger.Info("Jobs fetched", slog.Int("count", len(jobs)), slog.Int64("chat_id", chatId))

	reportService := w.svc.Report()
	file, err := reportService.GenerateReport(jobs, session.User.Name)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return err
	}

	if err := w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Text:   string(types.PromptReportGenerated),
	}); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_report.xlsx", session.User.Name)
	w.logger.Info("Report generated", slog.String("file", fileName), slog.Int64("chat_id", chatId))

	doc := tgbotapi.NewDocumentUpload(chatId, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: buf.Bytes(),
	})

	if err := w.mailbox.Send(w.appCtx, models.BotMsg{
		ChatId: chatId,
		Doc:    &doc,
	}); err != nil {
		return err
	}

	w.logger.Info("Report sent to user", slog.String("file", fileName), slog.Int64("chat_id", chatId))

	if session.User.Email != nil {
		w.logger.Info("Sending email to user", slog.String("email", *session.User.Email), slog.Int64("chat_id", chatId))
		if err := w.GoMailClient().SendEmail(session.User, file, len(jobs), fileName); err != nil {
			return err
		}
		w.logger.Info("Email sent to user", slog.String("email", *session.User.Email), slog.Int64("chat_id", chatId))
	}

	session.TelegramState = types.FINISHED
	return nil
}
