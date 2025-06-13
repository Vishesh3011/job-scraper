package cli

import (
	"fmt"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/repository"
	"log"
	"strings"
)

func GetUserInputFromCLI(app application.Application) (*models.UserInput, error) {
	var name string
	fmt.Println("Enter your name: ")
	if _, err := fmt.Scanln(&name); err != nil {
		log.Fatal(err)
	}

	var keyword string
	fmt.Println("Enter your interested job role (can add multiple keywords separated by commas): ")
	if _, err := fmt.Scanln(&keyword); err != nil {
		log.Fatal(err)
	}

	keywords := strings.Split(keyword, ",")

	var location string
	fmt.Println("Enter your interested location: ")
	if _, err := fmt.Scanln(&location); err != nil {
		log.Fatal(err)
	}

	var emailNotify string
	fmt.Println("Do you want to receive email notifications? (y/n): ")
	if _, err := fmt.Scanln(&emailNotify); err != nil {
		log.Fatal(err)
	}

	var email *string
	if strings.ToLower(emailNotify) == "y" {
		fmt.Print("Enter your email address: ")
		tempEmail := ""
		if _, err := fmt.Scanln(&tempEmail); err != nil {
			log.Fatal(err)
		}

		email = &tempEmail
	}

	if email != nil {
		exists, err := app.Queries().CheckUserExistsByEmail(app.Context(), *email)
		if err != nil {
			return nil, fmt.Errorf("error checking if user exists: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("user with email %s already exists", *email)
		} else {
			if err := app.Queries().CreateUser(app.Context(), repository.CreateUserParams{
				Name:     name,
				Location: location,
				Email:    *email,
				Keywords: keywords[0],
			}); err != nil {
				return nil, fmt.Errorf("error creating user: %w", err)
			}
		}
	}

	return models.NewUserInput(name, location, email, nil, nil, keywords), nil
}
