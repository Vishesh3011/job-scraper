package main

import (
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/cron"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	cron.NewCron(app).Start()
}
