package service

import (
	"fmt"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/utils"
)

type Accumulator struct {
	Jobs  []models.Job `json:"jobs"`
	Count int          `json:"count"`
}

func NewAccumulator(application application.Application, user *models.UserInput) (*Accumulator, error) {
	var jobs []models.Job
	for _, keyword := range user.Keywords {
		for _, geoId := range user.Locations {
			csrfToken := user.CsrfToken
			if csrfToken == nil {
				csrfToken = utils.ToPtr(application.Config().CsrfToken())
			}

			cookie := user.Cookie
			if cookie == nil {
				cookie = utils.ToPtr(application.Config().Cookie())
			}

			geoIds, err := client.GetLinkedInJobIds(geoId, keyword, *csrfToken, *cookie)
			if err != nil {
				return nil, fmt.Errorf("error fetching LinkedIn job IDs: %w", err)
			}

			for _, jobId := range geoIds {
				job, err := client.GetLinkedInJobDetails(jobId, *csrfToken, *cookie)
				if err != nil {
					return nil, fmt.Errorf("error fetching LinkedIn job details for ID %s: %w", jobId, err)
				}
				j := job.ToJob()
				jobs = append(jobs, *j)
			}
		}
	}
	return &Accumulator{
		Jobs:  jobs,
		Count: len(jobs),
	}, nil
}
