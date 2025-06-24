package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
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

	//app, err := application.NewApplication(appConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	bot, err := client.GetTelegramBot(appConfig.Token())
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %v", err)
	}

	updtConfig := tgbotapi.UpdateConfig{
		Timeout: 60,
	}
	updates, err := bot.GetUpdatesChan(updtConfig)
	if err != nil {
		log.Fatalf("Error getting updates from Telegram: %v", err)
	}

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Text {
			case "/start":
				msg.Text = "Welcome to the Job Scraper Bot! Please enter your name."
			case "/help":
				msg.Text = "This bot helps you find job listings based on your preferences. Please enter your name to get started."
			default:
				msg.Text = "You entered: " + update.Message.Text
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
