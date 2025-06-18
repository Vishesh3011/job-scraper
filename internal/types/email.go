package types

import (
	"strconv"
	"strings"
)

const From = "Job Scraper"
const Subject = "Your Daily Job Report"
const Body = `
Hello {{userName}},

Please find attached your daily job report. This report includes {{jobCnt}} new job postings that match your search criteria.

If you have any questions or suggestions, feel free to reach out.

Best regards,
Job Tracker Team
`

type EmailTemplate struct {
	From    string
	Subject string
	Body    string
}

func NewEmailTemplate(userName string, jobCnt int) *EmailTemplate {
	if userName == "" {
		userName = "User"
	}
	return &EmailTemplate{
		From:    From,
		Subject: Subject,
		Body:    strings.ReplaceAll(strings.ReplaceAll(Body, "{{userName}}", userName), "{{jobCnt}}", strconv.Itoa(jobCnt)),
	}
}
