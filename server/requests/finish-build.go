package requests

import (
	"net/http"

	"jomy.dev/mycilium/server/db"
)

func FinishHandler(w http.ResponseWriter, r *http.Request, allowedPlatform *string, token *string) {
	HandleStatusChange(w, r, allowedPlatform, token, db.StatusFinished)
}
