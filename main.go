package main

import (
	"log"
	"pressure-watcher-app/config"
	"pressure-watcher-app/handlers"
	"pressure-watcher-app/models"
	"pressure-watcher-app/scheduler"
	"pressure-watcher-app/weather_client"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()

	db, err := models.NewDatabase(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	wc := weather_client.NewWeatherClient(cfg.WeatherApiUrl)

	scheduler.Start(wc, db)

	app := fiber.New()
	app.Get("/pressure-check", handlers.GetLatestWeatherHandler(db))
	app.Listen(":8080")
}
