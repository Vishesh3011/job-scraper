package main

import (
	"github.com/joho/godotenv"
	"job-scraper.go/internal/cli"
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

	userInput, err := cli.GetUserInputFromCLI(app)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(app)

	jobs, err := svc.Accumulator().FetchJobs(userInput)
	if err != nil {
		log.Fatal(err)
	}

	file, err := svc.Report().GenerateReport(jobs, userInput.Name)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Report generated successfully!")
	}

	if userInput.Email != nil && *userInput.Email != "" {
		if err := client.SendEmail(userInput, file, len(jobs), app.Config().EmailHostName(), app.Config().EmailPort()); err != nil {
			log.Fatalf("Error sending email: %v", err)
		} else {
			log.Println("Email sent successfully!")
		}
	}
}
