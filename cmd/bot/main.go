package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/service"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := application.NewApplication(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := app.Clients().GetTelegramBot(appConfig.Token())
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %v", err)
	}

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: appConfig.Timeout(),
	})
	if err != nil {
		log.Fatalf("Error getting updates from Telegram: %v", err)
	}

	svc := service.NewService(app)

	if err := svc.Telegram().HandleTelegramUpdates(&updates); err != nil {
		log.Fatalf("Error handling Telegram updates: %v", err)
	}
}
