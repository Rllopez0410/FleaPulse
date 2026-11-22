package main

import (
    "database/sql"

    _ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "tarkov.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}