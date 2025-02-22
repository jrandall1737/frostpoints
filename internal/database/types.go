package database

type UserToken struct {
	AthleteId    int64  `json:"athlete_id" bson:"athlete_id"`
	AthleteName  string `json:"athlete_name" bson:"athlete_name"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	AccessToken  string `json:"access_token" bson:"access_token"`
	ExpiresAt    int64  `json:"expires_at" bson:"expires_at"`
	ExpiresIn    int    `json:"expires_in" bson:"expires_in"`
}
