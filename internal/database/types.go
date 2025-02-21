package database

type UserToken struct {
	AthleteId    int64  `json:"athlete_id"`
	AthleteName  string `json:"athlete_name"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
}
