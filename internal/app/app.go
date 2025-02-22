package app

import (
	"fmt"
	"net/http"

	"github.com/jrandall1737/frostpoints/internal/auth"
	"github.com/jrandall1737/frostpoints/internal/database"
	"github.com/jrandall1737/frostpoints/pkg/strava"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../assets/stravaConfig.yaml ../../assets/stravaSwagger.json

const VERIFY_TOKEN = "STRAVA"

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<html><body>
		<h1>Welcome to Strava OAuth Example</h1>
		<a href="/login">
			<img src="./assets/button.png" alt="Login with Strava";">
		</a>
	</body></html>`)
}

func StartApp(port string, myStravaConfig strava.StravaConfig, dbConnectionString string) {
	db := database.NewMongoDatabase(dbConnectionString)
	defer db.Disconnect()
	stravaAuth := auth.NewStravaAuth(db, myStravaConfig)
	stravaHandler := NewStravaWebhookHandler(db, stravaAuth)

	// Serve any assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/login", stravaAuth.HandleLogin)
	http.HandleFunc("/callback", stravaAuth.HandleCallback)
	http.HandleFunc("/webhook", stravaHandler.HandleWebhook)

	fmt.Printf("Starting server on localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
