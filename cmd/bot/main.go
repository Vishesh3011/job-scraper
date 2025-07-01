package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/utils"
	"log"
	"math/rand"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	app, err := application.NewApplication()
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error initializing application: %v", err))
		os.Exit(1)
	}

	bot, err := app.Clients().TelegramClient().GetTelegramBot(app.Config().TelegramConfig().Token())
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error creating Telegram bot: %v", err))
		os.Exit(1)
	}

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: app.Config().TelegramConfig().Timeout()})
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error getting updates from Telegram: %v", err))
		os.Exit(1)
	}

	app.Logger().Info("Telegram bot started. Listening for updates. Press Ctrl+C to stop.")
	go func() {
		svc := service.NewService(app).Telegram()
		if err := svc.HandleTelegramUpdates(bot, &updates); err != nil {
			app.Logger().Error(fmt.Sprintf("Error handling Telegram updates: %v", err))
			os.Exit(1)
		}
	}()

	<-utils.WaitForTermination(app.Cancel())

	msg := tgbotapi.NewMessage(int64(rand.Uint64()), "Bot is shutting down. See you later!")
	if _, err := bot.Send(msg); err != nil {
		app.Logger().Error(fmt.Sprintf("Failed to send shutdown message: %v", err))
	}
	app.Logger().Info("Bot shutdown gracefully")
}
