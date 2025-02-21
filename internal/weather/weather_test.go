package weather

import (
	"strings"
	"testing"
	"time"
)

func TestGetWeather(t *testing.T) {
	type args struct {
		longitude float64
		latitude  float64
		time      time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Test GetWeather",
			args:    args{40.5853, 105.084, time.Now()},
			want:    "Temperature at 8:00 AM MT:",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWeather(tt.args.longitude, tt.args.latitude, tt.args.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(got, tt.want) {
				t.Errorf("GetWeather() = %v, want %v", got, tt.want)
			}
		})
	}
}
