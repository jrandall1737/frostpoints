package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jrandall1737/frostpoints/internal/app"
	"github.com/jrandall1737/frostpoints/pkg/strava"
)

func main() {
	var port int // port of local demo server
	var myStravaConfig strava.StravaConfig

	// setup the credentials for your app
	// These need to be set to reflect your application
	// and can be found at https://www.strava.com/settings/api
	flag.IntVar(&myStravaConfig.ClientId, "id", 0, "Strava Client ID")
	flag.StringVar(&myStravaConfig.ClientSecret, "secret", "", "Strava Client Secret")
	flag.StringVar(&myStravaConfig.CallbackUrl, "callback", "localhost", "Strava Callback URL")
	flag.IntVar(&port, "port", 3009, "Strava Client Secret")

	flag.Parse()

	if myStravaConfig.ClientId == 0 || myStravaConfig.ClientSecret == "" {
		fmt.Println("\nPlease provide your application's client_id and client_secret.")
		fmt.Println("For example: go run oauth_example.go -id=9 -secret=longrandomsecret")
		fmt.Println(" ")

		flag.PrintDefaults()
		os.Exit(1)
	}

	if myStravaConfig.CallbackUrl == "localhost" {
		myStravaConfig.CallbackUrl = fmt.Sprintf("http://localhost:%d", port)
	}

	app.StartApp(port, myStravaConfig)

	// weatherAtTime, err := weather.GetWeather(40.7128, 74.0060, time.Now())
	// weatherAtTime, err := weather.GetWeather(40.5853, 105.084, time.Now())

	// if err != nil {
	// 	fmt.Println(weatherAtTime)
	// }
}
