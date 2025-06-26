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
	return &Service{
		user:        newUserService(app.Context(), repository.New(app.DBConn())),
		accumulator: newAccumulatorService(app.Clients(), app.Config()),
		report:      newReportService(),
		telegram:    newTelegramService(),
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
