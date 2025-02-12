package handlers

import (
	"pressure-watcher-app/models"

	"github.com/gofiber/fiber/v3"
)

func GetLatestWeatherHandler(db *models.Database) fiber.Handler {
	pressureThreshold := 10.0

	return func(c fiber.Ctx) error {
		pressureForecast, err := db.GetLatestPressureForecast()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve weather data",
			})
		}

		significantChange := pressureForecast.MaxPressure-pressureForecast.MinPressure > pressureThreshold

		return c.JSON(fiber.Map{
			"significantChange": significantChange,
		})
	}
}
