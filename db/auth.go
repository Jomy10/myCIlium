package db

func GetPlatformRights(token string) (*string, error) {
	sql := `
	SELECT platform
	FROM Token
	WHERE token = ?
	`

	rows, err := db.Query(sql, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var platform string
		err = rows.Scan(&platform)
		if err != nil {
			return nil, err
		}
		return &platform, nil
	} else {
		return nil, nil
	}
}

func GetRight(token string) (*string, error) {
	sql := `
	SELECT right
	FROM Token
	WHERE token = ?
	`

	rows, err := db.Query(sql, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var right string
		err = rows.Scan(&right)
		if err != nil {
			return nil, err
		}
		return &right, nil
	} else {
		return nil, nil
	}
}
