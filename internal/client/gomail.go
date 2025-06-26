package client

import (
	"crypto/tls"
	"github.com/xuri/excelize/v2"
	"gopkg.in/gomail.v2"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/types"
)

type goMailClient struct {
	hostName string
	port     int
}

func newGoMailClient(hostname string, port int) *goMailClient {
	return &goMailClient{
		hostName: hostname,
		port:     port,
	}
}

func (c *goMailClient) SendEmail(user *models.UserInput, file *excelize.File, jobCount int) error {
	dialer := gomail.NewDialer(c.hostName, c.port, "", "")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	template := types.NewEmailTemplate(user.Name, jobCount)

	m := gomail.NewMessage()
	m.SetHeader("From", template.From)
	m.SetHeader("To", *user.Email)
	m.SetHeader("Subject", template.Subject)
	m.SetBody("text/html", template.Body)
	if file != nil {
		m.Attach(file.Path)
	}
	return dialer.DialAndSend(m)
}
