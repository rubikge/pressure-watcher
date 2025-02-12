package models

import (
	"fmt"
	"pressure-watcher-app/weather_client"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PressureLog struct {
	ID          uint      `gorm:"primaryKey"`
	Timestamp   time.Time `gorm:"type:timestamp with time zone"`
	MinPressure float32   `gorm:"type:real"`
	MaxPressure float32   `gorm:"type:real"`
}

type Database struct {
	*gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&PressureLog{}); err != nil {
		return nil, fmt.Errorf("migration failed: %v", err)
	}

	return &Database{db}, nil
}

func (db *Database) SavePressure(pressureForecast *weather_client.PressureForecast) error {
	pressureLog := PressureLog{
		MinPressure: float32(pressureForecast.MinPressure),
		MaxPressure: float32(pressureForecast.MaxPressure),
		Timestamp:   time.Now(),
	}

	result := db.Create(&pressureLog)
	if result.Error != nil {
		return fmt.Errorf("failed to save pressure log: %v", result.Error)
	}

	return nil
}

func (db *Database) GetLatestPressureForecast() (*weather_client.PressureForecast, error) {
	var pressureLog PressureLog
	if err := db.Order("timestamp desc").First(&pressureLog).Error; err != nil {
		return nil, err
	}

	pressureForecast := weather_client.PressureForecast{
		MinPressure: float64(pressureLog.MinPressure),
		MaxPressure: float64(pressureLog.MaxPressure),
	}

	return &pressureForecast, nil
}
