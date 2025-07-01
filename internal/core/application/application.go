package application

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"job-scraper.go/internal/client"
	"job-scraper.go/internal/core/config"
	"job-scraper.go/internal/utils"
	"log/slog"
)

type Application interface {
	Context() context.Context
	Config() config.Config
	Clients() client.Client
	DBConn() *sql.DB
	Cancel() context.CancelFunc
	Logger() *slog.Logger
}

type application struct {
	context context.Context
	config  config.Config
	clients client.Client
	dbConn  *sql.DB
	cancel  context.CancelFunc
	logger  *slog.Logger
}

func NewApplication() (Application, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cid := uuid.New().String()
	ctx = context.WithValue(ctx, "cid", cid)
	logger := slog.Default().With("correlation_id", utils.GetCorrelationID(ctx))

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
		cancel:  cancel,
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

func (application *application) Cancel() context.CancelFunc {
	return application.cancel
}

func (application *application) Logger() *slog.Logger {
	return application.logger
}
