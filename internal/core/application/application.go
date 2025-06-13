package application

import (
	"context"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/repository"
)

type Application interface {
	Context() context.Context
	Queries() *repository.Queries
	Config() *config.Config
}

type application struct {
	context context.Context
	config  *config.Config
	queries *repository.Queries
}

func NewApplication(config *config.Config) (Application, error) {
	ctx, _ := context.WithCancel(context.Background())

	queries, err := config.Connect()
	if err != nil {
		return nil, err
	}
	return &application{
		context: ctx,
		config:  config,
		queries: queries,
	}, nil
}

func (application *application) Config() *config.Config {
	return application.config
}

func (application *application) Context() context.Context {
	return application.context
}

func (application *application) Queries() *repository.Queries {
	return application.queries
}
