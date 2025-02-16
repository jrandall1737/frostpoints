package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type tomorrowIo struct {
	baseUrl  string
	apiToken string
}

func newTomorrowIo() tomorrowIo {
	return tomorrowIo{
		"https://api.tomorrow.io/v4",
		"8B7SYSJrzjcNTt0luynAGARNS3NtvdTG",
	}
}

func (t tomorrowIo) GetTemperature(latitude float64, longitude float64, activityTime time.Time) (string, error) {
	roundedTime := t.roundTimeToNearestHour(activityTime)
	hourlyWeather, err := t.requestHourlyWeather(latitude, longitude)
	if err != nil {
		return "", err
	}

	temperature, err := t.getTemperatureAtTime(hourlyWeather, roundedTime)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%f", temperature), nil
}

func (t tomorrowIo) requestHourlyWeather(latitude float64, longitude float64) ([]HourlyData, error) {
	params := url.Values{}
	params.Add("apikey", t.apiToken)
	params.Add("location", fmt.Sprintf("%f, %f", latitude, longitude))
	params.Add("units", "imperial")

	url := fmt.Sprintf("%s/weather/history/recent?%s", t.baseUrl, params.Encode())

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	resp, _ := http.DefaultClient.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var weather WeatherData
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	return weather.Timelines.Hourly, nil
}

func (t tomorrowIo) roundTimeToNearestHour(activityTime time.Time) time.Time {
	return activityTime.Round(time.Hour)
}

func (t tomorrowIo) getTemperatureAtTime(hourlyWeather []HourlyData, targetTime time.Time) (float64, error) {
	for _, data := range hourlyWeather {
		if data.Time.Equal(targetTime) {
			return data.Values.TemperatureApparent, nil
		}
	}

	return 0, fmt.Errorf("no data found for time %v", targetTime)
}
