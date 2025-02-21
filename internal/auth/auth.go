package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jrandall1737/frostpoints/internal/database"
	"github.com/jrandall1737/frostpoints/pkg/strava"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type StravaAuth struct {
	oauth2Config *oauth2.Config
	db           *database.Database
}

type StravaTokenRequest struct {
	ClientId     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
}

type StravaTokenResponse struct {
	ExpiresAt    int64                 `json:"expires_at"`
	ExpiresIn    int                   `json:"expires_in"`
	RefreshToken string                `json:"refresh_token"`
	AccessToken  string                `json:"access_token"`
	Athlete      strava.AthleteSummary `json:"athlete"`
}

type StravaRefreshTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func NewStravaAuth(db *database.Database, config strava.StravaConfig) *StravaAuth {
	auth := &StravaAuth{db: db}
	auth.SetOauthConfig(config)
	return auth
}

// {"read,read_all,profile:read_all,profile:write,activity:read,activity:read_all,activity:write"},

func (s *StravaAuth) SetOauthConfig(config strava.StravaConfig) {
	s.oauth2Config = &oauth2.Config{
		ClientID:     strconv.Itoa(config.ClientId),
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/callback", config.CallbackUrl),
		Scopes:       []string{"activity:read_all,activity:write"},
		Endpoint:     endpoints.Strava,
	}
}

func (s *StravaAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	promptOpt := oauth2.SetAuthURLParam("approval_prompt", "auto")
	url := s.oauth2Config.AuthCodeURL("state-token", promptOpt)
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *StravaAuth) exchangeCodeForToken(code string) (*StravaTokenResponse, error) {
	// Prepare the form tokenExchangeRequestUrl for the POST request
	tokenExchangeRequestUrl := url.Values{}
	tokenExchangeRequestUrl.Set("client_id", s.oauth2Config.ClientID)
	tokenExchangeRequestUrl.Set("client_secret", s.oauth2Config.ClientSecret)
	tokenExchangeRequestUrl.Set("code", code)
	tokenExchangeRequestUrl.Set("grant_type", "authorization_code")

	// Make the POST request
	resp, err := http.Post(endpoints.Strava.TokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(tokenExchangeRequestUrl.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code for token: %d, %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into the TokenResponse struct
	var tokenResponse StravaTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	token := database.UserToken{
		AthleteId:    tokenResponse.Athlete.ID,
		AthleteName:  tokenResponse.Athlete.Firstname + " " + tokenResponse.Athlete.Lastname,
		RefreshToken: tokenResponse.RefreshToken,
		AccessToken:  tokenResponse.AccessToken,
		ExpiresAt:    tokenResponse.ExpiresAt,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}

	if token.AthleteId == 0 {
		return nil, fmt.Errorf("failed to get athlete ID: %s", body)
	}

	s.db.AddToken(token)

	return &tokenResponse, nil
}

func (s *StravaAuth) RefreshToken(refreshToken string) (*StravaTokenResponse, error) {
	// Prepare the form tokenExchangeRequestUrl for the POST request
	tokenExchangeRequestUrl := url.Values{}
	tokenExchangeRequestUrl.Set("client_id", s.oauth2Config.ClientID)
	tokenExchangeRequestUrl.Set("client_secret", s.oauth2Config.ClientSecret)
	tokenExchangeRequestUrl.Set("grant_type", "refresh_token")
	tokenExchangeRequestUrl.Set("refresh_token", refreshToken)

	// Make the POST request
	resp, err := http.Post("https://www.strava.com/api/v3/oauth/token?", "application/x-www-form-urlencoded", bytes.NewBufferString(tokenExchangeRequestUrl.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into the TokenResponse struct
	var tokenResponse StravaTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (s *StravaAuth) renderTemplate(w http.ResponseWriter, tmpl string, athlete strava.AthleteSummary) error {
	// Define the HTML template
	tmplContent := `
<!DOCTYPE html>
<html>
<head>
	<title>Successfully Registered!</title>
</head>
<body>
	<h1>Athlete Information</h1>
	<p><strong>ID:</strong> {{.ID}}</p>
	<p><strong>First Name:</strong> {{.Firstname}}</p>
	<p><strong>Last Name:</strong> {{.Lastname}}</p>
</body>
</html>
`
	tmplParsed, err := template.New(tmpl).Parse(tmplContent)
	if err != nil {
		return err
	}

	// Execute the template with the provided data
	return tmplParsed.Execute(w, athlete)
}

func (s *StravaAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in callback", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for a token
	stravaToken, err := s.exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Not able to exchange code. Please try again.", http.StatusBadRequest)
		log.Println("Error exchanging code for token:", err)
		return
	}

	// Store the data securely (e.g., in a database or encrypted file)
	log.Printf("Received: %+v", *stravaToken)

	err = s.renderTemplate(w, "athlete", stravaToken.Athlete)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}
