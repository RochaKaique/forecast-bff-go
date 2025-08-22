package forecast

import (
	"net/http"

	"github.com/RochaKaique/forecastgo/internal/coordinates"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func HandleForecast(client *http.Client, config *viper.Viper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		coordinateClient := coordinates.NewCooordinatesClient(client, config)
		zipcode := c.Params("zipcode")

		coordinates, err := coordinateClient.GetCoordinates(c.UserContext(), zipcode)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		
		forecasRequest := &ForecastRequest{
			Latitude: coordinates.Latitude,
			Longitude: coordinates.Longitude,
			StartDate: c.Query("start_date", ""),
			EndDate:   c.Query("end_date", ""),
			Units:     c.Query("units", ""),
			TempScale: c.Query("temp_scale"),
		}

		forecastClient := NewForecastClient(client, config)
		forecast, err := forecastClient.GetForecast(c.UserContext(), forecasRequest)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(forecast)
	}
}
