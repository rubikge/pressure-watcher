package api

import (
	"pressure-watcher-app/weather/service"

	"github.com/gofiber/fiber/v3"
)

type WeatherController struct {
	weatherService *service.WeatherService
}

func NewWeatherController(weatherService *service.WeatherService) *WeatherController {
	return &WeatherController{weatherService: weatherService}
}

func (wc *WeatherController) GetLatestWeatherHandler(c fiber.Ctx) error {
	significantChange, err := wc.weatherService.IsSignificantChange()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve weather data",
		})
	}
	return c.JSON(fiber.Map{
		"significantChange": significantChange,
	})
}
