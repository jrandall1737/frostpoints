package weather

import (
	"testing"
	"time"
)

func TestTomorrowIo(t *testing.T) {
	tomorrowIo := newTomorrowIo()

	temperature, err := tomorrowIo.GetTemperature(40.5853, -105.084, time.Now().Add(-10*time.Minute))

	if err != nil {
		t.Errorf("GetWeather() error = %v", err)
	}

	if temperature == "" {
		t.Errorf("GetWeather() temperature = %v", temperature)
	}
}

func TestRoundTime(t *testing.T) {
	tomorrowIo := newTomorrowIo()

	roundedTime := tomorrowIo.roundTimeToNearestHour(time.Date(2021, 1, 1, 1, 30, 0, 0, time.UTC))

	if roundedTime.Hour() != 2 {
		t.Errorf("roundTimeToNearestHour() = %v", roundedTime)
	}
}
