package handler

import (
	"context"
	"log/slog"

	"github.com/RochaKaique/forecastgo/internal/forecast/service"
	"github.com/gofiber/fiber/v2"
)

type ForecastService interface {
	ForecastByZip(ctx context.Context, zipcode string, qp service.QueryParams) (interface{}, error)
}

// func HandleForecast(client *http.Client, config *viper.Viper) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		coordinateClient := coordinates.NewCooordinatesClient(client, config)
// 		zipcode := c.Params("zipcode")

// 		coordinates, err := coordinateClient.GetCoordinates(c.UserContext(), zipcode)
// 		if err != nil {
// 			return c.SendStatus(fiber.StatusInternalServerError)
// 		}

// 		forecasRequest := &forecast.ForecastRequest{
// 			Latitude:  coordinates.Latitude,
// 			Longitude: coordinates.Longitude,
// 			StartDate: c.Query("start_date", ""),
// 			EndDate:   c.Query("end_date", ""),
// 			Units:     c.Query("units", ""),
// 			TempScale: c.Query("temp_scale"),
// 		}

// 		forecastClient := forecast.NewForecastClient(client, config)
// 		forecast, err := forecastClient.GetForecast(c.UserContext(), forecasRequest)
// 		if err != nil {
// 			return c.SendStatus(fiber.StatusInternalServerError)
// 		}
// 		return c.Status(fiber.StatusOK).JSON(forecast)
// 	}
// }

func Register(app fiber.Router, svc *service.Service) {
	app.Get("/:zipcode", func(c *fiber.Ctx) error {
		slog.Info("Entrou")
		zipcode := c.Params("zipcode")
		params := service.QueryParams{
			StartDate: c.Query("start_date", ""),
			EndDate:   c.Query("end_date", ""),
			Units:     c.Query("units", ""),
			TempScale: c.Query("temp_scale", ""),
		}

		out, err := svc.ForecastByZip(c.UserContext(), zipcode, params)
		if err == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(out)
	})
}
