package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"pressure-watcher/db"
	"pressure-watcher/weather"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Pressure Watcher Service started...")

	database, err := db.InitDB()
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	minPressure, maxPressure, significantChange, err := db.GetLastPressureLog(database)
	if err == nil {
		fmt.Printf("Min Pressure: %.2f hPa, Max Pressure: %.2f hPa\n", minPressure, maxPressure)

		if significantChange {
			fmt.Println("⚠ The pressure will significantly fluctuate in the next 24 hours!")
		}
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	pressureThreshold := 10.0

	FetchAndSavePressureData(database, pressureThreshold)

	for range ticker.C {
		FetchAndSavePressureData(database, pressureThreshold)
	}
}

func FetchAndSavePressureData(database *sql.DB, pressureThreshold float64) {
	pressures, err := weather.GetPressureForecast()
	if err != nil {
		log.Printf("Error fetching forecast: %v\n", err)
		return
	}

	if len(pressures) == 0 {
		log.Printf("No pressure data available")
		return
	}

	minPressure := pressures[0]
	maxPressure := pressures[0]

	for _, p := range pressures {
		if p < minPressure {
			minPressure = p
		}
		if p > maxPressure {
			maxPressure = p
		}
	}

	significantChange := maxPressure-minPressure > pressureThreshold

	db.SavePressureLog(database, minPressure, maxPressure, significantChange)
}
