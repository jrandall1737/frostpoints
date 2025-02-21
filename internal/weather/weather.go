package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type WeatherResponse struct {
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

func GetWeather(latitude float64, longitude float64, startTime time.Time) (*float64, error) {
	weather, err := queryWeatherApi(longitude, latitude)
	if err != nil {
		return nil, errors.New("could not get temperature")
	}

	closestTime, err := getClosestTimeToEvent(startTime, weather.Hourly.Time)
	if err != nil {
		return nil, errors.New("could not get closest time")
	}

	temperature, err := getTemperature(weather, closestTime)
	if err != nil {
		return nil, errors.New("could not get temperature")
	}

	log.Printf("Temperature at %s was %f\n", startTime.Format(time.RFC3339), temperature)
	return &temperature, nil
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

func getClosestTimeToEvent(eventTime time.Time, weatherTimes []string) (string, error) {
	// Define a layout for parsing (only HH:MM)
	layout := "2006-01-02T15:04"

	// Track the closest time
	var closestTime time.Time
	minDiff := time.Duration(1<<63 - 1) // Set to max possible duration initially

	for _, ts := range weatherTimes {
		// Parse the time string into a time.Time object
		parsedTime, err := time.Parse(layout, ts)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}

		// Compute absolute difference
		diff := eventTime.Sub(parsedTime)
		if diff < 0 {
			diff = -diff // Get absolute value
		}

		// Update closest time if this one is nearer
		if diff < minDiff {
			minDiff = diff
			closestTime = parsedTime
		}
	}

	// Output result
	fmt.Println("Current Time:", eventTime.Format("15:04"))
	fmt.Println("Closest Time:", closestTime.Format("15:04"))

	return closestTime.Format("15:04"), nil
}

func getTemperature(weather WeatherResponse, targetTime string) (float64, error) {
	for i, time := range weather.Hourly.Time {
		if time[len(time)-5:] == targetTime { // Match "HH:MM" to time string
			return weather.Hourly.Temperature2m[i], nil
		}
	}

	return 0, errors.New("could not find time in response")
}
