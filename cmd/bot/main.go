package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/utils"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := app.Clients().TelegramClient().GetTelegramBot(app.Config().TelegramConfig().Token())
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %v", err)
	}

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: app.Config().TelegramConfig().Timeout(),
	})
	if err != nil {
		log.Fatalf("Error getting updates from Telegram: %v", err)
	}

	svc := service.NewService(app).Telegram()
	if err := svc.HandleTelegramUpdates(bot, &updates); err != nil {
		log.Fatalf("Error handling Telegram updates: %v", err)
	}
	<-utils.WaitForTermination(app.Cancel())
}
