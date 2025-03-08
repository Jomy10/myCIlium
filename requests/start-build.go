package requests

import (
	"net/http"

	"jomy.dev/CI/db"
)

// Returns HTTP Conflict if the request is already building or built
func StartHandler(w http.ResponseWriter, r *http.Request, allowedPlatform *string, token string) {
	HandleStatusChange(w, r, allowedPlatform, token, db.StatusStarted)
}
