package db

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func getConnection() {
	var err error

	if Db != nil {
		return
	}

	// Init SQLite3 database
	Db, err = sql.Open("sqlite3", "./fiber.db")
	if err != nil {
		log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
	}

	log.Println("🚀 Connected Successfully to the Database")
}

func MakeMigrations() {
	getConnection()

	// Storage package can create this table for you at init time
	// but for the purpose of this example I created it manually
	// expanding its structure with an "u" column to better query
	// all user-related sessions.
	stmt := `CREATE TABLE IF NOT EXISTS sessions (
		k  VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		v  BLOB NOT NULL,
		e  BIGINT NOT NULL DEFAULT '0',
		u  TEXT
	);`

	_, err := Db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(64) PRIMARY KEY NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(64) NOT NULL,
		created_at TIMESTAMP DEFAULT DATETIME
	);`

	_, err = Db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

/*
https://noties.io/blog/2019/08/19/sqlite-toggle-boolean/index.html
*/
