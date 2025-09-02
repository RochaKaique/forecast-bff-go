package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/RochaKaique/forecastgo/internal/coordinates"
	"github.com/RochaKaique/forecastgo/internal/forecast"
)

type mockCoord struct {
	coord coordinates.Coordinates
	err   error
}

type mockForecast struct {
	out forecast.ForecastResponse
	err error
}

func (mc *mockCoord) GetCoordinates(ctx context.Context, zip string) (coordinates.Coordinates, error) {
	return mc.coord, mc.err
}

func (mc *mockForecast) GetForecast(ctx context.Context, request *forecast.ForecastRequest) (forecast.ForecastResponse, error) {
	return mc.out, mc.err
}

func TestForecastByZipSuccess(t *testing.T) {
	t.Parallel()
	mkCoord := &mockCoord{coord: coordinates.Coordinates{Latitude: "40.7", Longitude: "-74.0"}}
	mkForecast := &mockForecast{out: forecast.ForecastResponse{CurrentTemperature: 25.0}}
	svc := New(mkCoord, mkForecast, slog.Default())

	got, err := svc.ForecastByZip(context.Background(), "10001", QueryParams{
		StartDate: "2025-06-06", EndDate: "2025-06-07", Units: "metric", TempScale: "C",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.CurrentTemperature != 25.0 {
		t.Fatalf("want 25.0 got %v", got.CurrentTemperature)
	}
}

func TestErrorWhenCoordError(t *testing.T) {
	t.Parallel()
	svc := New(&mockCoord{err: assertErr{}}, &mockForecast{}, slog.Default())
	_, err := svc.ForecastByZip(context.Background(), "00000", QueryParams{})
    if err == nil {
        t.Fatalf("expected error")
    }
}

func TestErrorWhenForecastError(t *testing.T) {
	t.Parallel()
	svc := New(&mockCoord{}, &mockForecast{err: assertErr{}}, slog.Default())
	_, err := svc.ForecastByZip(context.Background(), "00000", QueryParams{})
    if err == nil {
        t.Fatalf("expected error")
    }
}

type assertErr struct{}
func (assertErr) Error() string { return "coords failed" }
