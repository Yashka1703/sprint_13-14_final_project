package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
); 
CREATE INDEX date_index ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {

	var install bool

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Println("DB file not found, creating")
		install = true
	} else if err != nil {
		return fmt.Errorf("file verification error: %w", err)
	}

	var err error

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("db opening error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db connection error: %w", err)
	}
	log.Println("DB is connected")

	if install {
		log.Println("Creating DB schema")
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("creating db schema error: %w", err)
		}
		log.Println("DB schema has been created")
	} else {
		log.Println("DB schema already exists")
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
