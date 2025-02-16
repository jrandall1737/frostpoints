package app

import (
	"fmt"
	"net/http"

	"github.com/jrandall1737/frostpoints/internal/auth"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<html><body>
		<h1>Welcome to Strava OAuth Example</h1>
		<a href="/login">
			<img src="./assets/button.png" alt="Login with Strava";">
		</a>
	</body></html>`)
}

func StartApp(port int) {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/login", auth.HandleLogin)
	http.HandleFunc("/callback", auth.HandleCallback)

	fmt.Printf("Starting server on localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
