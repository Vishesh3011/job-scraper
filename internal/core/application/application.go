package application

import (
	"context"
	"database/sql"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"log"
)

type Application interface {
	Context() context.Context
	Config() config.Config
	Clients() client.Client
	DBConn() *sql.DB
}

type application struct {
	context context.Context
	config  config.Config
	clients client.Client
	dbConn  *sql.DB
}

func NewApplication() (Application, error) {
	ctx, _ := context.WithCancel(context.Background())

	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
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
