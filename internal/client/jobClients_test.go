package client

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"job-scraper.go/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestJobClients(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	geoId := "104769905"
	token := os.Getenv("LINKEDIN_CSRF_TOKEN")
	cookie := os.Getenv("LINKEDIN_COOKIE")

	jc := newJobClient(&http.Client{})
	t.Run("GetLinkedInJobIds", func(t *testing.T) {
		cards, err := jc.GetLinkedInJobIds(geoId, "golang", token, cookie)
		if err != nil {
			t.Errorf("error getting linked injected job ids: %v", err)
		}
		assert.Greater(t, len(cards), 1)

		var jobs []models.Job
		for _, id := range cards {
			j, err := jc.GetLinkedInJobDetails(id, token, cookie)
			if err != nil {
				t.Errorf("error getting linked injected job details: %v", err)
			}
			jobs = append(jobs, *j.ToJob())
		}
		assert.Greater(t, len(cards), 1)
	})

	//t.Run("GlassdoorJobs", func(t *testing.T) {
	//	jobs := jc.GetGlassdoorJobs("sydney", "golang")
	//	fmt.Println(len(jobs))
	//})
	//
	//t.Run("JoraJobs", func(t *testing.T) {
	//	jobs := jc.GetJoraJobs("sydney", "golang")
	//	fmt.Println(len(jobs))
	//})
}
