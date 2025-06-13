package models

type UserInput struct {
	Name      string
	Email     *string
	Location  string
	Keywords  []string
	Cookie    *string
	CsrfToken *string
}

func NewUserInput(name, location string, email, cookie, csrfToken *string, keywords []string) *UserInput {
	return &UserInput{
		Name:      name,
		Email:     email,
		Location:  location,
		Keywords:  keywords,
		Cookie:    cookie,
		CsrfToken: csrfToken,
	}
}

type User struct {
}
