package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gopkg.in/gomail.v2"
	"io"
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

func (c *goMailClient) SendEmail(user *models.User, file *excelize.File, jobCount int, fileName string) error {
	dialer := gomail.NewDialer(c.hostName, c.port, "", "")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	template := types.NewEmailTemplate(user.Name, jobCount)

	m := gomail.NewMessage()
	m.SetHeader("From", template.From)
	m.SetHeader("To", *user.Email)
	m.SetHeader("Subject", template.Subject)
	m.SetBody("text/html", template.Body)

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return fmt.Errorf("failed to write excel to buffer: %w", err)
	}
	m.Attach(fmt.Sprintf("%s.xlsx", fileName), gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, &buf)
		return err
	}))

	return dialer.DialAndSend(m)
}
