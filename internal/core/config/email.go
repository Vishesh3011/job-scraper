package config

import (
	"errors"
	"strconv"
	"syscall"
)

type emailConfig struct {
	host string
	port int
}

func newEmailConfig() (*emailConfig, error) {
	host, found := syscall.Getenv("SMTP_HOST")
	if !found {
		return nil, errors.New("SMTP_HOST environment variable not set")
	}

	portStr, found := syscall.Getenv("SMTP_PORT")
	if !found {
		return nil, errors.New("SMTP_PORT environment variable not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("SMTP_PORT must be an integer")
	}

	return &emailConfig{
		host: host,
		port: port,
	}, nil
}

func (c *emailConfig) EmailHostName() string {
	return c.host
}

func (c *emailConfig) EmailPort() int {
	return c.port
}
