package models

type UserInput struct {
	Name        string
	Email       *string
	Location    string
	JobKeywords []string
}

func NewUserInput(name, location string, email *string, keywords []string) *UserInput {
	return &UserInput{
		Name:        name,
		Email:       email,
		Location:    location,
		JobKeywords: keywords,
	}
}
