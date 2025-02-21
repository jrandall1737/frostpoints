package strava

type StravaConfig struct {
	ClientId     int
	ClientSecret string
	CallbackUrl  string
}

type WebhookEvent struct {
	ObjectType     string            `json:"object_type"`     // Either "activity" or "athlete"
	ObjectID       int64             `json:"object_id"`       // ID of the activity or athlete
	AspectType     string            `json:"aspect_type"`     // Always "create," "update," or "delete"
	Updates        map[string]string `json:"updates"`         // Contains updates about the event
	OwnerID        int64             `json:"owner_id"`        // Athlete's ID
	SubscriptionID int               `json:"subscription_id"` // Push subscription ID receiving this event
	EventTime      int64             `json:"event_time"`      // The time the event occurred
}
