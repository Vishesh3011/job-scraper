package application

import (
	"context"
	"database/sql"
	"log/slog"

	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
)

type Application interface {
	Context() context.Context
	Config() config.Config
	Clients() client.Client
	DBConn() *sql.DB
	Logger() *slog.Logger
}

type application struct {
	context context.Context
	config  config.Config
	clients client.Client
	dbConn  *sql.DB
	logger  *slog.Logger
}

func NewApplication() (Application, error) {
	ctx := context.Background()
	logger := slog.Default()

	appConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	dbConn, err := appConfig.DBConfig().DBConn()
	if err != nil {
		return nil, err
	}

	return &application{
		context: ctx,
		config:  appConfig,
		clients: client.NewClient(appConfig),
		dbConn:  dbConn,
		logger:  logger,
	}, nil
}

func (application *application) Context() context.Context {
	return application.context
}

func (application *application) Config() config.Config {
	return application.config
}

func (application *application) Clients() client.Client {
	return application.clients
}

func (application *application) DBConn() *sql.DB {
	return application.dbConn
}

func (application *application) Logger() *slog.Logger {
	return application.logger
}
