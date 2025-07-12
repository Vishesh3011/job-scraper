package workers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/utils"
	"log"
	"math/rand"
	"time"
)

type worker struct {
	application.Application
}

func NewWorker(app application.Application) *worker {
	return &worker{app}
}

func (worker *worker) Start() {
	bot, err := worker.Clients().TelegramClient().GetTelegramBot(worker.Config().TelegramConfig().Token())
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %v", err)
	}

	mailbox := actor.NewMailbox[models.BotMsg]()

	processor := actor.New(&telegramSenderWorker{
		bot:     bot,
		logger:  worker.Logger(),
		mailbox: mailbox,
	})

	poller := actor.New(&telegramReceiverWorker{
		bot:     bot,
		logger:  worker.Logger(),
		mailbox: mailbox,
		app:     worker.Application,
	})

	cron := actor.New(&cronWorker{
		app: worker,
	})

	actors := actor.Combine(processor, poller, cron).Build()
	actors.Start()

	worker.Logger().Info(fmt.Sprintf("Worker started at %s...", time.Now().Format("2006-01-02T15:04:05 MST")))
	<-utils.WaitForTermination(worker.Cancel())

	msg := tgbotapi.NewMessage(int64(rand.Uint64()), "Worker shutting down. See you later!")
	if _, err := bot.Send(msg); err != nil {
		worker.Logger().Error(fmt.Sprintf("Failed to send shutdown message: %v", err))
	}
}
