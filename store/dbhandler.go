package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go-final-project/config"
	"log"
	"os"
	"path/filepath"
)

func OpenDB() (*sqlx.DB, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(basePath, config.GetDBFile())
	log.Println("Path to database file:", dbPath)

	//If an error occurs, the database file does not exist and needs to be created
	_, err = os.Stat(dbPath)
	install := false
	if err != nil {
		install = true
	}

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if install {
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			date TEXT NOT NULL, 
			title TEXT NOT NULL, 
			comment TEXT, 
			repeat TEXT CHECK(length(repeat) <= 128)
		);`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, err
		}
		log.Println("The `scheduler` table was created successfully.")
	}

	return db, nil
}
