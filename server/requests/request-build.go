package requests

import (
	"encoding/json"
	"io"

	log "github.com/sirupsen/logrus"

	"net/http"

	"jomy.dev/mycilium/server/db"
)

type response struct {
	Ids []db.PlatformId `json:"ids"`
}

func BuildRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Expected POST", http.StatusMethodNotAllowed)
		return
	}

	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "Couldn't read request", http.StatusInternalServerError)
		return
	}

	var data db.BuildRequest
	err = json.Unmarshal(json_data, &data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("New BuildRequest", data)

	ids, err := db.AddRequest(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error queuing request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	idsJson, err := json.Marshal(response{Ids: ids})
	if err != nil {
		log.Error(err)
		w.Write([]byte("{\"error\": \"InternalError: couldn't marshal json, but record was created\"}"))
		return
	}

	w.Write(idsJson)
}
