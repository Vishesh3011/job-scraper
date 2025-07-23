package models

import (
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
)

type Job struct {
	Title           string         `json:"title"`
	Company         string         `json:"company"`
	JobLocation     string         `json:"jobLocation"`
	Description     string         `json:"description"`
	Link            string         `json:"link"`
	ApplyLink       *string        `json:"applyLink"`
	StaffCount      *int           `json:"staffCount"`
	HeadquarterCity *string        `json:"headquarterCity"`
	ApplicantsCount *int           `json:"applicantsCount"`
	Platform        types.Platform `json:"platform"`
	ExpiryDate      *int64         `json:"expiryDate"`
	WorkType        string         `json:"workType"`
}

func ToJobFromJora(
	title, company, location, description, link string,
) *Job {
	return &Job{
		Title:           title,
		Company:         company,
		JobLocation:     location,
		Description:     description,
		Link:            link,
		ApplyLink:       utils.ToPtr(link),
		StaffCount:      nil,
		HeadquarterCity: nil,
		ApplicantsCount: nil,
		Platform:        types.JORA,
		ExpiryDate:      nil,
	}
}

func ToJobFromGlassdoor(
	title, company, location, description, link string,
) *Job {
	return &Job{
		Title:           title,
		Company:         company,
		JobLocation:     location,
		Description:     description,
		Link:            link,
		ApplyLink:       nil,
		StaffCount:      nil,
		HeadquarterCity: nil,
		ApplicantsCount: nil,
		Platform:        types.GLASSDOOR,
		ExpiryDate:      nil,
	}
}
