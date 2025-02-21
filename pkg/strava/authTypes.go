package strava

type AuthResponse struct {
	TokenType    string         `json:"token_type"`
	ExpiresAt    int            `json:"expires_at"`
	ExpiresIn    int            `json:"expires_in"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	Athlete      AthleteSummary `json:"athlete"`
}

type RefreshRequest struct {
	ClientId     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type RefreshResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
