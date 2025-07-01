package service

import (
	"job-scraper.go/internal/core/application"
	"job-scraper.go/internal/repository"
)

type Service struct {
	user        UserService
	accumulator AccumulatorService
	report      ReportService
	telegram    TelegramService
}

func NewService(app application.Application) *Service {
	queries := repository.New(app.DBConn())
	return &Service{
		user:        newUserService(app.Context(), queries),
		accumulator: newAccumulatorService(app.Clients(), app.Config()),
		report:      newReportService(),
		telegram:    newTelegramService(app.Context(), queries, app.Config(), app.Clients(), app.Logger()),
	}
}

func (service *Service) User() UserService {
	return service.user
}

func (service *Service) Accumulator() AccumulatorService {
	return service.accumulator
}

func (service *Service) Report() ReportService {
	return service.report
}

func (service *Service) Telegram() TelegramService {
	return service.telegram
}
