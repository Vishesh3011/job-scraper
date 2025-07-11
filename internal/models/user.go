package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"job-scraper.go/internal/repository"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
)

type User struct {
	Id        []byte   `json:"id"`
	Name      string   `json:"name"`
	Email     *string  `json:"email;omitempty"`
	Locations []string `json:"locations"`
	Keywords  []string `json:"keywords"`
	Cookie    string   `json:"cookie;omitempty"`
	CsrfToken string   `json:"csrf_token;omitempty"`
}

func NewUser(
	name string,
	email *string,
	locations []string,
	keywords []string,
	cookie string,
	csrfToken string,
) *User {
	id := uuid.New()
	return &User{
		Id:        id[:],
		Name:      name,
		Email:     email,
		Locations: locations,
		Keywords:  keywords,
		Cookie:    cookie,
		CsrfToken: csrfToken,
	}
}

func (user *User) ToCreateUserParam(eToken, eCookie string) repository.CreateUserParams {
	loc, _ := json.Marshal(user.Locations)
	keywords, _ := json.Marshal(user.Keywords)

	return repository.CreateUserParams{
		ID:        user.Id,
		Name:      user.Name,
		Email:     utils.ToSQLNullStr(user.Email),
		Location:  loc,
		Keywords:  keywords,
		Cookie:    eCookie,
		CsrfToken: eToken,
	}
}

func ToUser(jsu *repository.JobScraperUser) *User {
	var keywords []string
	_ = json.Unmarshal(jsu.Keywords, &keywords)

	var locations []string
	_ = json.Unmarshal(jsu.Location, &locations)

	return &User{
		Name:      jsu.Name,
		Email:     utils.NullStringToPtr(jsu.Email),
		Locations: locations,
		Keywords:  keywords,
		Cookie:    jsu.Cookie,
		CsrfToken: jsu.CsrfToken,
	}
}

type UserInput struct {
	Name      string
	Email     *string
	Locations []string
	Keywords  []string
	Cookie    string
	CsrfToken string
}

func NewUserInput(name, cookie, csrfToken string, email *string, keywords, locations []string) *UserInput {
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
	Cookie        string
	CsrfToken     string
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
