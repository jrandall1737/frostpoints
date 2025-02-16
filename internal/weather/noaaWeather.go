package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiToken = "EVvspZChBnQDZYENueLLHcfAXzrGcBpF"
	baseURL  = "https://www.ncei.noaa.gov/cdo-web/api/v2"
)

type StationResponse struct {
	Results []struct {
		ID string `json:"id"`
	} `json:"results"`
}

type WeatherDataCollection struct {
	Results []struct {
		Date  string  `json:"date"`
		Value float64 `json:"value"`
	} `json:"results"`
}

type APIResponse struct {
	Results  []WeatherStationResponse `json:"results"`
	Metadata struct {
		ResultSet struct {
			Limit  int `json:"limit"`
			Count  int `json:"count"`
			Offset int `json:"offset"`
		} `json:"resultset"`
	} `json:"metadata"`
}

type WeatherStationResponse struct {
	Elevation     float64 `json:"elevation"`
	MinDate       string  `json:"mindate"`
	MaxDate       string  `json:"maxdate"`
	Latitude      float64 `json:"latitude"`
	Name          string  `json:"name"`
	DataCoverage  float64 `json:"datacoverage"`
	ID            string  `json:"id"`
	ElevationUnit string  `json:"elevationUnit"`
	Longitude     float64 `json:"longitude"`
}

func GetNoaaWeather(longitude float64, latitude float64) error {
	stationId, err := getNearestStation(latitude, longitude)
	if err != nil {
		return err
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	targetTime := "08:00"
	return getWeatherData(stationId, yesterday, targetTime)
}

// getNearestStation finds the closest station to the given coordinates
func getNearestStation(longitude float64, latitude float64) (string, error) {
	boxSize := 10.0
	url := fmt.Sprintf("%s/stations?datasetid=GHCND&limit=5&extent=%f,%f,%f,%f", baseURL, latitude-boxSize, longitude-boxSize, latitude+boxSize, longitude+boxSize)
	// url := fmt.Sprintf("%s/stations?datasetid=GHCND&limit=1&sortfield=distance&sortorder=asc&extent=%f,%f,%f,%f", baseURL, latitude, longitude, latitude, longitude)
	// url := fmt.Sprintf("%s/stations/COOP:010008", baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("token", apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var stationResp APIResponse
	if err := json.Unmarshal(body, &stationResp); err != nil {
		return "", err
	}

	return stationResp.Results[0].ID, nil
}

func getWeatherData(stationId, date, targetTime string) error {
	url := fmt.Sprintf("%s/data?datasetid=GHCND&stationid=%s&startdate=%s&enddate=%s&units=standard&limit=1000", baseURL, stationId, date, date)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("token", apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var weatherData WeatherDataCollection
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return err
	}

	// Filter for relevant time and temperature
	for _, data := range weatherData.Results {
		if data.Date[:10] == date {
			fmt.Printf("Weather at %s: %.1f°F\n", targetTime, data.Value/10) // Convert to degrees
			return nil
		}
	}

	return fmt.Errorf("no weather data found")
}
