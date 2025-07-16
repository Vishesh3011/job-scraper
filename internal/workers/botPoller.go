package workers

import (
	"bytes"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/types"
	"log/slog"
	"strings"
	"sync"
)

type telegramReceiverWorker struct {
	bot          *tgbotapi.BotAPI
	app          application.Application
	mailbox      actor.MailboxSender[models.BotMsg]
	logger       *slog.Logger
	updates      tgbotapi.UpdatesChannel
	userSessions map[int64]*models.UserTelegramSession // Move sessions to worker level
	mu           sync.RWMutex
}

func (w *telegramReceiverWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	case update := <-w.updates:
		if err := w.handleSingleUpdate(update); err != nil {
			w.logger.Error(fmt.Sprintf("Error handling update: %v", err))
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

		return w.mailbox.Send(w.app.Context(), models.BotMsg{
			ChatId: chatId,
			Text:   "Welcome to the JobScraper Telegram Bot! Please enter your name",
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
		return w.mailbox.Send(w.app.Context(), models.BotMsg{
			ChatId: chatId,
			Text:   "Please enter your name",
		})
	}

	switch session.TelegramState {
	case types.AWAIT_USER_NAME:
		return w.handleAwaitUserName(session, chatId, msgTxt)
	case types.AWAIT_JOB_ROLES:
		return w.handleAwaitJobRoles(session, chatId, msgTxt)
	case types.AWAIT_GEO_IDS:
		return w.handleAwaitGeoIds(session, chatId, msgTxt)
	case types.AWAIT_COOKIE:
		return w.handleAwaitCookie(session, chatId, msgTxt)
	case types.AWAIT_CSRF_TOKEN:
		return w.handleAwaitCsrfToken(session, chatId, msgTxt)
	case types.AWAIT_EMAIL_NOTIFY:
		return w.handleAwaitEmailNotify(session, chatId, msgTxt)
	case types.AWAIT_EMAIL:
		return w.handleAwaitEmail(session, chatId, msgTxt)
	}
	return nil
}

func (w *telegramReceiverWorker) handleAwaitUserName(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Name = msgTxt
	w.logger.Info("User name received", slog.String("username", session.Name), slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_JOB_ROLES
	return w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "Please enter your interested job roles (separated by commas)",
	})
}

func (w *telegramReceiverWorker) handleAwaitJobRoles(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Keywords = strings.Split(msgTxt, ",")
	w.logger.Info("User job roles received", slog.String("keywords", msgTxt), slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_GEO_IDS
	return w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "Please enter your interested job location geo-id (separated by commas)",
	})
}

func (w *telegramReceiverWorker) handleAwaitGeoIds(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Locations = strings.Split(msgTxt, ",")
	w.logger.Info("User geo-ids received", slog.String("geo_ids", msgTxt), slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_COOKIE
	return w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "Please enter your cookie from linkedin",
	})
}

func (w *telegramReceiverWorker) handleAwaitCookie(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Cookie = msgTxt
	w.logger.Info("User cookie received", slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_CSRF_TOKEN
	return w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "Please enter your csrf token from linkedin",
	})
}

func (w *telegramReceiverWorker) handleAwaitCsrfToken(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.CsrfToken = msgTxt
	w.logger.Info("User csrf token received", slog.Int64("chat_id", chatId))

	session.TelegramState = types.AWAIT_EMAIL_NOTIFY
	return w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "Are you interested in daily email report for jobs (y/n) ?",
	})
}

func (w *telegramReceiverWorker) handleAwaitEmailNotify(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	w.logger.Info("User email notification preference received", slog.String("preference", msgTxt), slog.Int64("chat_id", chatId))

	if msgTxt == "y" || msgTxt == "Y" {
		session.TelegramState = types.AWAIT_EMAIL
		return w.mailbox.Send(w.app.Context(), models.BotMsg{
			ChatId: chatId,
			Text:   "Please enter your email: ",
		})
	} else {
		w.logger.Info("User opted out of email notifications", slog.String("preference", msgTxt), slog.Int64("chat_id", chatId))

		userService := service.NewService(w.app).User()
		user, err := userService.CreateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, nil, session.Keywords, session.Locations))
		if err != nil {
			return err
		}

		session.CreatedUser = user
		session.TelegramState = types.SEND_REPORT
		return w.handleSendReport(session, chatId)
	}
}

