package workers

import (
	"fmt"
	"job-scraper.go/internal/utils"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/service"
)

type cronWorker struct {
	svc    *service.Service
	logger *slog.Logger
	client.Client
	inC <-chan bool
}

func (w *cronWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	case msg, ok := <-w.inC:
		if !ok {
			return actor.WorkerEnd
		}
		if msg {
			if err := w.handleSendNotification(); err != nil {
				w.logger.Error(utils.PrepareLogMsg(fmt.Sprintf("failed to handle send report: %v", err)))
			}
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
			w.logger.Info(fmt.Sprintf("Email sent successfully to user %s", u.Name))
		}
	}
	return nil
}
