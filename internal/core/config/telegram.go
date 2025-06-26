package config

import (
	"errors"
	"strconv"
	"syscall"
)

type telegramConfig struct {
	token   string
	timeout int
}

func newTelegramConfig() (*telegramConfig, error) {
	token, found := syscall.Getenv("TELEGRAM_BOT_TOKEN")
	if !found {
		err := errors.New("TELEGRAM_BOT_TOKEN environment variable not set")
		return nil, err
	}

	timeoutStr, found := syscall.Getenv("TELEGRAM_BOT_TIMEOUT")
	if !found {
		err := errors.New("TELEGRAM_BOT_TIMEOUT environment variable not set")
		return nil, err
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return nil, err
	}

	return &telegramConfig{
		token:   token,
		timeout: timeout,
	}, nil
}

func (config *telegramConfig) Token() string {
	return config.token
}

func (config *telegramConfig) Timeout() int {
	return config.timeout
}
