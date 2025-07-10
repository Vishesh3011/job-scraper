package cli

import (
	"errors"
	"fmt"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
	"log"
	"strings"
)

func GetUserInputFromCLI(app application.Application) (*models.User, error) {
	var name string
	fmt.Println("Enter your name: ")
	if _, err := fmt.Scanln(&name); err != nil {
		return nil, err
	}

	var jobRoles string
	fmt.Println("Enter your interested programming languages for job role (can add multiple keywords separated by commas): ")
	if _, err := fmt.Scanln(&jobRoles); err != nil {
		return nil, err
	}
	keywords := strings.Split(jobRoles, ",")

	var geoIds []string
	fmt.Println("Enter your interested geo ids for locations from linkedin (can add multiple locations separated by commas): ")
	var geoIdsStr string
	if _, err := fmt.Scanln(&geoIdsStr); err != nil {
		return nil, err
	}
	geoIds = strings.Split(geoIdsStr, ",")

	cookie, err := utils.ReadMultilineInput("Enter your linkedin cookie (press Enter twice to end):")
	if err != nil {
		return nil, err
	}

	csrfToken, err := utils.ReadMultilineInput("Enter your CSRF token (press Enter twice to end):")
	if err != nil {
		return nil, err
	}

	var emailNotify string
	fmt.Println("Do you want to receive job notifications? (y/n): ")
	if _, err := fmt.Scanln(&emailNotify); err != nil {
		return nil, err
	}

	var email *string
	if strings.ToLower(emailNotify) == "y" {
		fmt.Print("Enter your email address: ")
		tempEmail := ""
		if _, err := fmt.Scanln(&tempEmail); err != nil {
			return nil, err
		}
		email = &tempEmail
	}

	createdUser := &models.User{}
	if email != nil {
		userService := service.NewService(app).User()
		user, err := userService.GetUserByEmail(*email)
		ui := models.NewUserInput(name, cookie, csrfToken, email, keywords, geoIds)
		if err != nil {
			if errors.Is(err, types.ErrRecordNotFound) {
				if _, err := userService.CreateUser(ui); err != nil {
					return nil, fmt.Errorf("error creating user: %w", err)
				}
			} else {
				return nil, err
			}
		}

		if user != nil {
			log.Print("User already exists! Would you like to update your job roles? (y/n): ")
			if _, err := userService.UpdateUser(ui); err != nil {
				return nil, fmt.Errorf("error updating user: %w", err)
			}
		}
		createdUser = user
	} else {
		user, err := service.NewService(app).User().CreateUser(models.NewUserInput(name, cookie, csrfToken, nil, keywords, geoIds))
		if err != nil {
			return nil, fmt.Errorf("error creating user: %w", err)
		}
		createdUser = user
	}

	return createdUser, nil
}
