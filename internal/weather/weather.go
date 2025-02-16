package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WeatherResponse struct {
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

func GetWeather(longitude float64, latitude float64, time time.Time) (string, error) {
	weather, err := queryWeatherApi(longitude, latitude)
	if err != nil {
		return "", errors.New("could not get temperature")
	}

	temperature, err := getTemperature(weather, "08:00")
	if err != nil {
		return "", errors.New("could not get temperature")
	}

	response := fmt.Sprintf("Temperature at 8:00 AM MT: %.2f°F\n", temperature)
	return response, nil
}

func queryWeatherApi(longitude float64, latitude float64) (WeatherResponse, error) {
	timezone := "America/Denver"
	forecastDays := 1
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&hourly=temperature_2m&timezone=%s&forecast_days=%d&temperature_unit=fahrenheit",
		latitude, longitude, timezone, forecastDays,
	)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return WeatherResponse{}, err
	}

	return weather, nil
}

func getTemperature(weather WeatherResponse, targetTime string) (float64, error) {
	for i, time := range weather.Hourly.Time {
		if time[len(time)-5:] == targetTime { // Match "HH:MM" to time string
			return weather.Hourly.Temperature2m[i], nil
		}
	}

	return 0, errors.New("could not find time in response")
}
