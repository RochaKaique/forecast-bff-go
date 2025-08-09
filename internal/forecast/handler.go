package forecast

import (
	"log"

	"github.com/RochaKaique/forecastgo/internal/coordinates"
	"github.com/gofiber/fiber/v2"
)

func GetForecast(http *coordinates.CoordinatesClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		zipcode := c.Params("zipcode")
		resp, err := http.GetCoordinates(c.UserContext(), zipcode)
		if err != nil {
			log.Fatal(err)
		}
		return c.Status(fiber.StatusAccepted).JSON(resp)
	}
}
