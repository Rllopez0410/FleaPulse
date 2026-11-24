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

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS items (
			id TEXT PRIMARY KEY UNIQUE,
			name TEXT
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func insertItemDB(id, name string, db *sql.DB) error {
	query := `
		INSERT into 
			items (id, name)
		VALUES
			(?, ?)
		ON CONFLICT(id) DO UPDATE SET name = excluded.name;
	`
	_, err := db.Exec(query, id, name)
	if err != nil {
		return err
	}

	return nil
}

