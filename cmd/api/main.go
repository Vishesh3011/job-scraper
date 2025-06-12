package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/accumulator"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/models"
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

	var name string
	fmt.Println("Enter your name: ")
	if _, err = fmt.Scanln(&name); err != nil {
		log.Fatal(err)
	}

	var jobKeyword string
	fmt.Println("Enter your interested job role: ")
	if _, err = fmt.Scanln(&jobKeyword); err != nil {
		log.Fatal(err)
	}

	var location string
	fmt.Println("Enter your interested location: ")
	if _, err = fmt.Scanln(&location); err != nil {
		log.Fatal(err)
	}

	var emailNotify string
	fmt.Println("Do you want to receive email notifications? (y/n): ")
	if _, err = fmt.Scanln(&emailNotify); err != nil {
		log.Fatal(err)
	}

	var email *string
	if emailNotify == "y" {
		fmt.Println("Enter your email address: ")
		if _, err = fmt.Scanln(&email); err != nil {
			log.Fatal(err)
		}
	}

	userInput := models.NewUserInput(name, jobKeyword, location, email)

	app := application.NewApplication(appConfig)
	_, err = accumulator.NewAccumulator(app, jobKeyword)
	if err != nil {
		log.Fatal(err)
	}
}
