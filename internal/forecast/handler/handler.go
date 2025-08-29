package handler

import (
	"context"

	"log/slog"

	"github.com/RochaKaique/forecastgo/internal/forecast"
	"github.com/RochaKaique/forecastgo/internal/forecast/service"
	"github.com/gofiber/fiber/v2"
)

type ForecastService interface {
	ForecastByZip(ctx context.Context, zipcode string, qp service.QueryParams) (forecast.ForecastResponse, error)
}

func Register(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		zipcode := c.Params("zipcode")
		params := service.QueryParams{
			StartDate: c.Query("start_date", ""),
			EndDate:   c.Query("end_date", ""),
			Units:     c.Query("units", ""),
			TempScale: c.Query("temp_scale", ""),
		}

		out, err := svc.ForecastByZip(c.UserContext(), zipcode, params)
		if err != nil {
			slog.ErrorContext(c.UserContext(), "err", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(out)
	}
}
