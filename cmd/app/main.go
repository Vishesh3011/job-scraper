package main

import (
	"github.com/joho/godotenv"
	"job-scraper.go/internal/cli"
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

	userInput, err := cli.GetUserInputFromCLI(app)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := service.NewAccumulator(app, userInput)
	if err != nil {
		log.Fatal(err)
	}

	if err := service.GenerateReport(resp.Jobs, userInput.Name); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Report generated successfully!")
	}
}
