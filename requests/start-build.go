package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"jomy.dev/CI/db"
)

type requestBody struct {
	requestId *int
}

func (req *requestBody) UnmarshalJSON(data []byte) error {
	var vals map[string]int
	if err := json.Unmarshal(data, &vals); err != nil {
		return err
	}

	for k, v := range vals {
		switch k {
		case "requestId":
			req.requestId = &v
		default:
			return errors.New(fmt.Sprintf("Key `%s` is invalid", k))
		}
	}

	if req.requestId == nil {
		return errors.New("No requestId specified")
	}

	return nil
}

type Executing struct {
	mu sync.Mutex
}

func NewExecutingCache() Executing {
	return Executing{
		mu: sync.Mutex{},
	}
}

var executing Executing = NewExecutingCache()

// Returns HTTP Conflict if the request is already building or built
func StartHandler(w http.ResponseWriter, r *http.Request, allowedPlatform *string, token string) {
	if allowedPlatform == nil {
		http.Error(w, "Bearer is not authorized for any platforms", http.StatusUnauthorized)
		return
	}

	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "Couldn't read request body", http.StatusInternalServerError)
		return
	}

	var req requestBody
	err = json.Unmarshal(json_data, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// lock BuildRequest
	executing.mu.Lock()
	defer executing.mu.Unlock()

	buildRequest, err := db.GetRequestById(*req.requestId)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error retrieving request", http.StatusInternalServerError)
		return
	}

	if buildRequest == nil {
		http.Error(w, "Request with specified id doesn't exist", http.StatusNotFound)
		return
	}

	if buildRequest.Platform != *allowedPlatform {
		http.Error(w, fmt.Sprintf("Not allowed to start request of platform %s, expected %s", buildRequest.Platform, *allowedPlatform), http.StatusUnauthorized)
		return
	}

	if buildRequest.Status != db.StatusRequested {
		if buildRequest.Status == db.StatusFinished {
			http.Error(w, fmt.Sprintf("Already built %d", *req.requestId), http.StatusConflict)
		} else {
			http.Error(w, fmt.Sprintf("Already building %d", *req.requestId), http.StatusConflict)
		}
		return
	}

	err = db.SetStatus(*req.requestId, db.StatusStarted, &token)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprint("Couldn't start building"), http.StatusInternalServerError)
		return
	}

	// OK -> you can start building
}
