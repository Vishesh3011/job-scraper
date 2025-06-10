package application

import "job-scraper.go/internal/core/config"

type Application interface {
	GetConfig() *config.Config
}

type application struct {
	config *config.Config
}

func NewApplication(config *config.Config) Application {
	return &application{config: config}
}

func (application *application) GetConfig() *config.Config {
	return application.config
}
