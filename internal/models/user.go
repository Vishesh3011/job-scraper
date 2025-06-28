package models

import (
	"encoding/json"
	"job-scraper.go/internal/repository"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
)

type User struct {
	Name      string   `json:"name"`
	Email     *string  `json:"email;omitempty"`
	Locations []string `json:"locations"`
	Keywords  []string `json:"keywords"`
	Cookie    *string  `json:"cookie;omitempty"`
	CsrfToken *string  `json:"csrf_token;omitempty"`
}

func NewUser(jsu *repository.JobScraperUser) *User {
	var keywords []string
	_ = json.Unmarshal([]byte(jsu.Keywords), &keywords)

	var locations []string
	_ = json.Unmarshal(jsu.Location, &locations)

	return &User{
		Name:      jsu.Name,
		Email:     &jsu.Email,
		Locations: locations,
		Keywords:  keywords,
		Cookie:    utils.NullStringToPtr(jsu.Cookie),
		CsrfToken: utils.NullStringToPtr(jsu.CsrfToken),
	}
}

type UserInput struct {
	Name      string
	Email     *string
	Locations []string
	Keywords  []string
	Cookie    *string
	CsrfToken *string
}

func NewUserInput(name string, email, cookie, csrfToken *string, keywords, locations []string) *UserInput {
	return &UserInput{
		Name:      name,
		Email:     email,
		Locations: locations,
		Keywords:  keywords,
		Cookie:    cookie,
		CsrfToken: csrfToken,
	}
}

type UserTelegramSession struct {
	Name          string
	Email         *string
	Locations     []string
	Keywords      []string
	Cookie        *string
	CsrfToken     *string
	TelegramState types.TELEGRAM_STATE
}

func (uTS *UserTelegramSession) ToUserInput() UserInput {
	return UserInput{
		Name:      uTS.Name,
		Email:     uTS.Email,
		Locations: uTS.Locations,
		Keywords:  uTS.Keywords,
		Cookie:    uTS.Cookie,
		CsrfToken: uTS.CsrfToken,
	}
}
