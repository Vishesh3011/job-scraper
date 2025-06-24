package config

import (
	"errors"
	"syscall"
)

type telegramConfig struct {
	token string
}

func newTelegramConfig() (*telegramConfig, error) {
	token, found := syscall.Getenv("TELEGRAM_BOT_TOKEN")
	if !found {
		err := errors.New("TELEGRAM_BOT_TOKEN environment variable not set")
		return nil, err
	}
	return &telegramConfig{
		token: token,
	}, nil
}

func (config *telegramConfig) Token() string {
	return config.token
}
