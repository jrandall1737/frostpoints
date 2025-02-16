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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var oauth2Config *oauth2.Config

type StravaTokenRequest struct {
	ClientId     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
}

type StravaTokenResponse struct {
	ExpiresAt    int     `json:"expires_at"`
	ExpiresIn    int     `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	AccessToken  string  `json:"access_token"`
	Athlete      Athlete `json:"athlete"`
}

type StravaRefreshTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func SetOauthConfig(clientId int, clientSecret string, port int) {
	oauth2Config = &oauth2.Config{
		ClientID:     strconv.Itoa(clientId),
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("http://localhost:%d/callback", port),
		Scopes:       []string{"read,activity:write"},
		Endpoint:     endpoints.Strava,
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	promptOpt := oauth2.SetAuthURLParam("approval_prompt", "auto")
	url := oauth2Config.AuthCodeURL("state-token", promptOpt)
	http.Redirect(w, r, url, http.StatusFound)
}

func exchangeCodeForToken(code string) (*StravaTokenResponse, error) {
	// Prepare the form tokenExchangeRequestUrl for the POST request
	tokenExchangeRequestUrl := url.Values{}
	tokenExchangeRequestUrl.Set("client_id", oauth2Config.ClientID)
	tokenExchangeRequestUrl.Set("client_secret", oauth2Config.ClientSecret)
	tokenExchangeRequestUrl.Set("code", code)
	tokenExchangeRequestUrl.Set("grant_type", "authorization_code")

	// Make the POST request
	resp, err := http.Post(endpoints.Strava.TokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(tokenExchangeRequestUrl.Encode()))
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

func renderTemplate(w http.ResponseWriter, tmpl string, athlete Athlete) error {
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
	<p><strong>Firstname:</strong> {{.Firstname}}</p>
	<p><strong>Lastname:</strong> {{.Lastname}}</p>
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

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in callback", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for a token
	stravaToken, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Not able to exchange code", http.StatusBadRequest)
	}

	// Store the data securely (e.g., in a database or encrypted file)
	log.Printf("Received: %+v", *stravaToken)

	err = renderTemplate(w, "athlete", stravaToken.Athlete)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}
