package config

import (
	"errors"
	"syscall"
)

type LinkedinConfig struct {
	CsrfToken, Cookie string
}

func newLinkedINConfig() (*LinkedinConfig, error) {
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

	return &LinkedinConfig{
		CsrfToken: csrfToken,
		Cookie:    cookie,
	}, nil
}
