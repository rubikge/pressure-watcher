package scheduler

import (
	"log"
	"pressure-watcher-app/models"
	"pressure-watcher-app/weather_client"
	"time"
)

func Start(client *weather_client.WeatherClient, db *models.Database) {
	go func() {
		for {
			processPressureData(client, db)
			time.Sleep(1 * time.Hour)
		}
	}()
}

func processPressureData(client *weather_client.WeatherClient, db *models.Database) {
	pressureData, err := client.GetPressureForecast()
	if err != nil {
		log.Printf("Error fetching pressure data: %v", err)
		return
	}

	err = db.SavePressure(&pressureData)
	if err != nil {
		log.Printf("Error saving to database: %v", err)
		return
	}

	log.Println("Data saved")
}
