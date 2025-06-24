package cli

import (
	"encoding/json"
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

	var jobRoles string
	fmt.Println("Enter your interested programming languages for job role (can add multiple keywords separated by commas): ")
	if _, err := fmt.Scanln(&jobRoles); err != nil {
		log.Fatal(err)
	}

	keywords := strings.Split(jobRoles, ",")
	keywordsJson, err := json.Marshal(keywords)
	if err != nil {
		return nil, fmt.Errorf("error marshalling keywords: %w", err)
	}

	var geoIds []string
	fmt.Println("Enter your interested geo ids for locations from linkedin (can add multiple locations separated by commas): ")
	var geoIdsStr string
	if _, err := fmt.Scanln(&geoIdsStr); err != nil {
		log.Fatal(err)
	}
	geoIds = strings.Split(geoIdsStr, ",")
	locationsJson, err := json.Marshal(geoIds)
	if err != nil {
		return nil, fmt.Errorf("error marshalling geo ids: %w", err)
	}

	var emailNotify string
	fmt.Println("Do you want to receive job notifications? (y/n): ")
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
			log.Print("User already exists! Would you like to update your job roles? (y/n): ")
			user, err := app.Queries().GetUserByEmail(app.Context(), *email)
			if err != nil {
				return nil, fmt.Errorf("error retrieving user by email: %w", err)
			}

			if err := app.Queries().UpdateUser(app.Context(), repository.UpdateUserParams{
				Name:      user.Name,
				Location:  user.Location,
				Keywords:  keywordsJson,
				Cookie:    user.Cookie,
				CsrfToken: user.CsrfToken,
				Email:     *email,
			}); err != nil {
				return nil, fmt.Errorf("error updating user: %w", err)
			}
		} else {
			if err := app.Queries().CreateUser(app.Context(), repository.CreateUserParams{
				Name:     name,
				Location: locationsJson,
				Email:    *email,
				Keywords: keywordsJson,
			}); err != nil {
				return nil, fmt.Errorf("error creating user: %w", err)
			}
		}
	}

	return models.NewUserInput(name, email, nil, nil, keywords, geoIds), nil
}
