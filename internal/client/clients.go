package client

import (
	"job-scraper.go/internal/core/config"
)

type Client interface {
	GoMailClient() *goMailClient
	LinkedInClient() *linkedInClient
	TelegramClient() *telegramClient
}

type client struct {
	*goMailClient
	*linkedInClient
	*telegramClient
}

func NewClient(config config.Config) Client {
	return &client{
		newGoMailClient(config.EmailConfig().EmailHostName(), config.EmailConfig().EmailPort()),
		newLinkedInClient(),
		newTelegramClient(config.TelegramConfig().Token()),
	}
}

func (c *client) GoMailClient() *goMailClient {
	return c.goMailClient
}

func (c *client) LinkedInClient() *linkedInClient {
	return c.linkedInClient
}

func (c *client) TelegramClient() *telegramClient {
	return c.telegramClient
}
