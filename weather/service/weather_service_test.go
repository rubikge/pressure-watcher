package service_test

import (
	"pressure-watcher-app/weather/service"
	"pressure-watcher-app/weather_client"
	"testing"
)

type MockDatabase struct {
	GetLatestPressureForecastFunc func() (*weather_client.PressureForecast, error)
}

func (m *MockDatabase) GetLatestPressureForecast() (*weather_client.PressureForecast, error) {
	return m.GetLatestPressureForecastFunc()
}

func TestIsSignificantChange(t *testing.T) {
	tests := []struct {
		name           string
		maxPressure    float64
		minPressure    float64
		expectedResult bool
	}{
		{"No significant change", 1015, 1010, false},
		{"Significant change", 1025, 1010, true},
		{"Exactly at threshold", 1020, 1010, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := &MockDatabase{
				GetLatestPressureForecastFunc: func() (*weather_client.PressureForecast, error) {
					return &weather_client.PressureForecast{
						MinPressure: tt.minPressure,
						MaxPressure: tt.maxPressure,
					}, nil
				},
			}

			weatherService := service.NewWeatherService(mockDb)

			result, _ := weatherService.IsSignificantChange()
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
