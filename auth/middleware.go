package auth

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"jomy.dev/CI/db"
)

func AuthMiddleware(next func(http.ResponseWriter, *http.Request, *string, string)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Expected Authorization header", http.StatusUnauthorized)
			return
		}
		components := strings.SplitN(auth, " ", 2)

		if len(components) != 2 {
			http.Error(w, "Bearer token not specified", http.StatusBadRequest)
			return
		}

		log.Println(components)
		if components[0] != "Bearer" {
			http.Error(w, "Expected Bearer token", http.StatusUnauthorized)
			return
		}

		token := components[1]

		platform, err := db.GetPlatformRights(token)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error retrieving platform rights", http.StatusInternalServerError)
			return
		}

		next(w, r, platform, token)
	})
}
