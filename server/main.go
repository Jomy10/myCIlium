package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"net/http"

	"jomy.dev/mycilium/server/auth"
	"jomy.dev/mycilium/server/db"
	"jomy.dev/mycilium/server/requests"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("USAGE: %s [PORT]\n", args[0])
		os.Exit(1)
	}

	port := args[1]
	if p, err := strconv.Atoi(port); err != nil || p < 0 || p > 65535 {
		fmt.Printf("'%s' is not a valid port number", port)
		os.Exit(1)
	}

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

	log.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
