package requests

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"jomy.dev/CI/db"
)

func BuildRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Expected POST", http.StatusMethodNotAllowed)
		return
	}

	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Couldn't read request", http.StatusInternalServerError)
		return
	}

	var data db.BuildRequest
	err = json.Unmarshal(json_data, &data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("New BuildRequest", data)

	err = db.AddRequest(data)
	if err != nil {
		http.Error(w, "Error queuing request", http.StatusInternalServerError)
		return
	}
}
