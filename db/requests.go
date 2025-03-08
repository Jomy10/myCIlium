package db

import (
	"encoding/json"
	"errors"
	"time"
)

// Specifies a request to build `repo` at `rev` for the specified `rlatforms`
type BuildRequest struct {
	repo string
	// If emptry -> master branch
	rev       string
	platforms []string
	requestor string
}

func (key *BuildRequest) UnmarshalJSON(data []byte) error {
	var vals map[string]any
	if err := json.Unmarshal(data, &vals); err != nil {
		return err
	}

	var ok bool
	for k, v := range vals {
		switch k {
		case "repo":
			key.repo, ok = v.(string)
			if !ok {
				return errors.New("Key `repo` should be a string")
			}
		case "rev":
			key.rev, ok = v.(string)
			if !ok {
				return errors.New("Key `rev` should be a string")
			}
		case "platforms":
			var platforms []any
			platforms, ok = v.([]any)
			if !ok {
				return errors.New("Key `platforms` should be an array of strings")
			}

			for _, platform := range platforms {
				var platformString string
				platformString, ok = platform.(string)
				if !ok {
					return errors.New("Key `platforms` should be an array of strings")
				}
				key.platforms = append(key.platforms, platformString)
			}
		default:
			return errors.New("Invalid key " + k)
		}
	}

	if len(key.repo) == 0 {
		return errors.New("No repo defined")
	}

	if len(key.platforms) == 0 {
		return errors.New("No platforms defined")
	}

	return nil
}

// Add request to database, everyone with the right to read request can then
// access /requests endpoint to poll requests and start them with /request-start
func AddRequest(req BuildRequest) error {
	// TODO: all in one statement
	sqlStr := `
	INSERT INTO BuildRequest (repo, revision, platform, status, requestedBy)
	VALUES (?, ?, ?, 1, ?);
	`

	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	for _, platform := range req.platforms {
		_, err = stmt.Exec(req.repo, req.rev, platform, req.requestor)
		if err != nil {
			return err
		}
	}

	return nil
}

// Request for a specific platform
type PlatformBuildRequest struct {
	Id       int    `json:"id"`
	Repo     string `json:"repo"`
	Revision string `json:"revision"`
}

// Get all open requests for a specific platform
func GetOpenRequests(platform string) ([]PlatformBuildRequest, error) {
	sql := `
	SELECT id, repo, revision
	FROM BuildRequest
	WHERE status = 1
	  AND platform = ?
	ORDER BY requested ASC;
	`

	rows, err := db.Query(sql, platform)
	if err != nil {
		return nil, err
	}

	var id int
	var repo string
	var rev string
	var reqs []PlatformBuildRequest
	for rows.Next() {
		err = rows.Scan(&id, &repo, &rev)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, PlatformBuildRequest{
			Id:       id,
			Repo:     repo,
			Revision: rev,
		})
	}

	return reqs, nil
}

type BuildRequestStatus struct {
	Repo      string
	Rev       string
	Platform  string
	Status    string
	Requested time.Time
	Updated   time.Time
}

func GetAllRequests() ([]BuildRequestStatus, error) {
	sql := `
	SELECT repo, revision, platform, pStatus.name, requested, statusDate
	FROM BuildRequest
	INNER JOIN par_Status pStatus on pStatus.id = BuildRequest.status
	ORDER BY requested DESC;
	`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	timeFmt := "2006-01-02 15:04:05"
	var repo string
	var rev string
	var platform string
	var status string
	var requested string
	var statusDate string
	var reqs []BuildRequestStatus
	for rows.Next() {
		err = rows.Scan(&repo, &rev, &platform, &status, &requested, &statusDate)
		if err != nil {
			return nil, err
		}
		var reqTime, updTime time.Time
		reqTime, err = time.ParseInLocation(timeFmt, requested, time.UTC)
		if err != nil {
			return nil, err
		}
		updTime, err = time.ParseInLocation(timeFmt, statusDate, time.UTC)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, BuildRequestStatus{
			Repo:      repo,
			Rev:       rev,
			Platform:  platform,
			Status:    status,
			Requested: reqTime,
			Updated:   updTime,
		})
	}

	return reqs, nil
}
