package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

type viewPayload struct {
	ClientId     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Subscription struct {
	ID            int    `json:"id"`
	ResourceState int    `json:"resource_state"`
	ApplicationId int    `json:"application_id"`
	CallbackUrl   string `json:"callback_url"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
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

	u, err := url.Parse("https://www.strava.com/api/v3/push_subscriptions")
	if err != nil {
		log.Println("Failed to parse URL:", err)
		os.Exit(1)
	}

	q := u.Query()
	q.Set("client_id", strconv.Itoa(myStravaId.ClientId))
	q.Set("client_secret", myStravaId.ClientSecret)
	u.RawQuery = q.Encode()

	respGet, err := http.Get(u.String())
	if err != nil {
		log.Println("GET request failed:", err)
		os.Exit(1)
	}
	defer respGet.Body.Close()

	bodyGet, err := io.ReadAll(respGet.Body)
	if err != nil {
		log.Println("Failed to read GET response body:", err)
		os.Exit(1)
	}

	var subscriptions []Subscription
	err = json.Unmarshal(bodyGet, &subscriptions)
	if err != nil {
		log.Println("Failed to unmarshal subscriptions:", err)
		os.Exit(1)
	}

	log.Println("Parsed Subscriptions:", subscriptions)

	// delete subscription
	if len(subscriptions) > 0 {
		subscriptionId := subscriptions[0].ID
		deleteUrl, err := url.Parse("https://www.strava.com/api/v3/push_subscriptions/" + strconv.Itoa(subscriptionId))
		if err != nil {
			log.Println("Failed to parse DELETE URL:", err)
			os.Exit(1)
		}

		q := deleteUrl.Query()
		q.Set("client_id", strconv.Itoa(myStravaId.ClientId))
		q.Set("client_secret", myStravaId.ClientSecret)
		deleteUrl.RawQuery = q.Encode()

		req, err := http.NewRequest("DELETE", deleteUrl.String(), nil)
		if err != nil {
			log.Println("Failed to create DELETE request:", err)
			os.Exit(1)
		}

		deleteResp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Failed to make DELETE request:", err)
			os.Exit(1)
		}
		defer deleteResp.Body.Close()

		deleteBody, err := io.ReadAll(deleteResp.Body)
		if err != nil {
			log.Println("Failed to read DELETE response body:", err)
			os.Exit(1)
		}

		log.Println("DELETE Response Status:", deleteResp.Status)
		log.Println("DELETE Response Body:", string(deleteBody))
	}

	// subscribe
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

	resp, err := http.Post("https://www.strava.com/api/v3/push_subscriptions",
		"application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("Failed to make POST request:", err)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		os.Exit(1)
	}

	log.Println("Response Status:", resp.Status)
	log.Println("Response Body:", string(body))

}
