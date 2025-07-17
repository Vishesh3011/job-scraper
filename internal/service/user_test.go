package service

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"testing"
)

func TestUserService(t *testing.T) {
	if err := godotenv.Load("/Users/vishesh/Desktop/vishesh./projects/go/job-scraper/.env"); err != nil {
		t.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}
	app, err := application.NewApplication()
	if err != nil {
		t.Fatal(err)
	}

	svc := NewService(app).User()
	t.Run("TestGetAllUsers", func(t *testing.T) {
		users, err := svc.GetAllUsers()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(len(users))
	})
}
