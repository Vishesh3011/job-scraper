package config

import (
	"errors"
	"syscall"
)

type Config interface {
	DBConfig() *dbConfig
	EmailConfig() *emailConfig
	TelegramConfig() *telegramConfig
	EncryptionKey() string
}

type config struct {
	*dbConfig
	*emailConfig
	*telegramConfig
	encryptionKey string
}

func NewConfig() (Config, error) {
	dbConfig, err := newDBConfig()
	if err != nil {
		return nil, err
	}

	emailConfig, err := newEmailConfig()
	if err != nil {
		return nil, err
	}

	telegramConfig, err := newTelegramConfig()
	if err != nil {
		return nil, err
	}

	key, found := syscall.Getenv("ENCRYPTION_KEY")
	if !found {
		err := errors.New("ENCRYPTION_KEY environment variable not set")
		return nil, err
	}

	return &config{
		dbConfig:       dbConfig,
		emailConfig:    emailConfig,
		telegramConfig: telegramConfig,
		encryptionKey:  key,
	}, nil
}

func (c *config) DBConfig() *dbConfig {
	return c.dbConfig
}

func (c *config) EmailConfig() *emailConfig {
	return c.emailConfig
}

func (c *config) TelegramConfig() *telegramConfig {
	return c.telegramConfig
}

func (c *config) EncryptionKey() string {
	return c.encryptionKey
}
