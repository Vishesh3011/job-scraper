package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/cli"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/service"
	"log"
	"os"
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

	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error marshalling response: %v", err)
	}

	if err := os.WriteFile("output.json", respJson, 0644); err != nil {
		log.Fatalf("Error writing response to file: %v", err)
	}
}
