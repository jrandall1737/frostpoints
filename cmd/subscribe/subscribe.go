package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type StravaId struct {
	ClientId     int
	ClientSecret string
}

type subscribePayload struct {
	ClientId     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackUrl  string `json:"callback_url"`
	VerifyToken  string `json:"verify_token"`
}

const VERIFY_TOKEN = "STRAVA"

func main() {
	var port int
	var myStravaId StravaId

	flag.IntVar(&myStravaId.ClientId, "id", 0, "Strava Client ID")
	flag.StringVar(&myStravaId.ClientSecret, "secret", "", "Strava Client Secret")
	flag.IntVar(&port, "port", 3009, "Strava Client Secret")

	flag.Parse()

	if myStravaId.ClientId == 0 || myStravaId.ClientSecret == "" {
		fmt.Println("\nPlease provide your application's client_id and client_secret.")
		fmt.Println("For example: go run oauth_example.go -id=9 -secret=longrandomsecret")
		fmt.Println(" ")

		flag.PrintDefaults()
		os.Exit(1)
	}

	subscribePayload := subscribePayload{
		ClientId:     myStravaId.ClientId,
		ClientSecret: myStravaId.ClientSecret,
		CallbackUrl:  fmt.Sprintf("http://localhost:%d/webhook", port),
		VerifyToken:  VERIFY_TOKEN,
	}

	jsonData, err := json.Marshal(subscribePayload)
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		os.Exit(1)
	}

	http.Post("https://www.strava.com/api/v3/push_subscriptions",
		"application/json",
		bytes.NewBuffer(jsonData))
}
