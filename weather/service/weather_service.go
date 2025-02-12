package service

import "pressure-watcher-app/models"

type WeatherServiceInterface interface {
	IsSignificantChange() (bool, error)
}

type WeatherService struct {
	db *models.Database
}

func NewWeatherService(db *models.Database) *WeatherService {
	return &WeatherService{db: db}
}

func (ws *WeatherService) IsSignificantChange() (bool, error) {
	pressureThreshold := 10.0

	pressureForecast, err := ws.db.GetLatestPressureForecast()
	if err != nil {
		return false, err
	}

	significantChange := pressureForecast.MaxPressure-pressureForecast.MinPressure > pressureThreshold

	return significantChange, nil
}
