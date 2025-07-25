package cli

import (
	"errors"
	"fmt"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
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

	svc := service.NewService(app)

	var location, geoIds []string
	var locations string
	fmt.Println("Enter your interested location (state/city) in Australia (can add multiple locations separated by commas): ")
	if _, err := fmt.Scanln(&locations); err != nil {
		return nil, err
	}
	location = strings.Split(locations, ",")
	for _, l := range location {
		id := svc.Location().FetchGeoIdBasedOnLocation(l)
		if id == "" {
			return nil, errors.New(fmt.Sprintf("Could not find location for %s", l))
		}
		geoIds = append(geoIds, id)
	}

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
		fmt.Println("Enter your email address: ")
		tempEmail := ""
		if _, err := fmt.Scanln(&tempEmail); err != nil {
			return nil, err
		}
		email = &tempEmail
	}

	if email != nil {
		userService := svc.User()
		user, err := userService.GetUserByEmail(*email)
		ui := models.NewUserInput(name, cookie, csrfToken, email, keywords, location)
		if err != nil {
			if errors.Is(err, types.ErrRecordNotFound) {
				u, err := userService.CreateUser(ui)
				if err != nil {
					return nil, err
				}
				return u, nil
			} else {
				return nil, err
			}
		}

		if user != nil {
			fmt.Println("User already exists! Updating the existing user...")
			u, err := userService.UpdateUser(ui)
			if err != nil {
				return nil, fmt.Errorf("error updating user: %w", err)
			}
			return u, nil
		}
	} else {
		user, err := service.NewService(app).User().CreateUser(models.NewUserInput(name, cookie, csrfToken, nil, keywords, location))
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}
