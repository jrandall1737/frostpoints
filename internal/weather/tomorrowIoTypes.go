package weather

import "time"

type WeatherData struct {
	Timelines Timelines `json:"timelines"`
	Location  Location  `json:"location"`
}

type Timelines struct {
	Hourly []HourlyData `json:"hourly"`
	Daily  []DailyData  `json:"daily"`
}

type HourlyData struct {
	Time   time.Time    `json:"time"`
	Values HourlyValues `json:"values"`
}

type DailyData struct {
	Time   time.Time   `json:"time"`
	Values DailyValues `json:"values"`
}

type HourlyValues struct {
	CloudBase                *float64 `json:"cloudBase"`
	CloudCeiling             *float64 `json:"cloudCeiling"`
	CloudCover               int      `json:"cloudCover"`
	DewPoint                 float64  `json:"dewPoint"`
	Evapotranspiration       float64  `json:"evapotranspiration"`
	FreezingRainIntensity    float64  `json:"freezingRainIntensity"`
	HailProbability          float64  `json:"hailProbability"`
	HailSize                 float64  `json:"hailSize"`
	Humidity                 int      `json:"humidity"`
	IceAccumulation          float64  `json:"iceAccumulation"`
	IceAccumulationLwe       float64  `json:"iceAccumulationLwe"`
	PrecipitationProbability float64  `json:"precipitationProbability"`
	PressureSurfaceLevel     float64  `json:"pressureSurfaceLevel"`
	RainAccumulation         float64  `json:"rainAccumulation"`
	RainAccumulationLwe      float64  `json:"rainAccumulationLwe"`
	RainIntensity            float64  `json:"rainIntensity"`
	SleetAccumulation        float64  `json:"sleetAccumulation"`
	SleetAccumulationLwe     float64  `json:"sleetAccumulationLwe"`
	SleetIntensity           float64  `json:"sleetIntensity"`
	SnowAccumulation         float64  `json:"snowAccumulation"`
	SnowAccumulationLwe      float64  `json:"snowAccumulationLwe"`
	SnowDepth                float64  `json:"snowDepth"`
	SnowIntensity            float64  `json:"snowIntensity"`
	Temperature              float64  `json:"temperature"`
	TemperatureApparent      float64  `json:"temperatureApparent"`
	UVHealthConcern          int      `json:"uvHealthConcern"`
	UVIndex                  int      `json:"uvIndex"`
	Visibility               float64  `json:"visibility"`
	WeatherCode              int      `json:"weatherCode"`
	WindDirection            float64  `json:"windDirection"`
	WindGust                 float64  `json:"windGust"`
	WindSpeed                float64  `json:"windSpeed"`
}

type DailyValues struct {
	CloudBaseAvg           float64 `json:"cloudBaseAvg"`
	CloudBaseMax           float64 `json:"cloudBaseMax"`
	CloudBaseMin           float64 `json:"cloudBaseMin"`
	SnowDepthAvg           float64 `json:"snowDepthAvg"`
	SnowDepthMax           float64 `json:"snowDepthMax"`
	SnowDepthMin           float64 `json:"snowDepthMin"`
	SnowDepthSum           float64 `json:"snowDepthSum"`
	TemperatureAvg         float64 `json:"temperatureAvg"`
	TemperatureMax         float64 `json:"temperatureMax"`
	TemperatureMin         float64 `json:"temperatureMin"`
	TemperatureApparentAvg float64 `json:"temperatureApparentAvg"`
	TemperatureApparentMax float64 `json:"temperatureApparentMax"`
	TemperatureApparentMin float64 `json:"temperatureApparentMin"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
