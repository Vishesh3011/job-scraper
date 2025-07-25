package workers

import (
	"context"
	"fmt"
	"job-scraper.go/internal/types"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron/v3"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/utils"
)

type worker struct {
	application.Application
}

func NewWorker(app application.Application) *worker {
	return &worker{app}
}

func (worker *worker) Start() {
	ctx, cancel := context.WithCancel(worker.Context())
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
	cronMailBox := actor.NewMailbox[bool]()
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
		appCtx:       ctx,
		updates:      updates,
		userSessions: make(map[int64]*models.UserTelegramSession),
		svc:          svc,
		Client:       worker.Clients(),
	})

	c := cron.New()
	if _, err := c.AddFunc(types.EmailCronTime, func() {
		worker.Logger().Info("Cron job triggered")
		if err := cronMailBox.Send(ctx, true); err != nil {
			worker.Logger().Error("Worker failed to send the cron trigger!")
		}
	}); err != nil {
		worker.Logger().Error("Cron job triggered failed: %v", err)
	}
	c.Start()
	defer c.Stop()

	cronWorker := actor.New(&cronWorker{
		svc:    svc,
		logger: worker.Logger(),
		Client: worker.Clients(),
		inC:    cronMailBox.ReceiveC(),
	})

	actors := actor.Combine(cronMailBox, mailbox, processor, poller, cronWorker).Build()
	actors.Start()
	defer actors.Stop()

	worker.Logger().Info(fmt.Sprintf("Worker started at %s...", time.Now().Format("2006-01-02T15:04:05 MST")))
	<-utils.WaitForTermination(cancel)
}
