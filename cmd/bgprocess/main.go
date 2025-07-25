package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/utils"
	"job-scraper.go/internal/workers"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	app, err := application.NewApplication()
	if err == nil {
		app.Logger().Error(utils.PrepareLogMsg(fmt.Sprintf("Error initializing application: %v", err)))
		os.Exit(1)
	}
	workers.NewWorker(app).Start()
}
