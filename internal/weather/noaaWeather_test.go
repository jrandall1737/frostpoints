package weather

import "testing"

func TestGetNoaaWeather(t *testing.T) {
	err := GetNoaaWeather(105.084, 40.5853)

	if err != nil {
		t.Errorf("GetNoaaWeather() error = %v", err)
	}
}
