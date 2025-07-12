package workers

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/service"
	"log"
)

type cronWorker struct {
	app application.Application
}

func (w *cronWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	default:
		c := cron.New()
		_, err := c.AddFunc("0 9 * * *", func() {
			if err := w.handleSendNotification(); err != nil {
				w.app.Logger().Error(fmt.Sprintf("failed to handle send report: %v", err))
			}
		})
		if err != nil {
			w.app.Logger().Error(fmt.Sprintf("failed to start the cron job: %v", err))
			return actor.WorkerEnd
		}
		c.Start()
		return actor.WorkerContinue
	}
}

func (w *cronWorker) handleSendNotification() error {
	svc := service.NewService(w.app)
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

		if err := w.app.Clients().GoMailClient().SendEmail(&u, file, len(jobs)); err != nil {
			return err
		} else {
			log.Printf("Email sent successfully to user %d", u.Name)
		}
	}
	return nil
}
