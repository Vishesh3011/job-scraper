package service

import (
	"context"
	"fmt"
	"job-scraper.go/internal/repository"
	"log/slog"
	"strings"
)

type LocationService interface {
	FetchGeoIdBasedOnLocation(location string) string
}

type locationService struct {
	context.Context
	*repository.Queries
	logger *slog.Logger
}

func newLocationService(ctx context.Context, q *repository.Queries, logger *slog.Logger) LocationService {
	return &locationService{ctx, q, logger}
}

func (l locationService) FetchGeoIdBasedOnLocation(location string) string {
	geoLoc, err := l.GetGeoLocationByLocation(l.Context, strings.ToLower(location))
	if err != nil {
		l.logger.Info(fmt.Sprintf("error fetching geolocation for location %s; err: %v", location, err))
		return ""
	}
	return geoLoc.GeoID
}
