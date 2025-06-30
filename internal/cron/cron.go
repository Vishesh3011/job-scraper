package cron

import (
	"github.com/robfig/cron/v3"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/utils"
	"log"
)

type Cron struct {
	application.Application
}

func NewCron(application application.Application) *Cron {
	return &Cron{application}
}

func (c *Cron) Start() {
	cr := cron.New()
	_, err := cr.AddFunc("0 9 * * *", func() {})
	if err != nil {
		log.Fatal(err)
	}
	if err := handleSendNotification(c.Application); err != nil {
		log.Fatalf("Error sending notifications: %v", err)
	}
	<-utils.WaitForTermination(c.Cancel())
}

func handleSendNotification(app application.Application) error {
	svc := service.NewService(app)
	user, err := svc.User().GetAllUsers()
	if err != nil {
		return err
	}

	for _, u := range user {
		if u.Email == nil {
			return err
		}
		ui := models.NewUserInput(u.Name, u.Email, u.CsrfToken, u.Cookie, u.Locations, u.Keywords)
		jobs, err := svc.Accumulator().FetchJobs(ui)
		if err != nil {
			return err
		}

		file, err := svc.Report().GenerateReport(jobs, u.Name)
		if err != nil {
			return err
		}

		if err := app.Clients().GoMailClient().SendEmail(ui, file, len(jobs)); err != nil {
			return err
		} else {
			log.Printf("Email sent successfully to user %d", u.Name)
		}
	}
	return nil
}
