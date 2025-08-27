package service

import (
	"context"
	"log/slog"

	"github.com/RochaKaique/forecastgo/internal/forecast"
)

type Service struct {
	coordProvider CoordinatesProvider
	fcProvider    ForecastProvider
	logger        *slog.Logger
}

func New(coord CoordinatesProvider, fc ForecastProvider, logger *slog.Logger) *Service {
	return &Service{
		coordProvider: coord,
		fcProvider:    fc,
		logger:        logger,
	}
}

type QueryParams struct {
	StartDate string
	EndDate   string
	Units     string
	TempScale string
}

func (s *Service) ForecastByZip(ctx context.Context, zipCode string, query QueryParams) (forecast.ForecastResponse, error) {
	coordinates, err := s.coordProvider.GetCoordinates(ctx, zipCode)
	if err != nil {
		s.logger.ErrorContext(ctx, "coordinates lookup failed", "zip", zipCode, "err", err)
		return forecast.ForecastResponse{}, err
	}

	req := forecast.ForecastRequest{
		Latitude:  coordinates.Latitude,
		Longitude: coordinates.Longitude,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
		Units:     query.Units,
		TempScale: query.TempScale,
	}

	out, err := s.fcProvider.GetForecast(ctx, &req)
	if err != nil {
		s.logger.ErrorContext(ctx, "forecast provider failed", "zip", zipCode, "err", err)
		return forecast.ForecastResponse{}, err
	}
	return out, nil
}
