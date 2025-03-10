package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"

	log "github.com/sirupsen/logrus"

	"net/http"

	"jomy.dev/mycilium/server/db"
)

func RequestsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePost(w, r)
	case "GET":
		handleGet(w, r)
	default:
		http.Error(w, "Expected a POST request for retrieving requests for a specific platform or a GET request for viewing all requests in a web browser", http.StatusMethodNotAllowed)
	}
}

type postRequest struct {
	platform *string
	status   *string
}

func (req *postRequest) UnmarshalJSON(data []byte) error {
	var vals map[string]string
	if err := json.Unmarshal(data, &vals); err != nil {
		return err
	}

	for k, v := range vals {
		switch k {
		case "platform":
			req.platform = &v
		case "status":
			req.status = &v
		default:
			return errors.New("Invalid key " + k)
		}
	}

	return nil
}

type postResult struct {
	Requests []db.PlatformBuildRequest `json:"requests"`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "Couldn't read request", http.StatusInternalServerError)
		return
	}

	var req postRequest
	err = json.Unmarshal(json_data, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var status *db.Status
	if req.status == nil {
		status = nil
	} else {
		var ok bool
		var temp_status db.Status
		temp_status, ok = db.ParseStatus(*req.status)
		status = &temp_status
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid status '%s'", *req.status), http.StatusBadRequest)
			return
		}
	}

	reqs, err := db.GetRequests(req.platform, status)
	if err != nil {
		log.Error(err)
		http.Error(w, "Couldn't retrieve requests", http.StatusInternalServerError)
		return
	}

	reqsJson, err := json.Marshal(postResult{Requests: reqs})
	if err != nil {
		log.Error(err)
		http.Error(w, "Couldn't marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(reqsJson)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	_ = r

	reqs, err := db.GetAllRequests()
	if err != nil {
		log.Error(err)
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	var t *template.Template
	t, err = template.New("get-requests").Parse(`
		<html>
			<head></head>
			<body>
				{{range .}}
					<div class="request">
						<h3 class="repo">{{ .Repo }}</h3>
						<ul>
							<li class="rev">revision: {{ .Rev }}</li>
							<li class="platform">platform: {{ .Platform }}</li>
							<li class="status">status: {{ .Status }}</li>
							<li class="requested">requested: {{ .Requested }}</li>
							<li class="updated">updated: {{ .Updated }}</li>
						</ul>
					</div>
				{{end}}
			</body>
		</html>
		`)
	if err != nil {
		log.Error(err)
		http.Error(w, "Failed creating template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, reqs)
	if err != nil {
		log.Error(err)
		http.Error(w, "Failed executing template", http.StatusInternalServerError)
		return
	}
}
