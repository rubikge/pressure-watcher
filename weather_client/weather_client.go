package weather_client

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type WeatherClient struct {
	Url string
}

type PressureForecast struct {
	MinPressure float64
	MaxPressure float64
}

type ForecastResponse struct {
	List []struct {
		Main struct {
			Pressure float64 `json:"pressure"`
		} `json:"main"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
}

func NewWeatherClient(url string) *WeatherClient {
	return &WeatherClient{Url: url}
}

func (wc *WeatherClient) GetPressureForecast() (PressureForecast, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(wc.Url)
	if err != nil {
		return PressureForecast{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return PressureForecast{}, errors.New("received non-OK HTTP status")
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(res.Body).Decode(&forecast); err != nil {
		return PressureForecast{}, err
	}

	now := time.Now()
	var pressures []float64

	for _, item := range forecast.List {
		t, parseErr := time.Parse("2006-01-02 15:04:05", item.DtTxt)
		if parseErr != nil {
			continue
		}
		if t.After(now) && t.Before(now.Add(24*time.Hour)) {
			pressures = append(pressures, item.Main.Pressure)
		}
	}

	if len(pressures) == 0 {
		return PressureForecast{}, errors.New("no pressure data available for the next 24 hours")
	}

	minPressure, maxPressure := pressures[0], pressures[0]
	for _, p := range pressures[1:] {
		if p < minPressure {
			minPressure = p
		}
		if p > maxPressure {
			maxPressure = p
		}
	}

	return PressureForecast{
		MinPressure: minPressure,
		MaxPressure: maxPressure,
	}, nil
}
