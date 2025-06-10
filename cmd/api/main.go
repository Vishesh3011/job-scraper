package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/accumulator"
	"job-scraper.go/internal/core/application"
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

	var keyword string
	fmt.Println("Enter your interested job role: ")
	if _, err = fmt.Scanln(&keyword); err != nil {
		log.Fatal(err)
	}

	app := application.NewApplication(appConfig)
	_, err = accumulator.NewAccumulator(app, keyword)
	if err != nil {
		log.Fatal(err)
	}
}
