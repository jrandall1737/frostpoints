package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jrandall1737/frostpoints/internal/auth"
	"github.com/jrandall1737/frostpoints/internal/database"
	"github.com/jrandall1737/frostpoints/internal/weather"
	"github.com/jrandall1737/frostpoints/pkg/strava"
)

type WebhookHandler struct {
	db   database.Database
	auth *auth.StravaAuth
}

func NewStravaWebhookHandler(db database.Database, auth *auth.StravaAuth) *WebhookHandler {
	return &WebhookHandler{db: db, auth: auth}
}

func (wh *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook received!")

	switch r.Method {
	case http.MethodGet:
		wh.handleWebhookGet(w, r)
	case http.MethodPost:
		wh.handleWebhookPost(w, r)
	default:
		log.Println("method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Webhook GET handler (for verification)
func (wh *WebhookHandler) handleWebhookGet(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	mode := query.Get("hub.mode")
	token := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")
	log.Printf("Webhook GET received! mode=%s, token=%s, challenge=%s\n", mode, token, challenge)

	// Check if mode and token are provided
	if mode != "" && token != "" {
		// Verify the token and mode
		if mode == "subscribe" && token == VERIFY_TOKEN {
			log.Println("WEBHOOK_VERIFIED")

			// Respond with the challenge
			response := map[string]string{"hub.challenge": challenge}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		// Return 403 Forbidden if token is invalid
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Return 400 Bad Request if parameters are missing
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

// Webhook POST handler
func (wh *WebhookHandler) handleWebhookPost(w http.ResponseWriter, r *http.Request) {
	// Log query parameters and body
	log.Println("Webhook POST received!")

	// Respond with 200 OK
	w.WriteHeader(http.StatusOK)

	// Read request body
	defer r.Body.Close()
	var event strava.WebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err == nil {
		log.Printf("Request Body: %#v\n", event)
		wh.handleEvent(event)
	} else {
		log.Println("Failed to parse JSON body:", err)
	}
}

func (wh *WebhookHandler) handleEvent(event strava.WebhookEvent) {
	// Don't handle athlete events
	// if event.AspectType != "create" {
	// 	return
	// }

	// get the athlete token
	userToken := wh.db.FindTokenById(event.OwnerID)
	if userToken == nil {
		log.Println("Failed to find user token")
		return
	}

	// swap the refresh token for an auth token
	authToken, err := wh.auth.RefreshToken(userToken.RefreshToken)
	if err != nil {
		log.Println("Failed to refresh token:", err)
		return
	}

	// get the activity
	activity, err := getActivity(authToken.AccessToken, event.ObjectID)
	if err != nil {
		log.Println("Failed to get activity:", err)
		return
	}

	log.Printf("Activity: %v\n", activity)

	// Don't update if the activity already has "Frost Points" in the description
	if strings.Contains(activity.Description, "Frost Points") {
		log.Println("Activity already has frost points, skipping")
		return
	}

	if activity.Manual {
		log.Println("Activity is manual, skipping")
		return
	}

	if activity.Type != "Run" && activity.Type != "TrailRun" && activity.Type != "Walk" {
		log.Println("Activity is not a run or walk, skipping")
		return
	}

	temperatureAtEvent, err := weather.GetWeather(activity.StartLocation[0], activity.StartLocation[1], activity.StartDateLocal)
	if err != nil {
		log.Println("Failed to get weather:", err)
		return
	}

	if *temperatureAtEvent >= 32.0 {
		log.Println("Temperature too high, skipping", err)
		return
	}

	frostPoints, err := getFrostPoints(*temperatureAtEvent, *activity)
	if err != nil {
		log.Println("Failed to get frost points:", err)
		return
	}

	err = updateActivity(authToken.AccessToken, event.ObjectID, *activity, int(*temperatureAtEvent), *frostPoints)
	if err != nil {
		log.Println("Failed to update:", err)
		return
	}
}

func getFrostPoints(temperatureAtEvent float64, activity strava.ActivitySummary) (*int, error) {
	frostPoints := int((32.0 - temperatureAtEvent) * activity.Distance * 0.621371 / 1000)

	log.Println("Frost Points:", frostPoints)

	return &frostPoints, nil
}
