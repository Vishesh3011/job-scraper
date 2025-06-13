package config

import (
	"errors"
	"syscall"
)

type linkedinConfig struct {
	csrfToken, cookie string
}

func newLinkedINConfig() (*linkedinConfig, error) {
	csrfToken, found := syscall.Getenv("LINKEDIN_CSRF_TOKEN")
	if !found {
		err := errors.New("LINKEDIN_CSRF_TOKEN environment variable not set")
		return nil, err
	}

	cookie, found := syscall.Getenv("LINKEDIN_COOKIE")
	if !found {
		err := errors.New("LINKEDIN_COOKIE environment variable not set")
		return nil, err
	}

	return &linkedinConfig{
		csrfToken: csrfToken,
		cookie:    cookie,
	}, nil
}

func (c *linkedinConfig) CsrfToken() string {
	return c.csrfToken
}

func (c *linkedinConfig) Cookie() string {
	return c.cookie
}
