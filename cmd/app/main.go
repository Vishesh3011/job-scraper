package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/cli"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/service"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	app, err := application.NewApplication()
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error creating application: %v", err))
		os.Exit(1)
	}

	userInput, err := cli.GetUserInputFromCLI(app)
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error getting user input: %v", err))
		os.Exit(1)
	}

	svc := service.NewService(app)
	jobs, err := svc.Accumulator().FetchJobs(userInput)
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error fetching jobs: %v", err))
		os.Exit(1)
	}

	file, err := svc.Report().GenerateReport(jobs, userInput.Name)
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error generating report: %v", err))
		os.Exit(1)
	} else {
		app.Logger().Info(fmt.Sprintf("Generated report: %s", file))
	}

	if userInput.Email != nil && *userInput.Email != "" {
		if err := app.Clients().GoMailClient().SendEmail(userInput, file, len(jobs)); err != nil {
			app.Logger().Error(fmt.Sprintf("Error sending email: %v", err))
		} else {
			app.Logger().Info(fmt.Sprintf("Email sent: %s", file))
		}
	}
}
