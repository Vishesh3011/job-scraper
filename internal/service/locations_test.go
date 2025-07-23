package service

import (
	"fmt"
	"github.com/joho/godotenv"
	"job-scraper.go/internal/core/application"
	"testing"
)

func TestLocationService(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}

	app, err := application.NewApplication()
	if err != nil {
		t.Fatal(err)
	}

	locSvc := NewService(app).Location()

	t.Run("GetGeoIdBasedOnLocation", func(t *testing.T) {
		geoId := locSvc.FetchGeoIdBasedOnLocation("sydney")
		fmt.Println(geoId)
	})
}
