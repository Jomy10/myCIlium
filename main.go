package main

import (
	log "github.com/sirupsen/logrus"

	"net/http"

	"jomy.dev/CI/auth"
	"jomy.dev/CI/db"
	"jomy.dev/CI/requests"
)

func main() {
	// Request for building
	http.HandleFunc("/request-build", auth.AuthCreateMiddleware(requests.BuildRequestHandler))
	// Retrieve requests
	http.HandleFunc("/requests", requests.RequestsHandler)
	// Attempt to start a request
	http.HandleFunc("/request-start", auth.AuthPlatformsMiddleware(requests.StartHandler))
	// Finish a request
	http.HandleFunc("/request-finish", auth.AuthPlatformsMiddleware(requests.FinishHandler))

	err := db.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDatabase()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
