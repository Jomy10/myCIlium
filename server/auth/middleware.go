package auth

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"jomy.dev/mycilium/server/db"
)

func getToken(w http.ResponseWriter, r *http.Request) *string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Expected Authorization header", http.StatusUnauthorized)
		return nil
	}
	components := strings.SplitN(auth, " ", 2)

	if len(components) != 2 {
		http.Error(w, "Bearer token not specified", http.StatusBadRequest)
		return nil
	}

	log.Println(components)
	if components[0] != "Bearer" {
		http.Error(w, "Expected Bearer token", http.StatusUnauthorized)
		return nil
	}

	return &components[1]
}

func AuthPlatformsMiddleware(next func(http.ResponseWriter, *http.Request, *string, *string)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(w, r)
		if token == nil {
			return
		}

		platform, err := db.GetPlatformRights(*token)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error retrieving platform rights", http.StatusInternalServerError)
			return
		}

		next(w, r, platform, token)
	})
}

func AuthCreateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(w, r)
		if token == nil {
			return
		}

		right, err := db.GetRight(*token)
		if err != nil {
			log.Error(err)
			http.Error(w, "Couldn't query database", http.StatusInternalServerError)
			return
		}

		if right == nil {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		if *right != "create" {
			http.Error(w, "Token cannot be used to create record", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
