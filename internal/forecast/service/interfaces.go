package service

import (
	"context"

	"github.com/RochaKaique/forecastgo/internal/coordinates"
	"github.com/RochaKaique/forecastgo/internal/forecast"
)

type CoordinatesProvider interface {
	GetCoordinates(ctx context.Context, postalCode string) (coordinates.Coordinates, error)
}

type ForecastProvider interface {
	GetForecast(ctx context.Context, request *forecast.ForecastRequest) (forecast.ForecastResponse, error)
}