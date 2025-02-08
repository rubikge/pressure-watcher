package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiKey = "9b74e72982ee57e583716174905a4ed5"
const city = "Tbilisi"
const apiURL = "https://api.openweathermap.org/data/2.5/forecast"

type ForecastResponse struct {
	List []struct {
		Main struct {
			Pressure float64 `json:"pressure"`
		} `json:"main"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
}

func GetPressureForecast() ([]float64, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", apiURL, city, apiKey)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var forecast ForecastResponse
	if err := json.NewDecoder(res.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	pressures := []float64{}
	now := time.Now()

	for _, item := range forecast.List {
		t, _ := time.Parse("2006-01-02 15:04:05", item.DtTxt)
		if t.After(now) && t.Before(now.Add(24*time.Hour)) {
			pressures = append(pressures, item.Main.Pressure)
		}
	}

	return pressures, nil
}
