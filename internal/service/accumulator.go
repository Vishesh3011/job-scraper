package service

import (
	"fmt"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/utils"
)

type AccumulatorService interface {
	FetchJobs(*models.UserInput) ([]models.Job, error)
}

type accumulatorService struct {
	*client.Client
	*config.Config
}

func newAccumulatorService(client *client.Client, config *config.Config) AccumulatorService {
	return &accumulatorService{
		Client: client,
		Config: config,
	}
}

func (a accumulatorService) FetchJobs(user *models.UserInput) ([]models.Job, error) {
	var jobs []models.Job
	for _, keyword := range user.Keywords {
		for _, geoId := range user.Locations {
			csrfToken := user.CsrfToken
			if csrfToken == nil {
				csrfToken = utils.ToPtr(a.Config.CsrfToken())
			}

			cookie := user.Cookie
			if cookie == nil {
				cookie = utils.ToPtr(a.Config.Cookie())
			}

			geoIds, err := a.GetLinkedInJobIds(geoId, keyword, *csrfToken, *cookie)
			if err != nil {
				return nil, fmt.Errorf("error fetching LinkedIn job IDs: %w", err)
			}

			for _, jobId := range geoIds {
				job, err := a.GetLinkedInJobDetails(jobId, *csrfToken, *cookie)
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
