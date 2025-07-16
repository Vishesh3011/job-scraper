package workers

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/service"
	"log"
	"log/slog"
)

type cronWorker struct {
	svc    *service.Service
	logger *slog.Logger
	client.Client
	c *cron.Cron
}

func (w *cronWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		if w.c != nil {
			w.c.Stop()
		}
		return actor.WorkerEnd
	default:
		if w.c == nil {
			w.c = cron.New()
			_, err := w.c.AddFunc("0 9 * * *", func() {
				if err := w.handleSendNotification(); err != nil {
					w.logger.Error(fmt.Sprintf("failed to handle send report: %v", err))
				}
			})
			if err != nil {
				w.logger.Error(fmt.Sprintf("failed to start the cron job: %v", err))
				return actor.WorkerEnd
			}
			w.c.Start()
		}
		return actor.WorkerContinue
	}
}

func (w *cronWorker) handleSendNotification() error {
	user, err := w.svc.User().GetAllUsers()
	if err != nil {
		return err
	}

	for _, u := range user {
		if u.Email == nil {
			return err
		}
		jobs, err := w.svc.Accumulator().FetchJobs(&u)
		if err != nil {
			return err
		}

		file, err := w.svc.Report().GenerateReport(jobs, u.Name)
		if err != nil {
			return err
		}

		if err := w.GoMailClient().SendEmail(&u, file, len(jobs), fmt.Sprintf("%s_report.xlsx", u.Name)); err != nil {
			return err
		} else {
			log.Printf("Email sent successfully to user %d", u.Name)
		}
	}
	return nil
}
