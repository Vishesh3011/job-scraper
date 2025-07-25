package service

import (
	"fmt"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/utils"
	"log/slog"
)

type AccumulatorService interface {
	FetchJobs(user *models.User) ([]models.Job, error)
}

type accumulatorService struct {
	client.Client
	config.Config
	logger *slog.Logger
}

func newAccumulatorService(client client.Client, config config.Config, logger *slog.Logger) AccumulatorService {
	return &accumulatorService{
		Client: client,
		Config: config,
		logger: logger,
	}
}

func (a accumulatorService) FetchJobs(user *models.User) ([]models.Job, error) {
	dToken, err := utils.DecryptStr(user.CsrfToken, a.EncryptionKey())
	if err != nil {
		return nil, err
	}

	dCookie, err := utils.DecryptStr(user.Cookie, a.EncryptionKey())
	if err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, keyword := range user.Keywords {
		for _, geoId := range user.Locations {
			geoIds, err := a.JobClient().GetLinkedInJobIds(geoId, keyword, dToken, dCookie)
			if err != nil {
				a.logger.Error(utils.PrepareLogMsg(fmt.Sprintf("Failed to get linkedin job ids: %v", err)))
				return nil, err
			}

			for _, jobId := range geoIds {
				job, err := a.JobClient().GetLinkedInJobDetails(jobId, dToken, dCookie)
				if err != nil {
					a.logger.Error(utils.PrepareLogMsg(fmt.Sprintf("Failed to get linkedin job details: %v", err)))
					return nil, err
				}
				j := job.ToJob()
				jobs = append(jobs, *j)
			}
		}
	}
	return jobs, nil
}
