package forecast

import (
	"github.com/RochaKaique/forecastgo/internal/coordinates"
	"github.com/gofiber/fiber/v2"
)

func GetForecast(http *coordinates.CoordinatesClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		zipcode := c.Params("zipcode")
		coordinates, err := http.GetCoordinates(c.UserContext(), zipcode)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(coordinates)
	}
}