func (w *telegramReceiverWorker) handleAwaitEmail(session *models.UserTelegramSession, chatId int64, msgTxt string) error {
	session.Email = &msgTxt

	userService := service.NewService(w.app).User() // Adjust based on your service structure
	user, err := userService.GetUserByEmail(*session.Email)
	if err != nil && !errors.Is(err, types.ErrRecordNotFound) {
		return err
	}

	if user == nil && errors.Is(err, types.ErrRecordNotFound) {
		user, err := userService.CreateUser(models.NewUserInput(session.Name, session.Cookie, session.CsrfToken, session.Email, session.Keywords, session.Locations))
		if err != nil {
			return err
		}
		session.CreatedUser = user
		session.TelegramState = types.SEND_REPORT

		w.logger.Info("New user registered", slog.String("email", *session.Email), slog.Int64("chat_id", chatId))
	} else {
		if err := w.mailbox.Send(w.app.Context(), models.BotMsg{
			ChatId: chatId,
			Text:   "Your account already exists! Updating your details....",
		}); err != nil {
			return err
		}

		session.CreatedUser = user
		if _, err := userService.UpdateUser(models.NewUserInput(user.Name, user.Cookie, user.CsrfToken, user.Email, user.Keywords, user.Locations)); err != nil {
			return err
		}

		if err := w.mailbox.Send(w.app.Context(), models.BotMsg{
			ChatId: chatId,
			Text:   "Your preferences are updated successfully to our service!",
		}); err != nil {
			return err
		}

		session.TelegramState = types.SEND_REPORT
		w.logger.Info("Existing user updated", slog.String("email", *session.Email), slog.Int64("chat_id", chatId))
	}
	if err := w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "You are registered successfully to our service! Sending report to you and your email...",
	}); err != nil {
		return err
	}
	return w.handleSendReport(session, chatId)
}

func (w *telegramReceiverWorker) handleSendReport(session *models.UserTelegramSession, chatId int64) error {
	accumulatorService := service.NewService(w.app).Accumulator()
	jobs, err := accumulatorService.FetchJobs(session.CreatedUser)
	if err != nil {
		return err
	}
	w.logger.Info("Jobs fetched", slog.Int("count", len(jobs)), slog.Int64("chat_id", chatId))

	reportService := service.NewService(w.app).Report() // Adjust based on your service structure
	file, err := reportService.GenerateReport(jobs, session.CreatedUser.Name)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return err
	}

	if err := w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Text:   "Report generated successfully! Sending to you and your email...",
	}); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_report.xlsx", session.CreatedUser.Name)
	w.logger.Info("Report generated", slog.String("file", fileName), slog.Int64("chat_id", chatId))

	doc := tgbotapi.NewDocumentUpload(chatId, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: buf.Bytes(),
	})

	if err := w.mailbox.Send(w.app.Context(), models.BotMsg{
		ChatId: chatId,
		Doc:    &doc,
	}); err != nil {
		return err
	}

	w.logger.Info("Report sent to user", slog.String("file", fileName), slog.Int64("chat_id", chatId))

	if session.CreatedUser.Email != nil {
		w.logger.Info("Sending email to user", slog.String("email", *session.CreatedUser.Email), slog.Int64("chat_id", chatId))
		if err := w.app.Clients().GoMailClient().SendEmail(session.CreatedUser, file, len(jobs)); err != nil {
			return err
		}
		w.logger.Info("Email sent to user", slog.String("email", *session.CreatedUser.Email), slog.Int64("chat_id", chatId))
	}

	session.TelegramState = types.FINISHED
	return nil
}
