package db

import (
	"database/sql"
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

type PlatformId struct {
	Id       int64  `json:"requestId"`
	Platform string `json:"platform"`
}

// Add request to database, everyone with the right to read request can then
// access /requests endpoint to poll requests and start them with /request-start
func AddRequest(req BuildRequest) ([]PlatformId, error) {
	// TODO: all in one statement
	sqlStr := `
	INSERT INTO BuildRequest (repo, revision, platform, status, requestedBy)
	VALUES (?, ?, ?, 1, ?);
	`

	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	var i int64
	var res sql.Result
	var ret []PlatformId
	for _, platform := range req.platforms {
		res, err = stmt.Exec(req.repo, req.rev, platform, req.requestor)
		if err != nil {
			return nil, err
		}
		i, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ret = append(ret, PlatformId{
			Id:       i,
			Platform: platform,
		})
	}

	return ret, nil
}

// Request for a specific platform
type PlatformBuildRequest struct {
	Id       int    `json:"id"`
	Repo     string `json:"repo"`
	Revision string `json:"revision"`
	Status   Status `json:"status"`
	Platform string `json:"platform"`
}

func LoadPlatformBuildRequest(rows *sql.Rows, req *PlatformBuildRequest) error {
	if req == nil {
		return errors.New("Request struct is nil")
	}
	return rows.Scan(&req.Id, &req.Repo, &req.Revision, &req.Status, &req.Platform)
}

func LoadPlatformBuildRequests(rows *sql.Rows) ([]PlatformBuildRequest, error) {
	var req PlatformBuildRequest
	var reqs []PlatformBuildRequest
	var err error
	for rows.Next() {
		err = LoadPlatformBuildRequest(rows, &req)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

// Get all open requests for a specific platform
func GetRequests(platform *string, status *Status) ([]PlatformBuildRequest, error) {
	sql := `
	SELECT id, repo, revision, status, platform
	FROM BuildRequest
	`

	var queryParams []any

	if status != nil {
		sql += "\nWHERE status = ?"
		queryParams = append(queryParams, *status)
	}

	if platform != nil {
		if len(queryParams) > 0 {
			sql += "\n  AND platform = ?"
		} else {
			sql += "\nWHERE platform = ?"
		}
		queryParams = append(queryParams, *platform)
	}

	sql += "ORDER BY requested ASC;"

	rows, err := db.Query(sql, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []PlatformBuildRequest = []PlatformBuildRequest{}
	reqs, err = LoadPlatformBuildRequests(rows)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func GetRequestById(id int) (*PlatformBuildRequest, error) {
	sql := `
	SELECT id, repo, revision, status, platform
	FROM BuildRequest
	WHERE id = ?;
	`

	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var req PlatformBuildRequest
	if !rows.Next() {
		return nil, nil // request does not exist
	}
	err = LoadPlatformBuildRequest(rows, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
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
	defer rows.Close()

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

type Status int

const (
	StatusRequested = 1
	StatusStarted   = 2
	StatusFinished  = 3
)

func ParseStatus(status string) (Status, bool) {
	switch status {
	case "requested":
		return StatusRequested, true
	case "started":
		return StatusStarted, true
	case "finished":
		return StatusFinished, true
	default:
		return -1, false
	}
}

func SetStatus(requestId int, status Status, token *string) error {
	var res sql.Result
	var err error
	var sql string

	if status == StatusStarted {
		sql = `
		UPDATE BuildRequest
		SET status = ?,
				startedBy = ?,
				statusDate = CURRENT_TIMESTAMP
		WHERE id = ?
		`
		if token == nil {
			return errors.New("Token shouldn't be nil for status Started")
		}
		res, err = db.Exec(sql, status, token, requestId)
	} else {
		sql = `
		UPDATE BuildRequest
		SET status = ?,
				statusDate = CURRENT_TIMESTAMP
		WHERE id = ?
		`
		res, err = db.Exec(sql, status, requestId)
	}

	if err != nil {
		return err
	}

	i, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if i != 1 {
		return errors.New("No rows updated")
	}

	return nil
}
