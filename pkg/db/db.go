package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	dbPath := dbFile
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if _, err := os.Stat(filepath.Join("..", dbFile)); err == nil {
			dbPath = filepath.Join("..", dbFile)
		}
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		DB.Close()
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`
	_, err = DB.Exec(query)
	if err != nil {
		DB.Close()
		return err
	}
	return nil
}
