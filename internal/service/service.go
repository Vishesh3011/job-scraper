package service

import (
	"job-scraper.go/internal/core/application"
)

type Service struct {
	accumulator AccumulatorService
	report      ReportService
	telegram    TelegramService
}

func NewService(app application.Application) *Service {
	return &Service{
		accumulator: newAccumulatorService(app),
		report:      newReportService(app),
		telegram:    newTelegramService(app),
	}
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
