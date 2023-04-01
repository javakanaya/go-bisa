package functionality

import "database/sql"

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", "postgres://postgres:root@localhost:5432/todorpl")
	if err != nil {
		return nil, err
	}

	return db, nil
}
