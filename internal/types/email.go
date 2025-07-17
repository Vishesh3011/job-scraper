package types

import (
	"strconv"
	"strings"
)

const From = "modivishesh30@gmail.com"
const Subject = "Your Daily Job Report"
const Body = "Hello {{userName}},\n\nPlease find attached your daily job report. This report includes {{jobCnt}} new job postings that match your search criteria.\n\nIf you have any questions or suggestions, feel free to reach out.\n\nBest regards,\nJob Scraper Team"

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
