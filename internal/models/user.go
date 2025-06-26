package models

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
