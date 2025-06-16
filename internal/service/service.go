package service

import (
	"context"
	"job-scraper.go/internal/repository"
)

type Service interface {
	Queries() *repository.Queries
	Context() *context.Context
}

type service struct {
	queries *repository.Queries
	context *context.Context
}

func NewService(queries *repository.Queries, ctx *context.Context) Service {
	return &service{
		queries: queries,
		context: ctx,
	}
}

func (s *service) Queries() *repository.Queries {
	return s.queries
}

func (s *service) Context() *context.Context {
	return s.context
}
