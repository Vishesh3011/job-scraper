package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/workers"
	"log"
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

	//bot, err := app.Clients().TelegramClient().GetTelegramBot(app.Config().TelegramConfig().Token())
	//if err != nil {
	//	app.Logger().Error(fmt.Sprintf("Error creating Telegram bot: %v", err))
	//	os.Exit(1)
	//}
	//
	//updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: app.Config().TelegramConfig().Timeout()})
	//if err != nil {
	//	app.Logger().Error(fmt.Sprintf("Error getting updates from Telegram: %v", err))
	//	os.Exit(1)
	//}
	//
	//app.Logger().Info("Telegram bot started. Listening for updates. Press Ctrl+C to stop.")
	//go func() {
	//	svc := service.NewService(app).Telegram()
	//	if err := svc.HandleTelegramUpdates(bot, &updates); err != nil {
	//		app.Logger().Error(fmt.Sprintf("Error handling Telegram updates: %v", err))
	//		os.Exit(1)
	//	}
	//}()
	//
	//<-utils.WaitForTermination(app.Cancel())

	workers.NewWorker(app).Start()
}
