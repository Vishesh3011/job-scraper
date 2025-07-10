package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/service"
	"job-scraper.go/internal/utils"
	"log"
	"os"
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
		c.Logger().Error(fmt.Sprintf("failed to send notification: %w", err))
		os.Exit(1)
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
		jobs, err := svc.Accumulator().FetchJobs(&u)
		if err != nil {
			return err
		}

		file, err := svc.Report().GenerateReport(jobs, u.Name)
		if err != nil {
			return err
		}

		if err := app.Clients().GoMailClient().SendEmail(&u, file, len(jobs)); err != nil {
			return err
		} else {
			log.Printf("Email sent successfully to user %d", u.Name)
		}
	}
	return nil
}
