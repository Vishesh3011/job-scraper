package service

import (
	"fmt"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/utils"
)

type AccumulatorService interface {
	FetchJobs(user *models.User) ([]models.Job, error)
}

type accumulatorService struct {
	client.Client
	config.Config
}

func newAccumulatorService(client client.Client, config config.Config) AccumulatorService {
	return &accumulatorService{
		Client: client,
		Config: config,
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
				return nil, fmt.Errorf("error fetching LinkedIn job IDs: %w", err)
			}

			for _, jobId := range geoIds {
				job, err := a.JobClient().GetLinkedInJobDetails(jobId, dToken, dCookie)
				if err != nil {
					return nil, fmt.Errorf("error fetching LinkedIn job details for ID %s: %w", jobId, err)
				}
				j := job.ToJob()
				jobs = append(jobs, *j)
			}
		}
	}
	return jobs, nil
}
