package client

import (
	"job-scraper.go/internal/core/config"
	"net/http"
)

type Client interface {
	GoMailClient() *goMailClient
	JobClient() *jobClient
	TelegramClient() *telegramClient
}

type client struct {
	*goMailClient
	*jobClient
	*telegramClient
}

func NewClient(config config.Config) Client {
	httpClient := &http.Client{}

	return &client{
		newGoMailClient(config.EmailConfig().EmailHostName(), config.EmailConfig().EmailPort()),
		newJobClient(httpClient),
		newTelegramClient(config.TelegramConfig().Token()),
	}
}

func (c *client) GoMailClient() *goMailClient {
	return c.goMailClient
}

func (c *client) JobClient() *jobClient {
	return c.jobClient
}

func (c *client) TelegramClient() *telegramClient {
	return c.telegramClient
}
