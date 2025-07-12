package workers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"log/slog"
)

type telegramReceiverWorker struct {
	bot     *tgbotapi.BotAPI
	app     application.Application
	mailbox actor.MailboxSender[models.BotMsg]
	logger  *slog.Logger
}

func (w *telegramReceiverWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	default:
		updates, err := w.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: w.app.Config().TelegramConfig().Timeout()})
		if err != nil {
			w.logger.Error(fmt.Sprintf("Error getting updates from Telegram: %v", err))
			return actor.WorkerEnd
		}

		go func() {
			if err := service.NewService(w.app).Telegram().HandleTelegramUpdates(w.bot, &updates); err != nil {
				w.logger.Error(fmt.Sprintf("Error handling Telegram updates: %v", err))
			}
		}()

		<-ctx.Done()
		return actor.WorkerContinue
	}
}
