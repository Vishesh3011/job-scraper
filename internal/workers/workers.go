package workers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/utils"
	"log"
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

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: worker.Config().TelegramConfig().Timeout(),
	})
	if err != nil {
		log.Fatalf("Error getting updates channel: %v", err)
	}

	mailbox := actor.NewMailbox[models.BotMsg]()
	svc := service.NewService(worker.Application)

	processor := actor.New(&telegramSenderWorker{
		bot:     bot,
		logger:  worker.Logger(),
		mailbox: mailbox,
	})

	poller := actor.New(&telegramReceiverWorker{
		bot:          bot,
		logger:       worker.Logger(),
		mailbox:      mailbox,
		appCtx:       worker.Context(),
		updates:      updates,
		userSessions: make(map[int64]*models.UserTelegramSession),
		svc:          svc,
		Client:       worker.Clients(),
	})

	cron := actor.New(&cronWorker{
		svc:    svc,
		logger: worker.Logger(),
		Client: worker.Clients(),
	})

	actors := actor.Combine(mailbox, processor, poller, cron).Build()
	actors.Start()

	worker.Logger().Info(fmt.Sprintf("Worker started at %s...", time.Now().Format("2006-01-02T15:04:05 MST")))
	<-utils.WaitForTermination(worker.Cancel())
}
