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

	user, err := cli.GetUserInputFromCLI(app)
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error getting user input: %v", err))
		os.Exit(1)
	}
	fmt.Println(user)

	svc := service.NewService(app)
	jobs, err := svc.Accumulator().FetchJobs(user)
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error fetching jobs: %v", err))
		os.Exit(1)
	}
	fmt.Println("hahahaha")
	fmt.Println(jobs)

	file, err := svc.Report().GenerateReport(jobs, user.Name)
	if err != nil {
		app.Logger().Error(fmt.Sprintf("Error generating report: %v", err))
		os.Exit(1)
	} else {
		app.Logger().Info(fmt.Sprintf("Generated report: %s", file))
	}
	fmt.Println("After report generation")

	if user.Email != nil && *user.Email != "" {
		if err := app.Clients().GoMailClient().SendEmail(user, file, len(jobs)); err != nil {
			app.Logger().Error(fmt.Sprintf("Error sending email: %v", err))
		} else {
			app.Logger().Info(fmt.Sprintf("Email sent: %s", file))
		}
	}
}
