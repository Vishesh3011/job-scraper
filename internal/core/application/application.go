package application

import (
	"context"
	"database/sql"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
)

type Application interface {
	Context() context.Context
	DBConn() *sql.DB
	Config() *config.Config
	Clients() *client.Client
}

type application struct {
	context context.Context
	dbConn  *sql.DB
	config  *config.Config
	clients *client.Client
}

func NewApplication(config *config.Config) (Application, error) {
	ctx, _ := context.WithCancel(context.Background())

	dbConn, err := config.DBConn()
	if err != nil {
		return nil, err
	}

	clients := &client.Client{
		GoMailClient: client.NewGoMailClient(
			config.EmailHostName(),
			config.EmailPort(),
		),
		LinkedInClient: client.NewLinkedInClient(
			config.CsrfToken(),
			config.Cookie(),
		),
		TelegramClient: client.NewTelegramClient(
			config.Token(),
		),
	}

	return &application{
		context: ctx,
		dbConn:  dbConn,
		config:  config,
		clients: clients,
	}, nil
}

func (application *application) Config() *config.Config {
	return application.config
}

func (application *application) Context() context.Context {
	return application.context
}

func (application *application) DBConn() *sql.DB {
	return application.dbConn
}

func (application *application) Clients() *client.Client {
	return application.clients
}
