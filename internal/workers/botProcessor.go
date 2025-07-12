package workers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/models"
	"log/slog"
)

type telegramSenderWorker struct {
	bot     *tgbotapi.BotAPI
	mailbox actor.MailboxReceiver[models.BotMsg]
	logger  *slog.Logger
}

func (w *telegramSenderWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	case msg := <-w.mailbox.ReceiveC():
		if msg.Doc != nil {
			_, err := w.bot.Send(msg.Doc)
			if err != nil {
				w.logger.Error("failed to send message", slog.String("error", err.Error()))
			}
		} else {
			_, err := w.bot.Send(tgbotapi.NewMessage(msg.ChatId, msg.Text))
			if err != nil {
				w.logger.Error("failed to send message", slog.String("error", err.Error()))
			}
		}
		return actor.WorkerContinue
	}
}
