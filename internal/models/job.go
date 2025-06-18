package models

import "job-scraper.go/internal/types"

type Job struct {
	Title           string         `json:"title"`
	Company         string         `json:"company"`
	JobLocation     string         `json:"jobLocation"`
	Description     string         `json:"description"`
	Link            string         `json:"link"`
	ApplyLink       string         `json:"applyLink"`
	StaffCount      int            `json:"staffCount"`
	HeadquarterCity string         `json:"headquarterCity"`
	ApplicantsCount int            `json:"applicantsCount"`
	Platform        types.Platform `json:"platform"`
	ExpiryDate      int64          `json:"expiryDate"`
	WorkType        string         `json:"workType"`
}
