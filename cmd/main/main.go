package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/jrandall1737/frostpoints/internal/app"
	"github.com/jrandall1737/frostpoints/pkg/strava"
)

var port string // port of local demo server
var myStravaConfig strava.StravaConfig
var dbConnectionString string

func main() {
	// setup the credentials for your app
	// These need to be set to reflect your application
	// and can be found at https://www.strava.com/settings/api
	flag.IntVar(&myStravaConfig.ClientId, "id", 0, "Strava Client ID")
	flag.StringVar(&myStravaConfig.ClientSecret, "secret", "", "Strava Client Secret")
	flag.StringVar(&myStravaConfig.CallbackUrl, "callback", "localhost", "Strava Callback URL")
	flag.StringVar(&port, "port", "3009", "Strava Client Secret")
	flag.StringVar(&dbConnectionString, "db", "", "Database connection string")

	flag.Parse()

	readEnvironmentVariables()

	if myStravaConfig.ClientId == 0 || myStravaConfig.ClientSecret == "" {
		fmt.Println("\nPlease provide your application's client_id and client_secret.")
		fmt.Println("For example: go run oauth_example.go -id=9 -secret=longrandomsecret")
		fmt.Println(" ")

		flag.PrintDefaults()
		os.Exit(1)
	}

	if myStravaConfig.CallbackUrl == "localhost" {
		myStravaConfig.CallbackUrl = "http://localhost:" + port
	}

	if dbConnectionString == "" {
		fmt.Println("\nPlease provide a connection string for the database.")
		os.Exit(1)
	}

	app.StartApp(port, myStravaConfig, dbConnectionString)

}

func readEnvironmentVariables() {
	// read environment variables
	value, exists := os.LookupEnv("PORT")
	if exists {
		fmt.Println("Using PORT environment variable")
		port = value
	}

	value, exists = os.LookupEnv("STRAVA_SECRET")
	if exists {
		fmt.Println("Using STRAVA_SECRET environment variable")
		myStravaConfig.ClientSecret = value
	}

	value, exists = os.LookupEnv("STRAVA_ID")
	if exists {
		fmt.Println("Using STRAVA_ID environment variable")
		id, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("STRAVA_ID must be an integer")
			os.Exit(1)
		}
		myStravaConfig.ClientId = id
	}

	value, exists = os.LookupEnv("DB_CONNECTION_STRING")
	if exists {
		fmt.Println("Using DB_CONNECTION_STRING environment variable")
		dbConnectionString = value
	}
}
