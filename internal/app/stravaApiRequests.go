package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jrandall1737/frostpoints/pkg/strava"
)

var stravaUrl = "https://www.strava.com/api/v3"

func getActivity(token string, objectId int64) (*strava.ActivitySummary, error) {
	client := &http.Client{}

	activityUrl := stravaUrl + "/activities/" + strconv.FormatInt(objectId, 10) + "?include_all_efforts=false"

	req, err := http.NewRequest("GET", activityUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get activity: %d", resp.StatusCode)
	}

	var activity strava.ActivitySummary
	err = json.NewDecoder(resp.Body).Decode(&activity)
	if err != nil {
		return nil, err
	}

	log.Printf("Activity: %v\n", activity)

	return &activity, nil
}

func updateActivity(token string, objectId int64, activity strava.ActivitySummary, temperature int, frostPoints int) error {
	client := &http.Client{}
	activityUrl := stravaUrl + "/activities/" + strconv.FormatInt(objectId, 10)

	// achievementString := fmt.Sprintf("This run earned %d Frost Points!", frostPoints)
	achievementString := fmt.Sprintf("❄ Frost Points: %d (%d°F)", frostPoints, temperature)

	var newDescription string
	if activity.Description == "" {
		newDescription = achievementString
	} else {
		newDescription = activity.Description + "\n" + achievementString
	}

	update := strava.UpdatableActivity{
		Description: newDescription,
	}

	bodyPayload, err := json.Marshal(update)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", activityUrl, bytes.NewBuffer(bodyPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update activity: %d", resp.StatusCode)
	}

	log.Println("Activity updated")

	return nil
}
